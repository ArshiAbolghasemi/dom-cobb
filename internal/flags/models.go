package models

import (
	"time"

	"gorm.io/gorm"
)


type FeatureFlag struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;size:255;not null" json:"name"`
	IsActive bool   `gorm:"not null;default:false" json:"is_active"`
}

type FlagDependency struct {
	FlagID          uint       `gorm:"primaryKey;not null" json:"flag_id"`
	DependsOnFlagID uint       `gorm:"primaryKey;not null" json:"depends_on_flag_id"`
	CreatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (FeatureFlag) TableName() string {
	return "feature_flags"
}

func (FlagDependency) TableName() string {
	return "flag_dependencies"
}
