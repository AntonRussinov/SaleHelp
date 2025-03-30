package domain

import (
	"fmt"
)

type ForbiddenError struct {
	Message string
	Err     error
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *ForbiddenError) Unwrap() error { return e.Err }

func NewForbiddenErr(message string) error {
	return &APPError{
		Code: ErrorCodeForbidden,
		Err: &ForbiddenError{
			Message: message,
			Err:     nil,
		},
	}
}

func NewForbiddenWithCauseErr(message string, err error) error {
	return &APPError{
		Code: ErrorCodeForbidden,
		Err: &ForbiddenError{
			Message: message,
			Err:     err,
		},
	}
}
