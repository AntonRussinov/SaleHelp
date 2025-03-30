package domain

import (
	"fmt"
)

type (
	ValidationError struct {
		Err     error
		Message string
		Fields  []ValidationField
	}

	ValidationField struct {
		Field   string
		Message string
		Details string
	}
)

func NewValidationErr(message string) error {
	return &APPError{
		Code: ErrorCodeValidation,
		Err: &ValidationError{
			Err:     nil,
			Message: message,
			Fields:  nil,
		},
	}
}

func NewValidationWithCauseErr(err error, message string) error {
	return &APPError{
		Code: ErrorCodeValidation,
		Err: &ValidationError{
			Err:     err,
			Message: message,
			Fields:  nil,
		},
	}
}

func NewValidateFieldsErr(message string, fields []ValidationField) error {
	return &APPError{
		Code: ErrorCodeValidation,
		Err: &ValidationError{
			Err:     nil,
			Message: message,
			Fields:  fields,
		},
	}
}

func (e *ValidationError) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s: %v. %+q\n", e.Message, e.Err, e.Fields)
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
