package auth

import (
	"context"
	"fmt"
	"net/http"
	"salehelp/domain"
	"salehelp/pkg/auth/azure"
	"slices"
)

type (
	Config struct {
		AppID      string
		TennantID  string
		EnableAuth bool
	}

	UserData struct {
		ID    int
		Name  string
		Email string
		// Roles []domain.UserRole
	}
	Authenticator interface {
		AuthenticateWithQuery(req *http.Request) (*UserData, string, error)
		AuthenticateWithHeader(req *http.Request) (*UserData, string, error)
		AuthenticateWithCookie(req *http.Request) (*UserData, string, error)
		Authenticate(ctx context.Context, token string) (*UserData, string, error)
	}

	TokenValidator interface {
		ValidateToken(ctx context.Context, tokenString string) (map[string]any, error)
	}

	AzureAuthenticator struct {
		Audience  string
		Issuer    string
		Validator TokenValidator
	}

	key string
)

const (
	userContextKey        key = "user"
	accessTokenContextKey key = "accessToken"
	idTokenContextKey     key = "idToken"

	issuerTemplate string = "https://login.microsoftonline.com/%s/v2.0"

	SecureCookieName   = "session"
	InsecureCookieName = "session"
)

func NewAzureAuthenticator(cfg *Config) (*AzureAuthenticator, error) {
	cache, err := azure.NewCache(cfg.TennantID, cfg.AppID)
	if err != nil {
		return nil, err
	}

	az := azure.NewAuthentificator(cache)

	return &AzureAuthenticator{
		Validator: az,
		Audience:  cfg.AppID,
		Issuer:    fmt.Sprintf(issuerTemplate, cfg.TennantID),
	}, nil
}

func (a *AzureAuthenticator) AuthenticateWithHeader(req *http.Request) (*UserData, string, error) {
	token, err := tokenHeader(req)
	if err != nil {
		return nil, "", err
	}

	ctx := req.Context()

	return a.Authenticate(ctx, token)
}

func (a *AzureAuthenticator) AuthenticateWithQuery(req *http.Request) (*UserData, string, error) {
	token, err := tokenQueryParam(req)
	if err != nil {
		return nil, "", err
	}

	ctx := req.Context()

	return a.Authenticate(ctx, token)
}

func (a *AzureAuthenticator) AuthenticateWithCookie(req *http.Request) (*UserData, string, error) {
	token, err := tokenCookie(req, SecureCookieName)
	if err != nil {
		return nil, "", err
	}

	ctx := req.Context()

	return a.Authenticate(ctx, token)
}

func (a *AzureAuthenticator) Authenticate(ctx context.Context, token string) (*UserData, string, error) {
	claims, err := a.Validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, "", err
	}

	err = a.ValidateClaims(claims)
	if err != nil {
		return nil, "", err
	}

	user, err := GetUserFromClaims(claims)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func extractAudienceFromClaims(claims map[string]any) ([]string, error) {
	audiences := []string{}

	aud, ok := claims["aud"]
	if !ok {
		return nil, InvalidClaimsError("could not parse audience")
	}

	switch a := aud.(type) {
	case string:
		audiences = append(audiences, a)
	case []string:
		audiences = append(audiences, a...)
	case []any:
		for _, i := range a {
			v, ok := i.(string)
			if !ok {
				continue
			}

			audiences = append(audiences, v)
		}
	default:
		return nil, InvalidClaimsError("invalid audience")
	}

	return audiences, nil
}

func (a *AzureAuthenticator) ValidateClaims(claims map[string]any) error {
	iss, ok := claims["iss"].(string)
	if !ok {
		return InvalidClaimsError("could not parse issuer")
	}

	if iss != a.Issuer {
		return InvalidClaimsError("invalid issuer")
	}

	audiences, err := extractAudienceFromClaims(claims)
	if err != nil {
		return err
	}

	if !slices.Contains(audiences, a.Audience) {
		return InvalidClaimsError("invalid audience")
	}

	return nil
}

func GetUserFromClaims(claims map[string]any) (*UserData, error) {
	name, ok := claims["name"].(string)
	if !ok {
		return nil, InvalidClaimsError("could not parse name")
	}

	email, ok := claims["preferred_username"].(string)
	if !ok {
		return nil, InvalidClaimsError("could not parse email")
	}

	userID := 1

	// userRoles := extractUserRoles(claims)

	// if len(userRoles) == 0 {
	// 	return nil, InvalidClaimsError("no user roles provided")
	// }

	user := &UserData{
		ID:    userID,
		Name:  name,
		Email: email,
		// Roles: userRoles,
	}

	return user, nil
}

// func extractUserRoles(claims map[string]any) []domain.UserRole {
// 	userRoles := []domain.UserRole{}

// 	tokenRoles, ok := claims["roles"].([]any)
// 	if !ok {
// 		return userRoles
// 	}

// 	for _, r := range tokenRoles {
// 		role, ok := r.(string)
// 		if !ok {
// 			continue
// 		}

// 		userRole := domain.UserRole(role)
// 		if userRole.IsValid() {
// 			userRoles = append(userRoles, userRole)
// 		}
// 	}

// 	return userRoles
// }

func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil
	}

	return user
}

func SetUserToContext(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func SetAccessTokenToContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, token)
}

func SetIDTokenToContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, idTokenContextKey, token)
}

func GetAccessTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(accessTokenContextKey).(string)
	if !ok {
		return ""
	}

	return token
}

func GetIDTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(idTokenContextKey).(string)
	if !ok {
		return ""
	}

	return token
}
