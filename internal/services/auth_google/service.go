package auth_google

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"salehelp/domain"
	"salehelp/pkg/logger/sl"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

const linkTokenTTL = 10 * time.Minute
const linkTokenLength = 6

// charset для кода привязки (только цифры и заглавные буквы, без 0/O)
const linkTokenCharset = "123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

type (
	UserRepo interface {
		GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		GetByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error)
		Create(ctx context.Context, u *domain.User) error
		SetGoogleID(ctx context.Context, userID int, googleID string) error
		SetTelegramChatID(ctx context.Context, userID int, chatID int64) error
		ClearTelegramChatID(ctx context.Context, userID int) error
	}

	LinkTokensRepo interface {
		Create(ctx context.Context, t *domain.LinkToken) error
		ConsumeForWeb(ctx context.Context, token string, now time.Time) (telegramChatID int64, err error)
	}

	ShopsRepoForLink interface {
		ReassignShopsToUser(ctx context.Context, fromUserID, toUserID int) error
	}

	Service struct {
		log            *slog.Logger
		userRepo       UserRepo
		linkTokensRepo LinkTokensRepo
		shopsRepo      ShopsRepoForLink
		googleCID      string
		jwtSecret      []byte
		jwtExpire      time.Duration
		bl             *tokenBlocklist
	}

	tokenBlocklist struct {
		mu      sync.Mutex
		entries map[string]time.Time // token → expires_at
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
	Role   string `json:"role,omitempty"`
}

func NewService(log *slog.Logger, userRepo UserRepo, linkTokensRepo LinkTokensRepo, shopsRepo ShopsRepoForLink, googleClientID, jwtSecret string, jwtExpire time.Duration) (*Service, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("jwt_secret is required")
	}
	if jwtExpire <= 0 {
		jwtExpire = 24 * time.Hour
	}
	bl := &tokenBlocklist{entries: make(map[string]time.Time)}
	go bl.cleanupLoop()
	return &Service{
		log:            log.With(sl.Module("AuthGoogle")),
		userRepo:       userRepo,
		linkTokensRepo: linkTokensRepo,
		shopsRepo:      shopsRepo,
		googleCID:      googleClientID,
		jwtSecret:      []byte(jwtSecret),
		jwtExpire:      jwtExpire,
		bl:             bl,
	}, nil
}

func (bl *tokenBlocklist) revoke(token string, expiresAt time.Time) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.entries[token] = expiresAt
}

func (bl *tokenBlocklist) isRevoked(token string) bool {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	exp, ok := bl.entries[token]
	return ok && time.Now().Before(exp)
}

func (bl *tokenBlocklist) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		bl.mu.Lock()
		for tok, exp := range bl.entries {
			if now.After(exp) {
				delete(bl.entries, tok)
			}
		}
		bl.mu.Unlock()
	}
}

// VerifyGoogleIDToken проверяет id_token от Google и возвращает sub (google_id) и email.
func (s *Service) VerifyGoogleIDToken(ctx context.Context, idToken string) (sub, email string, err error) {
	if s.googleCID == "" {
		return "", "", fmt.Errorf("google auth not configured")
	}
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
		tok, err := s.issueJWT(u.ID, email, s.userRole(u))
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
			tok, err := s.issueJWT(u.ID, email, s.userRole(u))
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
	u.Role = "user"
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}
	tok, err := s.issueJWT(u.ID, email, u.Role)
	return u, tok, err
}

func (s *Service) userRole(u *domain.User) string {
	if u == nil || u.Role == "" {
		return "user"
	}
	return u.Role
}

func (s *Service) issueJWT(userID int, email, role string) (string, error) {
	if role == "" {
		role = "user"
	}
	now := time.Now()
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   fmt.Sprintf("%d", userID),
		},
		UserID: userID,
		Email:  email,
		Role:   role,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.jwtSecret)
}

