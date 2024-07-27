package azure

import (
	"errors"
	"fmt"
)

var ErrInvalidToken = errors.New("invalid token")

func InvalidTokenError(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidToken, msg)
}
