package config

import (
	"context"
	"net/http"

	"github.com/Fenroe/chirpy/internal/auth"
)

func (C *Config) Revoke(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	err = C.Queries.RevokeRefreshToken(context.Background(), token)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	res.WriteHeader(204)
}
