package config

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fenroe/chirpy/internal/auth"
)

type RefreshResBody struct {
	Token string `json:"token"`
}

func (C *Config) Refresh(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	refreshToken, err := C.Queries.GetRefreshToken(context.Background(), token)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		res.WriteHeader(401)
		return
	}
	if refreshToken.RevokedAt.Valid {
		res.WriteHeader(401)
		return
	}
	signedString, err := auth.MakeJWT(refreshToken.UserID, C.JWTSecret, time.Duration(3600)*time.Second)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	resBody := RefreshResBody{
		Token: signedString,
	}
	res.WriteHeader(200)
	json, _ := json.Marshal(resBody)
	res.Write(json)
}
