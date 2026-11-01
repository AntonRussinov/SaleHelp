package ai

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNoQuota(t *testing.T) {
	assert.True(t, errors.Is(ErrNoQuota, ErrNoQuota))
	wrapped := errors.Join(ErrNoQuota, errors.New("other"))
	assert.True(t, errors.Is(wrapped, ErrNoQuota))
}
