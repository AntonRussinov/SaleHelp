package auth_google

import (
	"context"
	"time"

	"salehelp/domain"
)

type (
	Service interface {
		LoginByGoogle(ctx context.Context, idToken string) (*domain.User, string, error)
		RegisterEmail(ctx context.Context, email, password string) (*domain.User, string, error)
		LoginEmail(ctx context.Context, email, password string) (*domain.User, string, error)
		ValidateJWT(tokenString string) (userID int, email string, err error)
		ValidateJWTAndRequireAdmin(tokenString string) (userID int, email string, err error)
		CreateLinkToken(ctx context.Context, userID int) (linkCode string, expiresAt time.Time, err error)
		ConfirmTelegramLink(ctx context.Context, userID int, code string) error
		RevokeToken(tokenString string) error
	}

	Handler struct {
		service     Service
		botUsername string
	}
)

// LoginByGoogleRequest тело запроса POST /auth/google.
type LoginByGoogleRequest struct {
	IDToken string `json:"idToken" validate:"required" example:"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJleGFtcGxlIn0.ZXhhbXBsZQ"`
}

// LoginByGoogleResponse ответ после успешного логина.
type LoginByGoogleResponse struct {
	AccessToken string      `json:"accessToken" example:"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJleGFtcGxlIn0.ZXhhbXBsZQ"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID    int    `json:"id"              example:"5"`
	Email string `json:"email,omitempty" example:"user@example.com"`
}

// RegisterRequest тело POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"    example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8"    example:"mypassword123"`
}

// LoginRequest тело POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required"       example:"mypassword123"`
}

func NewHandler(service Service, botUsername string) *Handler {
	return &Handler{service: service, botUsername: botUsername}
}
