package prettylog_test

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"salehelp/pkg/logger/handlers/prettylog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	opts := &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}
	handler := prettylog.NewHandler(opts)

	assert.NotNil(t, handler)
	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))

	handler = prettylog.NewHandler(nil)

	assert.NotNil(t, handler)
	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))
}

func TestColorize(t *testing.T) {
	t.Parallel()

	coloredText := prettylog.Colorize(prettylog.Red, "test")
	expected := "\033[31mtest\033[0m"
	assert.Equal(t, expected, coloredText)
}

func TestHandle(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log.SetOutput(&buf)

	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := prettylog.NewHandler(opts)

	testcases := []struct {
		name   string
		level  slog.Level
		output string
	}{
		{name: "debug", level: slog.LevelDebug, output: "\033[90mDEBUG:\033[0m"},
		{name: "info", level: slog.LevelInfo, output: "\033[36mINFO:\033[0m"},
		{name: "warn", level: slog.LevelWarn, output: "\033[93mWARN:\033[0m"},
		{name: "error", level: slog.LevelError, output: "\033[91mERROR:\033[0m"},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			buf.Reset()

			rec := slog.NewRecord(time.Now(), testcase.level, "Test message", 0)
			rec.AddAttrs(slog.Attr{Key: "key", Value: slog.StringValue("value")})

			err := handler.Handle(context.Background(), rec)
			require.NoError(t, err)

			assert.Contains(t, buf.String(), testcase.output)
			assert.Contains(t, buf.String(), "Test message")
			assert.Contains(t, buf.String(), `"key": "value"`)
		})
	}
}

func TestWithAttrs(t *testing.T) {
	t.Parallel()

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := prettylog.NewHandler(opts)

	newHandler := handler.WithAttrs([]slog.Attr{{Key: "key", Value: slog.StringValue("value")}})
	assert.NotNil(t, newHandler)
}

func TestWithGroup(t *testing.T) {
	t.Parallel()

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := prettylog.NewHandler(opts)

	newHandler := handler.WithGroup("group")
	assert.NotNil(t, newHandler)
}

func TestComputeAttrs(t *testing.T) {
	t.Parallel()

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := prettylog.NewHandler(opts)

	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "Test message", 0)
	rec.AddAttrs(slog.Attr{Key: "key", Value: slog.StringValue("value")})

	attrs, err := handler.ComputeAttrs(context.Background(), rec)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"key": "value",
	}

	assert.Equal(t, expected, attrs)
}

func TestSuppressDefaults(t *testing.T) {
	t.Parallel()

	fn := func(groups []string, attr slog.Attr) slog.Attr {
		return attr
	}
	suppressFn := prettylog.SuppressDefaults(fn)

	// Проверка подавления дефолтных атрибутов
	attr := suppressFn(nil, slog.Attr{Key: slog.TimeKey})
	assert.Equal(t, slog.Attr{}, attr)

	attr = suppressFn(nil, slog.Attr{Key: "customKey", Value: slog.StringValue("value")})
	assert.Equal(t, "customKey", attr.Key)
	assert.Equal(t, slog.StringValue("value"), attr.Value)
}
