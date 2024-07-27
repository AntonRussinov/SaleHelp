package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"salehelp/domain"
	"salehelp/pkg/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TokenValidator struct{}

func TestAuthenticateWithHeader(t *testing.T) {
	t.Parallel()

	userData := auth.UserData{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:  "Test User",
		Email: "test@test.com",
		Roles: []domain.UserRole{domain.UserRoleAdmin},
	}

	tests := []struct {
		name              string
		headerKey         string
		headerValue       string
		accessTokenHeader string
		accessTokenValue  string
		wantUserData      *auth.UserData
		wantAccessToken   string
		wantErr           bool
	}{
		{name: "no header", headerKey: "wrongHeaderKey", headerValue: "", wantErr: true},
		{name: "empty token", headerKey: "Authorization", headerValue: "", wantErr: true},
		{name: "empty token with prefix", headerKey: "Authorization", headerValue: "Bearer ", wantErr: true},
		{name: "cant validate", headerKey: "Authorization", headerValue: "Bearer 123123", wantErr: true},
		{name: "empty claims", headerKey: "Authorization", headerValue: `Bearer {}`, wantErr: true},
		{
			name: "wrong issueer", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "WrongTestIssuer", "aud": "TestAudience"}`, wantErr: true,
		},
		{
			name: "wrong audience", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": ["WrongTestAudience"]}`, wantErr: true,
		},
		{
			name: "missing audience", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer"}`, wantErr: true,
		},
		{
			name: "missing user data", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience"}`, wantErr: true,
		},
		{
			name: "missing user data", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience"}`, wantErr: true,
		},
		{
			name: "missing user name", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", ,
			"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001",
			"roles": ["SandStorm.Admin"]}`, wantErr: true,
		},
		{
			name: "missing user email", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"oid": "00000000-0000-0000-0000-000000000001",
			"roles": ["SandStorm.Admin"]}`, wantErr: true,
		},
		{
			name: "missing user id", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com",
			"roles": ["SandStorm.Admin"]}`, wantErr: true,
		},
		{
			name: "wrong format user id", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com", "oid": "123123",
			"roles": ["SandStorm.Admin"]}`, wantErr: true,
		},
		{
			name: "missing roles", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001",
			}`, wantErr: true,
		},
		{
			name: "not array roles", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001", "roles": {}}`, wantErr: true,
		},
		{
			name: "empty roles", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001",
			"roles": []}`, wantErr: true,
		},
		{
			name: "valid", headerKey: "Authorization",
			headerValue: `Bearer {"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
			"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001",
			"roles": ["SandStorm.Admin", "WrongRole", 1]}`, wantUserData: &userData, wantErr: false,
		},
	}

	authenticator := setupAuthenticator(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(tt.headerKey, tt.headerValue)

			user, _, err := authenticator.AuthenticateWithHeader(req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserData, user)
			}
		})
	}
}

func TestAuthenticateWithQuery(t *testing.T) {
	t.Parallel()

	userData := auth.UserData{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:  "Test User",
		Email: "test@test.com",
		Roles: []domain.UserRole{domain.UserRoleAdmin},
	}

	tests := []struct {
		name         string
		queryString  string
		wantUserData *auth.UserData
		wantErr      bool
	}{
		{
			name: "no token", queryString: "", wantErr: true,
		},
		{
			name: "succes", queryString: `token={"iss": "TestIssuer", "aud": "TestAudience", "name": "Test User",
		"preferred_username": "test@test.com", "oid": "00000000-0000-0000-0000-000000000001", 
		"roles": ["SandStorm.Admin", "WrongRole", 1]}`, wantUserData: &userData, wantErr: false,
		},
	}

	authenticator := setupAuthenticator(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery += tt.queryString

			user, _, err := authenticator.AuthenticateWithQuery(req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserData, user)
			}
		})
	}
}

func TestGetAccessTokenFromContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithoutCancel(context.Background())
	token := auth.GetAccessTokenFromContext(ctx)
	require.Empty(t, token)

	ctx = auth.SetAccessTokenToContext(ctx, "token")
	token = auth.GetAccessTokenFromContext(ctx)
	require.Equal(t, "token", token)
}

func setupAuthenticator(t *testing.T) auth.Authenticator {
	t.Helper()

	return &auth.AzureAuthenticator{
		Audience:  "TestAudience",
		Issuer:    "TestIssuer",
		Validator: &TokenValidator{},
	}
}

func (*TokenValidator) ValidateToken(_ context.Context, tokenString string) (map[string]any, error) {
	res := map[string]any{}

	err := json.Unmarshal([]byte(tokenString), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
