package mock

import (
	"github.com/ArshiAbolghasemi/dom-cobb/internal/flags"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetFlagByName(name string) (*flags.FeatureFlag, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) GetFlagByIds(ids []uint) ([]*flags.FeatureFlag, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) GetFlagById(id uint) (*flags.FeatureFlag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) GetFlagDependencies(flag *flags.FeatureFlag) ([]*flags.FeatureFlag, error) {
	args := m.Called(flag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) GetFlagDependents(flag *flags.FeatureFlag) ([]*flags.FeatureFlag, error) {
	args := m.Called(flag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) GetFeatureFlagLogs(flag *flags.FeatureFlag, page, size uint) ([]*logger.LogEntry, uint, uint, error) {
	args := m.Called(flag, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(uint), args.Get(2).(uint), args.Error(3)
	}
	return args.Get(0).([]*logger.LogEntry), args.Get(1).(uint), args.Get(2).(uint), args.Error(3)
}

func (m *MockRepository) CreateFlag(name string, isActive bool, dependencies []uint) (*flags.FeatureFlag, error) {
	args := m.Called(name, isActive, dependencies)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*flags.FeatureFlag), args.Error(1)
}

func (m *MockRepository) UpdateFlag(flag *flags.FeatureFlag, isActive bool) error {
	args := m.Called(flag, isActive)
	return args.Error(0)
}
