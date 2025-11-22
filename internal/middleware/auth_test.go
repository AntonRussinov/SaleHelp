package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"salehelp/domain"
	"salehelp/internal/middleware"
	"salehelp/pkg/auth"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	MockAuthenticator struct {
		mock.Mock
	}
	MockUserService struct {
		mock.Mock
	}
)

func (m *MockAuthenticator) AuthenticateWithQuery(req *http.Request) (*auth.UserData, string, error) {
	args := m.Called(req)
	userData, ok := args.Get(0).(*auth.UserData)

	if ok {
		return userData, args.String(1), args.Error(2)
	}

	return nil, args.String(1), args.Error(2)
}

func (m *MockAuthenticator) AuthenticateWithHeader(req *http.Request) (*auth.UserData, string, error) {
	args := m.Called(req)
	userData, ok := args.Get(0).(*auth.UserData)

	if ok {
		return userData, args.String(1), args.Error(2)
	}

	return nil, args.String(1), args.Error(2)
}

func (m *MockAuthenticator) AuthenticateWithCookie(req *http.Request) (*auth.UserData, string, error) {
	args := m.Called(req)
	userData, ok := args.Get(0).(*auth.UserData)

	if ok {
		return userData, args.String(1), args.Error(2)
	}

	return nil, args.String(1), args.Error(2)
}

func (m *MockAuthenticator) Authenticate(_ context.Context, _ string) (*auth.UserData, string, error) {
	return nil, "0", nil
}

