package auth_google

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"salehelp/domain"
	"salehelp/pkg/logger/sl"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

type (
	UserRepo interface {
		GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		Create(ctx context.Context, u *domain.User) error
		SetGoogleID(ctx context.Context, userID int, googleID string) error
	}

	Service struct {
		log        *slog.Logger
		userRepo   UserRepo
		googleCID  string // Google OAuth2 Client ID (audience для id_token)
		jwtSecret  []byte
		jwtExpire  time.Duration
	}
)

// GoogleClaims — поля из проверенного Google id_token (достаём sub, email).
type GoogleClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

// JWTClaims — наши claims в JWT после логина.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int    `json:"user_id"`
	Email  string `json:"email,omitempty"`
}

func NewService(log *slog.Logger, userRepo UserRepo, googleClientID, jwtSecret string, jwtExpire time.Duration) (*Service, error) {
	if googleClientID == "" {
		return nil, fmt.Errorf("google_auth client_id is required")
	}
	if jwtSecret == "" {
		return nil, fmt.Errorf("jwt_secret is required")
	}
	if jwtExpire <= 0 {
		jwtExpire = 24 * time.Hour
	}
	return &Service{
		log:       log.With(sl.Module("AuthGoogle")),
		userRepo:  userRepo,
		googleCID: googleClientID,
		jwtSecret: []byte(jwtSecret),
		jwtExpire: jwtExpire,
	}, nil
}

// VerifyGoogleIDToken проверяет id_token от Google и возвращает sub (google_id) и email.
func (s *Service) VerifyGoogleIDToken(ctx context.Context, idToken string) (sub, email string, err error) {
	payload, err := idtoken.Validate(ctx, idToken, s.googleCID)
	if err != nil {
		return "", "", fmt.Errorf("invalid google id_token: %w", err)
	}
	sub, _ = payload.Claims["sub"].(string)
	email, _ = payload.Claims["email"].(string)
	if sub == "" {
		return "", "", fmt.Errorf("google id_token missing sub")
	}
	return sub, email, nil
}

// LoginByGoogle находит или создаёт пользователя по Google id_token и возвращает user + JWT.
func (s *Service) LoginByGoogle(ctx context.Context, idToken string) (*domain.User, string, error) {
	sub, email, err := s.VerifyGoogleIDToken(ctx, idToken)
	if err != nil {
		return nil, "", err
	}

	// 1) по google_id
	u, err := s.userRepo.GetByGoogleID(ctx, sub)
	if err != nil {
		return nil, "", fmt.Errorf("get user by google_id: %w", err)
	}
	if u != nil {
		tok, err := s.issueJWT(u.ID, email)
		return u, tok, err
	}

	// 2) по email (привязка к существующему аккаунту)
	if email != "" {
		u, err = s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, "", fmt.Errorf("get user by email: %w", err)
		}
		if u != nil {
			if err := s.userRepo.SetGoogleID(ctx, u.ID, sub); err != nil {
				s.log.Warn("set google_id for existing user failed", "user_id", u.ID, "err", err)
			}
			tok, err := s.issueJWT(u.ID, email)
			return u, tok, err
		}
	}

	// 3) создаём нового пользователя
	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}
	subPtr := sub
	u = &domain.User{
		Email:    emailPtr,
		GoogleID: &subPtr,
	}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}
	tok, err := s.issueJWT(u.ID, email)
	return u, tok, err
}

func (s *Service) issueJWT(userID int, email string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   fmt.Sprintf("%d", userID),
		},
		UserID: userID,
		Email:  email,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.jwtSecret)
}

// ValidateJWT проверяет наш JWT и возвращает user_id и email из claims.
func (s *Service) ValidateJWT(tokenString string) (userID int, email string, err error) {
	var claims JWTClaims
	tok, err := jwt.ParseWithClaims(tokenString, &claims, func(*jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("invalid jwt: %w", err)
	}
	if !tok.Valid {
		return 0, "", fmt.Errorf("invalid jwt")
	}
	return claims.UserID, claims.Email, nil
}
