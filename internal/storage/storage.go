package storage

import (
	"errors"
	"fmt"
)

type (
	DuplicateError struct {
		FieldName string
		Func      string
		Err       error
	}

	ForeignKeyError struct {
		FieldName string
		Func      string
		Err       error
	}
)

var (
	ErrNotFound   = errors.New("not found")
	ErrFailedScan = errors.New("failed to scan")
)

func (e *DuplicateError) Error() string {
	return "repository." + e.Func + ": " + "duplicate key \"" + e.FieldName + "\": " + e.Err.Error()
}

func (e *DuplicateError) Unwrap() error { return e.Err }

func FailedScanError(msg string) error {
	return fmt.Errorf("%w: %s", ErrFailedScan, msg)
}

func (e *ForeignKeyError) Error() string {
	return "repository." + e.Func + ": " + "foreign key not exist \"" + e.FieldName + "\": " + e.Err.Error()
}

func (e *ForeignKeyError) Unwrap() error { return e.Err }
