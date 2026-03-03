package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"salehelp/api"
	"salehelp/domain"
	"salehelp/domain/apperror"
	"salehelp/pkg/logger/sl"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func NewError(code int, msg string) *api.HTTPError {
	return &api.HTTPError{
		Code:      code,
		Message:   msg,
		Timestamp: time.Now().UTC(),
	}
}

func WrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func handleValidationErrors(c echo.Context, valErrs validator.ValidationErrors) error {
	apiErr := NewError(http.StatusBadRequest, "validation failed")
	apiErr.ValidationErrors = []api.ValidationError{}

	fieldMessages := []string{}

	for _, err := range valErrs {
		if err.Tag() == "required" {
			fieldName := err.Field()
			msg := fmt.Sprintf("%s is required", fieldName)
			apiErr.ValidationErrors = append(apiErr.ValidationErrors, api.ValidationError{
				Message: msg,
				Field:   fieldName,
			})

			fieldMessages = append(fieldMessages, msg)
		}
	}

	if len(fieldMessages) > 0 {
		apiErr.Message = apiErr.Message + ": " + strings.Join(fieldMessages, "; ")
	}

	apiErr.Details = valErrs.Error()

	return c.JSON(http.StatusBadRequest, apiErr)
}

func handleServiceValidationErrors(c echo.Context, valErrs domain.ValidationErrors) error {
	apiErr := NewError(http.StatusBadRequest, "validation failed")
	apiErr.ValidationErrors = []api.ValidationError{}

	fieldMessages := []string{}

	for _, err := range valErrs {
		apiErr.ValidationErrors = append(apiErr.ValidationErrors, api.ValidationError{
			Message: err.Message,
			Field:   err.Field,
		})

		fieldMessages = append(fieldMessages, err.Message)
	}

	if len(fieldMessages) > 0 {
		apiErr.Message = apiErr.Message + ": " + strings.Join(fieldMessages, "; ")
	}

	apiErr.Details = valErrs.Error()

	return c.JSON(http.StatusBadRequest, apiErr)
}

func handlerServiceError(c echo.Context, err *domain.ServiceError) error {
	var notFounderr *domain.NotFoundError
	if errors.As(err, &notFounderr) {
		apiErr := NewError(http.StatusNotFound, fmt.Sprintf("%s not found", notFounderr.Entity))
		return c.JSON(http.StatusNotFound, apiErr)
	}

	var appLogicErr *domain.ApplicationLogicError
	if errors.As(err, &appLogicErr) {
		apiErr := NewError(http.StatusUnprocessableEntity, appLogicErr.Msg)
		return c.JSON(http.StatusUnprocessableEntity, apiErr)
	}

	var forbiddenErr *domain.ForbiddenError
	if errors.As(err, &forbiddenErr) {
		apiErr := NewError(http.StatusForbidden, forbiddenErr.Message)
		return c.JSON(http.StatusForbidden, apiErr)
	}

	var validationErrs domain.ValidationErrors
	if errors.As(err, &validationErrs) {
		return handleServiceValidationErrors(c, validationErrs)
	}

	var duplicateErr *domain.DuplicateError
	if errors.As(err, &duplicateErr) {
		apiErr := NewError(http.StatusConflict, duplicateErr.Msg)
		return c.JSON(http.StatusConflict, apiErr)
	}

	return internalError(c)
}

func HandleAppValidationError(c echo.Context, err *apperror.ValidationError) error {
	apiErr := NewError(http.StatusBadRequest, err.Message)

	apiErr.Details = err.Error()

	return c.JSON(http.StatusBadRequest, apiErr)
}

func HandleAppError(c echo.Context, err *apperror.APPError) error {
	var validationErr *apperror.ValidationError
	if errors.As(err, &validationErr) {
		return HandleAppValidationError(c, validationErr)
	}

	var notFoundErr *apperror.NotFoundError
	if errors.As(err, &notFoundErr) {
		apiErr := NewError(http.StatusNotFound, fmt.Sprintf("%s not found", notFoundErr.Entity))
		apiErr.Details = notFoundErr.Message()

		return c.JSON(http.StatusNotFound, apiErr)
	}

	var chatErr *apperror.ChatError
	if errors.As(err, &chatErr) {
		apiErr := NewError(http.StatusUnprocessableEntity, chatErr.Message)
		return c.JSON(http.StatusUnprocessableEntity, apiErr)
	}

	return internalError(c)
}

func handleError(err error, c echo.Context) error {
	var appErr *apperror.APPError
	if errors.As(err, &appErr) {
		return HandleAppError(c, appErr)
	}

	var serviceErr *domain.ServiceError
	if errors.As(err, &serviceErr) {
		return handlerServiceError(c, serviceErr)
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return handleValidationErrors(c, ve)
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		msg := fmt.Sprintf("%v", echoErr.Message)
		return c.JSON(echoErr.Code, NewError(echoErr.Code, msg))
	}

	return internalError(c)
}

func ErrorHandler(log *slog.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		reqID := c.Response().Header().Get(echo.HeaderXRequestID)

		attributes := []any{}
		attributes = append(attributes, slog.String("request_id", reqID))
		attributes = append(attributes, slog.String("path", c.Path()))
		attributes = append(attributes, slog.String("method", c.Request().Method))
		attributes = append(attributes, sl.Err(err))

		// 404/405 уже логируются RequestLogger; не дублируем как ERROR
		if echoErr, ok := err.(*echo.HTTPError); ok && echoErr.Code >= 400 && echoErr.Code < 500 {
			log.Info("api request client error", attributes...)
		} else {
			log.Error("failed to process api request", attributes...)
		}

		if err := handleError(err, c); err != nil {
			attributes = append(attributes, sl.Err(err))
			log.Error("failed to send error response", attributes...)

			_ = c.NoContent(http.StatusInternalServerError)
		}
	}
}

func internalError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, NewError(http.StatusInternalServerError, "Internal error"))
}
