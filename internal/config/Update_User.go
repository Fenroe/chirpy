package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Fenroe/chirpy/internal/auth"
	"github.com/Fenroe/chirpy/internal/database"
)

func (C *Config) UpdateUser(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	reqBody := CreateUserParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(CreateUserError{Error: fmt.Sprint(err)})
		res.Write(errorRes)
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
		errorRes, _ := json.Marshal(ValidateChirpErrors{Error: err.Error()})
		res.Write(errorRes)
		return
	}
	hash, _ := auth.HashPassword(reqBody.Password)
	params := database.UpdateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hash,
		ID:             userID,
	}
	user, err := C.Queries.UpdateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(CreateUserError{Error: fmt.Sprint(err)})
		res.Write(errorRes)
		return
	}
	resUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	responseJson, _ := json.Marshal(resUser)
	res.WriteHeader(200)
	res.Write(responseJson)
}
