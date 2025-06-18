package flags

import (
	"errors"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/postgres"
	"gorm.io/gorm"
)

func GetFlagByName(name string) (*FeatureFlag, error) {
	var flag FeatureFlag
	db := postgres.GetDB()
	err := db.Where("name = ?", name).First(&flag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &flag, nil
}
