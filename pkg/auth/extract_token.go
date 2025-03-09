package auth

import (
	"net/http"
	"strings"
)

func tokenHeader(req *http.Request) (string, error) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		return "", ErrNoAuthToken
	}

	return token, nil
}

func tokenQueryParam(req *http.Request) (string, error) {
	token := req.URL.Query().Get("token")
	if token == "" {
		return "", ErrNoAuthToken
	}

	return token, nil
}

func tokenCookie(req *http.Request, cookieName string) (string, error) {
	token, err := req.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	if token == nil {
		return "", nil
	}

	return token.Value, nil
}
