package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Fenroe/chirpy/internal/auth"
)

type LoginParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (C *Config) Login(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	body := LoginParams{}
	expiresInSeconds := body.ExpiresInSeconds
	if expiresInSeconds < 1 || expiresInSeconds > 3600 {
		expiresInSeconds = 3600
	}
	fmt.Println(expiresInSeconds)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(ErrorJSON{Error: err.Error()})
		res.Write(errorRes)
		return
	}
	user, err := C.Queries.GetUserByEmail(context.Background(), body.Email)
	if err != nil {
		res.WriteHeader(404)
		errorRes, _ := json.Marshal(ErrorJSON{Error: err.Error()})
		res.Write(errorRes)
		return
	}
	err = auth.CheckPasswordHash(body.Password, user.HashedPassword)
	if err != nil {
		res.WriteHeader(401)
		errorRes, _ := json.Marshal(ErrorJSON{Error: "Incorrect email or password"})
		res.Write(errorRes)
		return
	}
	token, err := auth.MakeJWT(user.ID, C.JWTSecret, (time.Duration(expiresInSeconds) * time.Second))
	if err != nil {
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(ErrorJSON{Error: "Error signing JWT"})
		res.Write(errorRes)
		return
	}
	userJSON := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	res.WriteHeader(200)
	userRes, _ := json.Marshal(userJSON)
	res.Write(userRes)
}
