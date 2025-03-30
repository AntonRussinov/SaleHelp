package domain

import (
	"fmt"
)

// Message returns string without unwrapping error, so it can be passed to client
func (e *NotFoundError) Message() string {
	return fmt.Sprintf("%s with %s=%s not found", e.Entity, e.ParamName, e.ParamValue)
}

func NewNotFoundErr(entity, paramName, paramValue string) error {
	return &APPError{
		Code: ErrorCodeNotFound,
		Err: &NotFoundError{
			Entity:     entity,
			ParamName:  paramName,
			ParamValue: paramValue,
		},
	}
}

func NewNotFoundWithCauseErr(entity, paramName, paramValue string, err error) error {
	return &APPError{
		Code: ErrorCodeNotFound,
		Err: &NotFoundError{
			Entity:     entity,
			ParamName:  paramName,
			ParamValue: paramValue,
			Err:        err,
		},
	}
}
