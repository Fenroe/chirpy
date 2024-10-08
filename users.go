package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type createUserReqBody struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	reqBody := createUserReqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(errorVals{Error: err})
		res.Write(errorRes)
		return
	}
	user, err := cfg.databaseQueries.CreateUser(context.Background(), reqBody.Email)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(errorVals{Error: err})
		res.Write(errorRes)
		return
	}
	resUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	responseJson, _ := json.Marshal(resUser)
	res.WriteHeader(201)
	res.Write(responseJson)
}
