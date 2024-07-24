package sl_test

import (
	"errors"
	"log/slog"
	"salehelp/pkg/logger/sl"
	"testing"

	"github.com/stretchr/testify/require"
)

var ErrSimple = errors.New("simple error")

func TestErr(t *testing.T) {
	t.Parallel()

	err := sl.Err(ErrSimple)

	require.Equal(t, "error", err.Key)
	require.Equal(t, slog.StringValue("simple error"), err.Value)
}

func TestModule(t *testing.T) {
	t.Parallel()

	module := sl.Module("module")

	require.Equal(t, "module", module.Key)
	require.Equal(t, slog.StringValue("module"), module.Value)
}

func TestMethod(t *testing.T) {
	t.Parallel()

	method := sl.Method("method")

	require.Equal(t, "op", method.Key)
	require.Equal(t, slog.StringValue("method"), method.Value)
}
