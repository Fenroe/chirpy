package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Fenroe/chirpy/internal/auth"
	"github.com/Fenroe/chirpy/internal/database"
	"github.com/google/uuid"
)

type ValidateChirpParams struct {
	Body string `json:"body"`
}

type ValidateChirpErrors struct {
	Error string `json:"error"`
}

type ValidateChripRes struct {
	CleanedBody string `json:"cleaned_body"`
}

type CreateChirpRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

var profanity = []string{"kerfuffle", "sharbert", "formax"}

func validateChirp(chirp string) (cleanedChirp string, isValid bool) {
	isValid = len(chirp) <= 140
	if !isValid {
		return chirp, isValid
	}
	words := strings.Fields(chirp)
	for i, value := range words {
		lower := strings.ToLower(value)
		for _, badWord := range profanity {
			if lower == badWord {
				words[i] = "****"
			}
		}
	}
	cleanedChirp = strings.Join(words, " ")
	return cleanedChirp, isValid
}

func (C *Config) CreateChirp(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	body := ValidateChirpParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(ValidateChirpErrors{Error: fmt.Sprint(err)})
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
	cleanedChirp, isValid := validateChirp(body.Body)
	if !isValid {
		res.WriteHeader(400)
		errorRes, _ := json.Marshal(ValidateChirpErrors{Error: "Chirp is too long"})
		res.Write(errorRes)
		return
	}
	values := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
	}
	newChirp, err := C.Queries.CreateChirp(context.Background(), values)
	if err != nil {
		res.WriteHeader(400)
		errorRes, _ := json.Marshal(ValidateChirpErrors{Error: fmt.Sprint(err)})
		res.Write(errorRes)
		return
	}
	resBody := CreateChirpRes{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	res.WriteHeader(201)
	newResBody, _ := json.Marshal(resBody)
	res.Write(newResBody)
}
