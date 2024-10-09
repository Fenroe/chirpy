package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (token string, err error) {
	token = headers.Get("Authorization")
	if token != "" {
		token = strings.ReplaceAll(token, "ApiKey ", "")
		return
	}
	err = errors.New("authorization header missing")
	return
}
