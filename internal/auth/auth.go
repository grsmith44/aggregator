package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuth = errors.New("no authorization header included")

func GetApiKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuth
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 || parts[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return parts[1], nil
}
