package domain

import "fmt"

type UnexpectedError struct {
	Err     error
	Message string
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *UnexpectedError) Unwrap() error {
	return e.Err
}

func NewUnexpectedErr(message string, err error) error {
	return &APPError{
		Code: ErrorCodeUnexpected,
		Err: &NotProcessedError{
			Err:     err,
			Message: message,
		},
	}
}
