package cron_test

import (
	"context"
	"log/slog"
	"salehelp/internal/cron"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockLogger is a mock implementation of a logger for testing purposes.
type MockLogger struct {
	mock.Mock
}

func TestNewCronRunner(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := &gorm.DB{}

	runner, err := cron.NewCronRunner(ctx, db, slogdiscard.NewDiscardLogger())
	require.NoError(t, err)
	require.NotNil(t, runner)
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) With(args ...interface{}) *slog.Logger {
	return nil
}
