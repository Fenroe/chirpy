package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (token string, err error) {
	token = headers.Get("Authorization")
	if token != "" {
		token = strings.ReplaceAll(token, "Bearer ", "")
		return
	}
	err = errors.New("authorization header missing")
	return
}
