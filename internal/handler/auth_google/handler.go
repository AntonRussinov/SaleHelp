package auth_google

import (
	"context"

	"salehelp/domain"
)

type (
	Service interface {
		LoginByGoogle(ctx context.Context, idToken string) (*domain.User, string, error)
		ValidateJWT(tokenString string) (userID int, email string, err error)
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

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
