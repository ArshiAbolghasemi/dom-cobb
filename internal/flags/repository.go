package flags

import "github.com/ArshiAbolghasemi/dom-cobb/internal/database/postgres"

func GetFlagByName(name string) (*FeatureFlag, error) {
	var flag FeatureFlag
	db := postgres.GetDB()

	if err := db.Where("name = ?", name).First(&flag).Error; err != nil {
		return nil, err
	}

	return &flag, nil
}
