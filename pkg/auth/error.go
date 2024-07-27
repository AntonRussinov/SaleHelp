package auth

import (
	"errors"
	"fmt"
)

var (
	ErrNoAuthToken   = errors.New("no auth token provided")
	ErrInvalidClaims = errors.New("invalid claims")
)

func InvalidClaimsError(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidClaims, msg)
}