func (m *MockUserService) Login(ctx context.Context, data *auth.UserData) (*domain.User, error) {
	args := m.Called(ctx, data)

	user, ok := args.Get(0).(*domain.User)
	if ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func TestAuthenticationQuery(t *testing.T) {
	t.Parallel()

	e := echo.New()

	authenticator := new(MockAuthenticator)
	userService := new(MockUserService)

	middlewareFunc := middleware.AuthenticationQuery(authenticator, userService, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userData := &auth.UserData{ID: 1}
	user := &domain.User{ID: 1}
	token := "test-token"

	authenticator.On("AuthenticateWithQuery", req).Return(userData, token, nil)
	userService.On("Login", mock.Anything, userData).Return(user, nil)

	handler := middlewareFunc(func(c echo.Context) error {
		assert.Equal(t, user, auth.GetUserFromContext(c.Request().Context()))
		return c.String(http.StatusOK, "test")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	authenticator.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestAuthenticationHeader(t *testing.T) {
	t.Parallel()

	e := echo.New()

	authenticator := new(MockAuthenticator)
	userService := new(MockUserService)

	cookieUpdateFuncCalled := false
	cookieUpdateFn := func(c echo.Context, token string) {
		cookieUpdateFuncCalled = true
	}

	middlewareFunc := middleware.AuthenticationHeader(authenticator, userService, cookieUpdateFn)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accesstoken", "access-token")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userData := &auth.UserData{ID: 1}
	user := &domain.User{ID: 1}
	token := "test-token"

	authenticator.On("AuthenticateWithHeader", req).Return(userData, token, nil)
	userService.On("Login", mock.Anything, userData).Return(user, nil)

	handler := middlewareFunc(func(c echo.Context) error {
		assert.Equal(t, user, auth.GetUserFromContext(c.Request().Context()))
		assert.Equal(t, "access-token", auth.GetAccessTokenFromContext(c.Request().Context()))

		return c.String(http.StatusOK, "test")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
	assert.True(t, cookieUpdateFuncCalled)

	authenticator.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestAuthenticationCookie(t *testing.T) {
	t.Parallel()

	e := echo.New()

	authenticator := new(MockAuthenticator)
	userService := new(MockUserService)

	middlewareFunc := middleware.AuthenticationCookie(authenticator, userService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userData := &auth.UserData{ID: 1}
	user := &domain.User{ID: 1}

	authenticator.On("AuthenticateWithCookie", req).Return(userData, "", nil)
	userService.On("Login", mock.Anything, userData).Return(user, nil)

	handler := middlewareFunc(func(c echo.Context) error {
		assert.Equal(t, user, auth.GetUserFromContext(c.Request().Context()))
		return c.String(http.StatusOK, "test")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	authenticator.AssertExpectations(t)
	userService.AssertExpectations(t)
}

// func TestAuthenticationWebSockets(t *testing.T) {
// 	t.Parallel()

// 	faker := gofakeit.New(1)
// 	userID := uuid.MustParse(faker.UUID())

// 	tests := []struct {
// 		name                     string
// 		headerValue              string
// 		authenticatorCalled      bool
// 		authenticatorReturn      *auth.UserData
// 		authenticatorError       error
// 		userServiceCalled        bool
// 		userServiceReturn        *domain.User
// 		userServiceError         error
// 		wantErr                  bool
// 		wantUserInContext        *domain.User
// 		wantAccessTokenInContext string
// 		wantIDTokenInContext     string
// 	}{
// 		{name: "returns error on empty header", wantErr: true},
// 		{name: "returns error if it don't have 3 parts", headerValue: "1,2,3,4", wantErr: true},
// 		{
// 			name:        "returns error if 2d header value is not access token",
// 			headerValue: "sandstorm,wrongtoken_wrongTokenValue,idToken_idTokenValue",
// 			wantErr:     true,
// 		},
// 		{
// 			name:        "returns error if 3d header is not idToken",
// 			headerValue: "sandstorm,accessToken_accessTokenValue,wrongToken_value",
// 			wantErr:     true,
// 		},
// 		{
// 			name:                "returns error if authenticator returns error",
// 			headerValue:         "sandstorm,accessToken_accessValue,idToken_idValue",
// 			authenticatorCalled: true,
// 			authenticatorError:  assert.AnError,
// 			wantErr:             true,
// 		},
// 		{
// 			name:                "returns error if user service returns error",
// 			headerValue:         "sandstorm,accessToken_accessValue,idToken_idValue",
// 			authenticatorCalled: true,
// 			authenticatorReturn: &auth.UserData{ID: userID},
// 			userServiceCalled:   true,
// 			userServiceError:    assert.AnError,
// 			wantErr:             true,
// 		},
// 		{
// 			name:                     "sets user and tokens in context on success",
// 			headerValue:              "sandstorm,accessToken_accessValue,idToken_idValue",
// 			authenticatorCalled:      true,
// 			authenticatorReturn:      &auth.UserData{ID: userID},
// 			userServiceCalled:        true,
// 			userServiceReturn:        &domain.User{ID: userID},
// 			wantErr:                  false,
// 			wantUserInContext:        &domain.User{ID: userID},
// 			wantAccessTokenInContext: "accessValue",
// 			wantIDTokenInContext:     "idValue",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			mockAuth := auth_mocks.NewMockAuthenticator(t)
// 			mockUserService := middleware_mocks.NewMockUserService(t)

// 			e := echo.New()
// 			req := httptest.NewRequest(http.MethodGet, "/", nil)
// 			req.Header.Set("Sec-WebSocket-Protocol", tt.headerValue)

// 			rec := httptest.NewRecorder()
// 			c := e.NewContext(req, rec)

// 			mware := middleware.AuthenticationWebSockets(mockAuth, mockUserService)

// 			nextCalled := false

// 			var (
// 				accessTokenInContext string
// 				idTokenInContext     string
// 				userInContext        *domain.User
// 			)

// 			mwWithNextHandler := mware(func(c echo.Context) error {
// 				nextCalled = true
// 				accessTokenInContext = auth.GetAccessTokenFromContext(c.Request().Context())
// 				idTokenInContext = auth.GetIDTokenFromContext(c.Request().Context())
// 				userInContext = auth.GetUserFromContext(c.Request().Context())

// 				return nil
// 			})

// 			if tt.authenticatorCalled {
// 				mockAuth.EXPECT().Authenticate(mock.Anything, "idValue").Return(tt.authenticatorReturn, "", tt.authenticatorError)
// 			}

// 			if tt.userServiceCalled {
// 				mockUserService.EXPECT().Login(mock.Anything, tt.authenticatorReturn).Return(tt.userServiceReturn, tt.userServiceError)
// 			}

// 			err := mwWithNextHandler(c)

// 			if tt.wantErr {
// 				require.Error(t, err)
// 				require.False(t, nextCalled)
// 				require.Empty(t, auth.GetAccessTokenFromContext(c.Request().Context()))
// 				require.Empty(t, auth.GetIDTokenFromContext(c.Request().Context()))
// 				require.Nil(t, auth.GetUserFromContext(c.Request().Context()))
// 			} else {
// 				require.NoError(t, err)
// 				require.True(t, nextCalled)
// 				require.Equal(t, tt.wantAccessTokenInContext, accessTokenInContext)
// 				require.Equal(t, tt.wantIDTokenInContext, idTokenInContext)
// 				require.Equal(t, tt.wantUserInContext, userInContext)
// 			}
// 		})
// 	}
// }
