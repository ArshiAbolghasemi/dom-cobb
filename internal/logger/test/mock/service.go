package mock

import (
	"github.com/ArshiAbolghasemi/dom-cobb/internal/logger"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(entry *logger.LogEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockLogger) LogBatch(entries []*logger.LogEntry) error {
	args := m.Called(entries)
	return args.Error(0)
}
