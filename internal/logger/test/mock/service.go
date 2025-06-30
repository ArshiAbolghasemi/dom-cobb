package mock

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(message string, metadata map[string]any) error {
	args := m.Called(message, metadata)
	return args.Error(0)
}
