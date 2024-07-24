package apperror

import (
	"errors"
	"fmt"
)

// APP Error used only to wrap any other custom erros to ease error handling
type (
	APPError struct {
		Err error
	}

	ValidationError struct {
		Err     error
		Message string
	}

	NotFoundError struct {
		Entity     string
		ParamName  string
		ParamValue string
		Err        error
	}

	NotProcessedError struct {
		Err     error
		Message string
	}

	ChatError struct {
		IsBlocking bool
		Message    string
		Err        error
		Suppressed bool
	}
)

var (
	// generic validation error for validation errors without specific cause
	ErrValidation = errors.New("validation error")

	ErrUnsupportedModel          = errors.New("unsupported model")
	ErrDocumentsNotLoaded        = errors.New("documents is not loaded yet. please try later")
	ErrDocumentsProcessingFailed = errors.New("documents cant be processed. please try another files")
)

func (e *APPError) Error() string {
	return fmt.Sprintf("application error: %v", e.Err)
}

func (e *APPError) Unwrap() error {
	return e.Err
}

func (e *ValidationError) Error() string {
	return printWithMessage(e.Err, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s=%s not found: %v", e.Entity, e.ParamName, e.ParamValue, e.Err)
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// Message returns string without unwrapping error, so it can be passed to client
func (e *NotFoundError) Message() string {
	return fmt.Sprintf("%s with %s=%s not found", e.Entity, e.ParamName, e.ParamValue)
}

func (e *NotProcessedError) Error() string {
	return printWithMessage(e.Err, e.Message)
}

func (e *NotProcessedError) Unwrap() error {
	return e.Err
}

func (e *ChatError) Error() string {
	return printWithMessage(e.Err, e.Message)
}

func (e *ChatError) Unwrap() error {
	return e.Err
}

func Validation(message string) error {
	return &APPError{
		Err: &ValidationError{
			Err:     ErrValidation,
			Message: message,
		},
	}
}

func ValidationWithCause(err error, message string) error {
	return &APPError{
		Err: &ValidationError{
			Err:     err,
			Message: message,
		},
	}
}

func NotFound(entity, paramName, paramValue string) error {
	return &APPError{
		Err: &NotFoundError{
			Entity:     entity,
			ParamName:  paramName,
			ParamValue: paramValue,
		},
	}
}

func NotFoundWithCause(entity, paramName, paramValue string, err error) error {
	return &APPError{
		Err: &NotFoundError{
			Entity:     entity,
			ParamName:  paramName,
			ParamValue: paramValue,
			Err:        err,
		},
	}
}

func NotProcessed(message string) error {
	return &APPError{
		Err: &NotProcessedError{
			Err:     ErrValidation,
			Message: message,
		},
	}
}

func NotProcessedWithCause(err error, message string) error {
	return &APPError{
		Err: &NotProcessedError{
			Err:     err,
			Message: message,
		},
	}
}

func ChatAICallError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: false,
			Message:    "Ai call error: " + err.Error(),
			Err:        err,
		},
	}
}

func ChatFilesIsLoadingError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: true,
			Message:    "Files are still being loaded. Please try later",
			Err:        err,
		},
	}
}

func ChatLoadingFilesError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: true,
			Message:    "Error loading files. Please try reloading the page. If the problem persists, kindly reach out to our Support team",
			Err:        err,
		},
	}
}

func ChatConnectionError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: true,
			Message:    "Connection Error. Please try again later. If the issue persists, consider contacting our Support Team.",
			Err:        err,
		},
	}
}

func ChatQuotaExceededError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: true,
			Message:    "Quota Exceeded. To increase your limit, please submit a request.",
			Err:        err,
		},
	}
}

func ChatTemporaryConnectionError(err error) error {
	return &APPError{
		Err: &ChatError{
			IsBlocking: false,
			Message:    "Connection Error. Please try again. If the problem persists, feel free to contact our Support Team.",
			Err:        err,
		},
	}
}

func ChatSupressedError(err error) error {
	return &APPError{
		Err: &ChatError{
			Err:        err,
			Suppressed: true,
		},
	}
}

func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func printWithMessage(err error, message string) string {
	return fmt.Sprintf("%s: %v", message, err)
}
