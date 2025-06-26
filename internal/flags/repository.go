package flags

import (
	"errors"
	"sync"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/postgres"
	"gorm.io/gorm"
)

type IRepository interface {
	GetFlagByName(name string) (*FeatureFlag, error)
	GetFlagByIds(flagIds []uint) ([]FeatureFlag, error)
	CreateFlag(name string, active bool, dependecnyFlagIds []uint) (*FeatureFlag, error)
}

type Repository struct {
	db *gorm.DB
}

var (
	repo     IRepository
	onceRepo sync.Once
)

func GetRepository() IRepository {
	onceRepo.Do(func() {
		repo = &Repository{
			db: postgres.GetDB(),
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

func (r *Repository) GetFlagByIds(flagIds []uint) ([]FeatureFlag, error) {
	var flags []FeatureFlag
	err := r.db.Where("id IN ?", flagIds).Find(&flags).Error
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