// ValidateJWT проверяет наш JWT и возвращает user_id и email из claims.
func (s *Service) ValidateJWT(tokenString string) (userID int, email string, err error) {
	if s.bl.isRevoked(tokenString) {
		return 0, "", fmt.Errorf("token revoked")
	}
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

// RevokeToken добавляет токен в blocklist до его естественного истечения.
func (s *Service) RevokeToken(tokenString string) error {
	var claims JWTClaims
	tok, err := jwt.ParseWithClaims(tokenString, &claims, func(*jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return fmt.Errorf("invalid token")
	}
	s.bl.revoke(tokenString, claims.ExpiresAt.Time)
	return nil
}

// ValidateJWTAndRequireAdmin проверяет JWT и возвращает userID, email только если в claims роль admin. Иначе — ошибка.
func (s *Service) ValidateJWTAndRequireAdmin(tokenString string) (userID int, email string, err error) {
	if s.bl.isRevoked(tokenString) {
		return 0, "", fmt.Errorf("token revoked")
	}
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
	if claims.Role != "admin" {
		return 0, "", fmt.Errorf("admin role required")
	}
	return claims.UserID, claims.Email, nil
}

const bcryptCost = 12

// RegisterEmail создаёт пользователя по email и паролю, возвращает user и JWT.
func (s *Service) RegisterEmail(ctx context.Context, email, password string) (*domain.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, "", fmt.Errorf("email is required")
	}
	if len(password) < 8 {
		return nil, "", fmt.Errorf("password must be at least 8 characters")
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("get user by email: %w", err)
	}
	if existing != nil {
		return nil, "", fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}
	hashStr := string(hash)
	u := &domain.User{
		Email:        &email,
		PasswordHash: &hashStr,
	}
	u.Role = "user"
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, "", fmt.Errorf("create user: %w", err)
	}
	tok, err := s.issueJWT(u.ID, email, u.Role)
	return u, tok, err
}

// LoginEmail проверяет email и пароль, возвращает user и JWT.
func (s *Service) LoginEmail(ctx context.Context, email, password string) (*domain.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, "", fmt.Errorf("email and password are required")
	}

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("get user by email: %w", err)
	}
	if u == nil || u.PasswordHash == nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	emailVal := ""
	if u.Email != nil {
		emailVal = *u.Email
	}
	tok, err := s.issueJWT(u.ID, emailVal, s.userRole(u))
	return u, tok, err
}

// CreateLinkToken создаёт одноразовый код для привязки Telegram к аккаунту. Возвращает код и время истечения.
func (s *Service) CreateLinkToken(ctx context.Context, userID int) (linkCode string, expiresAt time.Time, err error) {
	if userID <= 0 {
		return "", time.Time{}, fmt.Errorf("user_id is required")
	}
	code, err := generateLinkCode(linkTokenLength)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate link code: %w", err)
	}
	expiresAt = time.Now().Add(linkTokenTTL)
	t := &domain.LinkToken{
		Token:     code,
		UserID:    &userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	if err := s.linkTokensRepo.Create(ctx, t); err != nil {
		return "", time.Time{}, fmt.Errorf("create link token: %w", err)
	}
	return code, expiresAt, nil
}

func generateLinkCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = linkTokenCharset[int(b[i])%len(linkTokenCharset)]
	}
	return string(b), nil
}

// ConfirmTelegramLink подтверждает привязку Telegram со стороны сайта: пользователь ввёл в боте email,
// получил код, зашёл на сайт под этим email и вводит код. Токен создан в боте (telegram_chat_id), потребляется здесь.
func (s *Service) ConfirmTelegramLink(ctx context.Context, userID int, code string) error {
	if userID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return fmt.Errorf("code is required")
	}

	chatID, err := s.linkTokensRepo.ConsumeForWeb(ctx, code, time.Now())
	if err != nil {
		return fmt.Errorf("consume link token: %w", err)
	}
	if chatID == 0 {
		return fmt.Errorf("code not found or expired")
	}

	telegramUser, err := s.userRepo.GetByTelegramChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get user by telegram chat: %w", err)
	}

	if telegramUser != nil && telegramUser.ID != userID {
		if err := s.shopsRepo.ReassignShopsToUser(ctx, telegramUser.ID, userID); err != nil {
			return fmt.Errorf("reassign shops: %w", err)
		}
		if err := s.userRepo.ClearTelegramChatID(ctx, telegramUser.ID); err != nil {
			s.log.Warn("clear old user telegram_chat_id failed", "user_id", telegramUser.ID, "err", err)
		}
	}
	if err := s.userRepo.SetTelegramChatID(ctx, userID, chatID); err != nil {
		return fmt.Errorf("set telegram_chat_id: %w", err)
	}
	return nil
}
