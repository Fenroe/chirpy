package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Fenroe/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (C *Config) DeleteChirp(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	chirpID, _ := uuid.Parse(req.PathValue("chirpID"))
	chirp, err := C.Queries.GetChirp(context.Background(), chirpID)
	if err != nil {
		res.WriteHeader(404)
		body, _ := json.Marshal(ErrorJSON{err.Error()})
		res.Write(body)
		return
	}
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		res.WriteHeader(401)
		errorRes, _ := json.Marshal(ValidateChirpErrors{Error: err.Error()})
		res.Write(errorRes)
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, C.JWTSecret)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	nullUserID := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}
	if nullUserID != chirp.UserID {
		res.WriteHeader(403)
		return
	}
	err = C.Queries.DeleteChirp(context.Background(), chirp.ID)
	if err != nil {
		fmt.Print(err.Error())
		res.WriteHeader(500)
	}
	res.WriteHeader(204)
}
