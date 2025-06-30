package mock

import (
	"time"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/flags"
	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

type FeatureFlagOption func(*flags.FeatureFlag)

func WithIsActive(isActive bool) FeatureFlagOption {
	return func(f *flags.FeatureFlag) {
		f.IsActive = isActive
	}
}

func WithName(name string) FeatureFlagOption {
	return func(f *flags.FeatureFlag) {
		f.Name = name
	}
}

func WithId(id uint) FeatureFlagOption {
	return func(f *flags.FeatureFlag) {
		f.Model.ID = id
	}
}

func CreateFeatureFlagByIds(ids []uint, opts ...FeatureFlagOption) []*flags.FeatureFlag {
	var f []*flags.FeatureFlag
	for _, id := range ids {
		flag := CreateFeatureFlag(WithId(id))

		for _, opt := range opts {
			opt(flag)
		}

		f = append(f, flag)
	}
	return f
}

func CreateFeatureFlag(opts ...FeatureFlagOption) *flags.FeatureFlag {
	flag := &flags.FeatureFlag{
		Model: gorm.Model{
			ID:        gofakeit.Uint(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     gofakeit.Word(),
		IsActive: gofakeit.Bool(),
	}

	for _, opt := range opts {
		opt(flag)
	}

	return flag
}
