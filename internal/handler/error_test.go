package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"salehelp/api"
	"salehelp/domain"
	"salehelp/domain/apperror"
	"salehelp/internal/handler"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type badSerializer struct{}

var (
	ErrCantSerialize = errors.New("cant serialize")
	ErrSimple        = errors.New("simple error")
	ErrEmpty         = errors.New("")
)

func (badSerializer) Serialize(_ echo.Context, _ any, _ string) error {
	return ErrCantSerialize
}

func (badSerializer) Deserialize(_ echo.Context, _ any) error {
	return ErrCantSerialize
}

func setupHandlerBadSerializer(t *testing.T) (*echo.Echo, func(err error, c echo.Context)) {
	t.Helper()

	log := slogdiscard.NewDiscardLogger()

	e := echo.New()
	e.JSONSerializer = badSerializer{}
	errHandler := handler.ErrorHandler(log)

	return e, errHandler
}

func setupHandler(t *testing.T) (*echo.Echo, func(err error, c echo.Context)) {
	t.Helper()

	log := slogdiscard.NewDiscardLogger()

	e := echo.New()
	errHandler := handler.ErrorHandler(log)

	return e, errHandler
}

func TestErrorHandlerInvalidError(t *testing.T) {
	t.Parallel()

	e, handle := setupHandlerBadSerializer(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handle(ErrSimple, c)
	require.Empty(t, rec.Body.String())
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		expectCode  int
		bodyContent string
	}{
		{name: "simple error", err: ErrSimple, expectCode: http.StatusInternalServerError, bodyContent: ""},
		{name: "empty error", err: ErrEmpty, expectCode: http.StatusInternalServerError, bodyContent: "Internal error"},
		//	{name: "not found error", err: domain.NewNotFoundError("", "template", "id", "1", storage.ErrNotFound), expectCode: http.StatusNotFound, bodyContent: "not found"},
		// {name: "forbidden", err: domain.NewForbiddenError("", "cant delete", storage.ErrCantDeletePredefined), expectCode: http.StatusForbidden, bodyContent: "cant delete"},
		{name: "echo error", err: &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "echo error bad request",
		}, expectCode: http.StatusBadRequest, bodyContent: "echo error bad request"},
		{name: "unknown domain error", err: &domain.ServiceError{Fn: "some fn", Err: ErrSimple, Msg: "unknown error"}, expectCode: http.StatusInternalServerError, bodyContent: "Internal error"},
		// {name: "duplicate error", err: domain.NewDuplicateError("", "already exist", "name", &storage.DuplicateError{
		// 	Err: assert.AnError,
		// }), expectCode: http.StatusConflict, bodyContent: "already exist"},
	}

	e, handle := setupHandler(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handle(tt.err, c)

			respBody := rec.Body.String()
			require.Equal(t, tt.expectCode, rec.Code)
			require.Contains(t, respBody, tt.bodyContent)

			var apiErr api.HTTPError
			err := json.Unmarshal([]byte(respBody), &apiErr)
			require.NoError(t, err)

			require.Equal(t, tt.expectCode, apiErr.Code)
			require.Nil(t, apiErr.ValidationErrors)
		})
	}
}

func TestErrorHandlerValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		validatedStruct any
		tag             string
		expectCode      int
		bodyContent     string
		expectedErrors  []api.ValidationError
	}{
		{
			name: "two required fields missing",
			validatedStruct: struct {
				Name string `validate:"required"`
				Age  int    `validate:"required"`
			}{Name: "", Age: 0},
			expectCode:  http.StatusBadRequest,
			bodyContent: "validation failed",
			expectedErrors: []api.ValidationError{
				{Field: "Name", Message: "Name is required"},
				{Field: "Age", Message: "Age is required"},
			},
		},
		{
			name: "one require field missing",
			validatedStruct: struct {
				Name string `validate:"required"`
				Age  int    `validate:"required"`
			}{Name: "John", Age: 0},
			expectCode:  http.StatusBadRequest,
			bodyContent: "validation failed",
			expectedErrors: []api.ValidationError{
				{Field: "Age", Message: "Age is required"},
			},
		},
	}

	e, handle := setupHandler(t)
	val := validator.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := val.Struct(tt.validatedStruct)

			handle(err, c)

			respBody := rec.Body.String()
			require.Equal(t, tt.expectCode, rec.Code)
			require.Contains(t, respBody, tt.bodyContent)

			var apiErr api.HTTPError
			err = json.Unmarshal([]byte(respBody), &apiErr)
			require.NoError(t, err)

			require.Equal(t, tt.expectCode, apiErr.Code)
			require.Equal(t, len(tt.expectedErrors), len(apiErr.ValidationErrors))

			for i, expectedErr := range tt.expectedErrors {
				require.Equal(t, expectedErr.Field, apiErr.ValidationErrors[i].Field)
				require.Equal(t, expectedErr.Message, apiErr.ValidationErrors[i].Message)
			}
		})
	}
}

func TestErrorHandlerServiceValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectCode     int
		bodyContent    string
		expectedErrors []api.ValidationError
	}{
		{name: "no detailed errors", err: domain.NewValidationErrors("", "name", []domain.FieldValidationError{}), expectCode: http.StatusBadRequest, bodyContent: "validation failed", expectedErrors: []api.ValidationError{}},
		{name: "2 detailed errors", err: domain.NewValidationErrors("", "name", []domain.FieldValidationError{
			{Field: "name", Message: "name is required"},
			{Field: "name", Message: "name is too short"},
		}), expectCode: http.StatusBadRequest, bodyContent: "validation failed", expectedErrors: []api.ValidationError{
			{Field: "name", Message: "name is required"},
			{Field: "name", Message: "name is too short"},
		}},
	}

	e, handle := setupHandler(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handle(tt.err, c)

			respBody := rec.Body.String()
			require.Equal(t, tt.expectCode, rec.Code)
			require.Contains(t, respBody, tt.bodyContent)

			var apiErr api.HTTPError
			err := json.Unmarshal([]byte(respBody), &apiErr)
			require.NoError(t, err)

			require.Equal(t, tt.expectCode, apiErr.Code)
			require.Equal(t, len(tt.expectedErrors), len(apiErr.ValidationErrors))

			for i, expectedErr := range tt.expectedErrors {
				require.Equal(t, expectedErr.Field, apiErr.ValidationErrors[i].Field)
				require.Equal(t, expectedErr.Message, apiErr.ValidationErrors[i].Message)
			}
		})
	}
}

func TestErrorHandlerValidationErrorsEmptySlice(t *testing.T) {
	t.Parallel()

	expectedCode := http.StatusBadRequest

	e, handle := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handle(validator.ValidationErrors{}, c)

	respBody := rec.Body.String()
	require.Equal(t, expectedCode, rec.Code)
	require.Contains(t, respBody, "validation failed")

	var apiErr api.HTTPError
	err := json.Unmarshal([]byte(respBody), &apiErr)
	require.NoError(t, err)

	require.Equal(t, expectedCode, apiErr.Code)
	require.NotNil(t, apiErr.ValidationErrors)
	require.Empty(t, apiErr.ValidationErrors)
}

func TestHandleAppValidationError(t *testing.T) {
	t.Parallel()

	e, _ := setupHandler(t)

	validationErr := &apperror.ValidationError{Message: "validation error message"}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HandleAppValidationError(c, validationErr)
	require.NoError(t, err)

	respBody := rec.Body.String()
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, respBody, "validation error message")

	var apiErr api.HTTPError
	err = json.Unmarshal([]byte(respBody), &apiErr)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, apiErr.Code)
	require.Nil(t, apiErr.ValidationErrors)
}

func TestHandleAppError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		expectCode  int
		bodyContent string
	}{
		{
			name:        "validation error",
			err:         &apperror.ValidationError{Message: "validation error message"},
			expectCode:  http.StatusBadRequest,
			bodyContent: "validation error message",
		},
		{
			name:        "not found error",
			err:         &apperror.NotFoundError{Entity: "resource"},
			expectCode:  http.StatusNotFound,
			bodyContent: "resource not found",
		},
		{
			name:        "chat error",
			err:         &apperror.ChatError{Message: "chat error message"},
			expectCode:  http.StatusUnprocessableEntity,
			bodyContent: "chat error message",
		},
		{
			name:        "unknown error",
			err:         assert.AnError,
			expectCode:  http.StatusInternalServerError,
			bodyContent: "Internal error",
		},
	}

	e, _ := setupHandler(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var appErr *apperror.APPError
			if errors.As(tt.err, &appErr) {
				err := handler.HandleAppError(c, appErr)
				require.NoError(t, err)
			} else {
				appErr = &apperror.APPError{Err: tt.err}
				err := handler.HandleAppError(c, appErr)
				require.NoError(t, err)
			}

			// err := handler.HandleAppError(c, tt.err)
			// require.NoError(t, err)

			respBody := rec.Body.String()
			require.Equal(t, tt.expectCode, rec.Code)
			require.Contains(t, respBody, tt.bodyContent)

			var apiErr api.HTTPError
			err := json.Unmarshal([]byte(respBody), &apiErr)
			require.NoError(t, err)

			require.Equal(t, tt.expectCode, apiErr.Code)
			require.Nil(t, apiErr.ValidationErrors)
		})
	}
}
