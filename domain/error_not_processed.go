package domain

import (
	"fmt"
)

type NotProcessedError struct {
	Err     error
	Message string
}

func (e *NotProcessedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *NotProcessedError) Unwrap() error {
	return e.Err
}

func NewNotProcessedErr(message string) error {
	return &APPError{
		Code: ErrorCodeNotProcessed,
		Err: &NotProcessedError{
			Err:     nil,
			Message: message,
		},
	}
}

func NewNotProcessedWithCauseErr(err error, message string) error {
	return &APPError{
		Code: ErrorCodeNotProcessed,
		Err: &NotProcessedError{
			Err:     err,
			Message: message,
		},
	}
}
