package logger_test

import (
	"context"
	"log/slog"
	"os"
	"salehelp/pkg/logger"
	"salehelp/pkg/logger/handlers/prettylog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		env              string
		expectedHandler  interface{}
		expectedLogLevel slog.Level
	}{
		{
			name:             "local environment",
			env:              "local",
			expectedHandler:  &prettylog.Handler{},
			expectedLogLevel: slog.LevelDebug,
		},
		{
			name:             "dev environment",
			env:              "dev",
			expectedHandler:  &slog.JSONHandler{},
			expectedLogLevel: slog.LevelDebug,
		},
		{
			name:             "prod environment",
			env:              "prod",
			expectedHandler:  &slog.JSONHandler{},
			expectedLogLevel: slog.LevelDebug,
		},
		{
			name:             "default environment",
			env:              "unknown",
			expectedHandler:  &slog.JSONHandler{},
			expectedLogLevel: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			log := logger.New(tt.env)
			require.NotNil(t, log)

			handler := log.Handler()
			assert.IsType(t, tt.expectedHandler, handler)

			switch h := handler.(type) {
			case *prettylog.Handler:
				assert.True(t, h.Enabled(context.Background(), tt.expectedLogLevel))
			case *slog.JSONHandler:
				assert.True(t, h.Enabled(context.Background(), tt.expectedLogLevel))
			}
		})
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "log_test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	oldStdout := os.Stdout
	os.Stdout = tmpFile

	defer func() { os.Stdout = oldStdout }()

	log := logger.New("prod")
	require.NotNil(t, log)

	handler := log.Handler()
	jsonHandler, ok := handler.(*slog.JSONHandler)
	assert.True(t, ok, "handler should be slog.JSONHandler")
	assert.True(t, jsonHandler.Enabled(context.Background(), slog.LevelDebug))
}
