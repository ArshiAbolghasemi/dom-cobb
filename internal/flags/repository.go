package flags

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/mongodb"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/postgres"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type IRepository interface {
	GetFlagByName(name string) (*FeatureFlag, error)
	GetFlagByIds(flagIds []uint) ([]*FeatureFlag, error)
	GetFlagById(flagId uint) (*FeatureFlag, error)
	GetFlagDependencies(flag *FeatureFlag) ([]*FeatureFlag, error)
	GetFlagDependents(flag *FeatureFlag) ([]*FeatureFlag, error)
	GetFeatureFlagLogs(flag *FeatureFlag, limit, offset uint) ([]*logger.LogEntry, uint, uint, error)
	CreateFlag(name string, active bool, dependecnyFlagIds []uint) (*FeatureFlag, error)
	UpdateFlag(flag *FeatureFlag, active bool) error
}

type Repository struct {
	db         *gorm.DB
	collection *mongo.Collection
}

var (
	repo     IRepository
	onceRepo sync.Once
)

func GetRepository() IRepository {
	onceRepo.Do(func() {
		loggerCollection, err := logger.GetCollection()
		if err != nil {
			panic("Failed to get logger collection: " + err.Error())
		}
		repo = &Repository{
			db: postgres.GetDB(),
			collection: mongodb.GetCollection(loggerCollection),
		}
	})
	return repo
}

func (r *Repository) GetFlagByName(name string) (*FeatureFlag, error) {
	var flag FeatureFlag
	err := r.db.Where("name = ?", name).First(&flag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &flag, nil
}

func (r *Repository) GetFlagByIds(flagIds []uint) ([]*FeatureFlag, error) {
	var flags []*FeatureFlag
	err := r.db.Where("id IN ?", flagIds).Find(&flags).Error
	if err != nil {
		return nil, err
	}

	return flags, nil
}

func (r *Repository) GetFlagById(flagId uint) (*FeatureFlag, error) {
	var flag FeatureFlag
	err := r.db.Where("id = ?", flagId).First(&flag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("flag with id %d not found", flagId)
		}
		return nil, err
	}
	return &flag, nil
}

func (r *Repository) GetFlagDependencies(flag *FeatureFlag) ([]*FeatureFlag, error) {
	var flags []*FeatureFlag

	err := r.db.Table("feature_flags").
		Joins("JOIN flag_dependencies ON feature_flags.id = flag_dependencies.depends_on_flag_id").
		Where("flag_dependencies.flag_id = ?", flag.ID).
		Find(&flags).Error

	if err != nil {
		return nil, err
	}

	return flags, nil
}

func (r *Repository) GetFlagDependents(flag *FeatureFlag) ([]*FeatureFlag, error) {
	var flags []*FeatureFlag

	err := r.db.Table("feature_flags").
		Joins("JOIN flag_dependencies ON feature_flags.id = flag_dependencies.flag_id").
		Where("flag_dependencies.depends_on_flag_id = ?", flag.ID).
		Find(&flags).Error

	if err != nil {
		return nil, err
	}

	return flags, nil
}

func (r *Repository) CreateFlag(name string, active bool, dependecnyFlagIds []uint) (*FeatureFlag, error) {
	flag := FeatureFlag{
		Name:     name,
		IsActive: active,
	}

	if len(dependecnyFlagIds) == 0 {
		err := r.db.Create(&flag).Error
		return &flag, err
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(&flag).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var dependencyFlags []FlagDependency
	for _, depFlagID := range dependecnyFlagIds {
		dependencyFlags = append(dependencyFlags, FlagDependency{
			FlagID:          flag.ID,
			DependsOnFlagID: depFlagID,
		})
	}
	if err := tx.Create(&dependencyFlags).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &flag, nil
}

func (r *Repository) UpdateFlag(flag *FeatureFlag, active bool) error {
	err := r.db.Model(flag).Update("is_active", active).Error
	if err != nil {
		return err
	}

	flag.IsActive = active

	return nil
}

func (r *Repository) GetFeatureFlagLogs(flag *FeatureFlag, limit, offset uint) ([]*logger.LogEntry, uint, uint, error) {
	ctx := context.Background()
	pager := &mongodb.Pager{
		Limit: limit,
		Offset: offset,
	}
	
	filter := bson.M{"metadata.flag_id": flag.ID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}
	pager.SetTotal(uint(total))
		
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)
	
	var logs []*logger.LogEntry
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, 0, err
	}
	
	return logs, pager.Total, pager.TotalPages, nil
}
