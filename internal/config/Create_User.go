package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Fenroe/chirpy/internal/auth"
	"github.com/Fenroe/chirpy/internal/database"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type CreateUserError struct {
	Error string `json:"error"`
}

func (C *Config) CreateUser(res http.ResponseWriter, req *http.Request) {
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
	hash, _ := auth.HashPassword(reqBody.Password)
	params := database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hash,
	}
	user, err := C.Queries.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(CreateUserError{Error: fmt.Sprint(err)})
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
