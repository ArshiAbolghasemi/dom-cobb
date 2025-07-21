package flags

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
	GetFeatureFlagLogs(flag *FeatureFlag, page, size uint) ([]*logger.LogEntry, uint, uint, error)
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
			db:         postgres.GetDB(),
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
	if active {
		return r.activateFlag(flag)
	}
	return r.deactivateFlag(flag)
}

func (r *Repository) activateFlag(flag *FeatureFlag) error {
	err := r.db.Model(flag).Update("is_active", true).Error
	if err != nil {
		return err
	}
	flag.IsActive = true
	return nil
}

func (r *Repository) deactivateFlag(flag *FeatureFlag) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	allTransitiveDependents, err := r.getAllTransitiveDependents(flag)
	if err != nil {
		tx.Rollback()
		return err
	}

	flagIDs := []uint{flag.ID}
	for _, dependent := range allTransitiveDependents {
		flagIDs = append(flagIDs, dependent.ID)
	}

	err = tx.Model(&FeatureFlag{}).Where("id IN ? AND is_active = true", flagIDs).Update("is_active", false).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	logEntries := make([]*logger.LogEntry, 0, len(allTransitiveDependents))
	for _, flagDependent := range allTransitiveDependents {
		flagDependent.IsActive = false
		logEntries = append(logEntries, &logger.LogEntry{
			Message: "Flag is auto disabled",
			Metadata: map[string]any{
				"flag_id":           flagDependent.ID,
				"dependecy_flag_id": flag.ID,
			},
			Timestamp: time.Now(),
		})
	}
	logger.NewService().LogBatch(logEntries)

	flag.IsActive = false
	return nil
}

func (r *Repository) getAllTransitiveDependents(flag *FeatureFlag) ([]*FeatureFlag, error) {
	var dependentFlags []*FeatureFlag
	err := r.db.Raw(`
		WITH RECURSIVE dependents AS (
			SELECT flag_id as id
			FROM flag_dependencies 
			WHERE depends_on_flag_id = ?
			
			UNION ALL

			SELECT fd.flag_id as id
			FROM flag_dependencies fd
			INNER JOIN dependents d ON fd.depends_on_flag_id = d.id
		)
		SELECT f.* FROM feature_flags f
		INNER JOIN dependents d ON f.id = d.id
		ORDER BY f.id
	`, flag.ID).Scan(&dependentFlags).Error
	return dependentFlags, err
}

func (r *Repository) GetFeatureFlagLogs(flag *FeatureFlag, page, size uint) ([]*logger.LogEntry, uint, uint, error) {
	ctx := context.Background()
	pager := &mongodb.Pager{
		Page: page,
		Size: size,
	}

	filter := bson.M{"metadata.flag_id": flag.ID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}
	pager.SetTotal(uint(total))

	findOptions := options.Find()
	findOptions.SetLimit(int64(pager.GetLimit()))
	findOptions.SetSkip(int64(pager.GetOffset()))
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
