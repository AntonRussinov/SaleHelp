package auth_google

import (
	"context"

	"salehelp/domain"
)

type (
	Service interface {
		LoginByGoogle(ctx context.Context, idToken string) (*domain.User, string, error)
		RegisterEmail(ctx context.Context, email, password string) (*domain.User, string, error)
		LoginEmail(ctx context.Context, email, password string) (*domain.User, string, error)
		ValidateJWT(tokenString string) (userID int, email string, err error)
		ValidateJWTAndRequireAdmin(tokenString string) (userID int, email string, err error)
	}

	Handler struct {
		service Service
	}
)

// LoginByGoogleRequest тело запроса POST /auth/google.
type LoginByGoogleRequest struct {
	IDToken string `json:"idToken" validate:"required"`
}

// LoginByGoogleResponse ответ после успешного логина.
type LoginByGoogleResponse struct {
	AccessToken string       `json:"accessToken"`
	User       UserResponse  `json:"user"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email,omitempty"`
}

// RegisterRequest тело POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest тело POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
