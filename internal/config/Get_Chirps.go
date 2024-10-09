package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fenroe/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpJSON struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

type ErrorJSON struct {
	Error string `json:"error"`
}

func (C *Config) GetChirps(res http.ResponseWriter, req *http.Request) {
	// Extract author_id and sort parameters
	authorIDString := req.URL.Query().Get("author_id")
	authorID := uuid.NullUUID{}
	if authorIDString != "" {
		id, _ := uuid.Parse(authorIDString)
		authorID.UUID = id
		authorID.Valid = true
	}
	fmt.Println(authorID)
	sort := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	// Decide query based on sort parameter
	if sort == "desc" {
		chirps, err = C.Queries.GetChirpsDesc(context.Background(), authorID)
	} else {
		chirps, err = C.Queries.GetChirps(context.Background(), authorID)
	}

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(400)
		resBody, _ := json.Marshal(ErrorJSON{Error: fmt.Sprint(err)})
		res.Write(resBody)
		return
	}

	// Processing and sending response
	chirpsJSONSlice := []ChirpJSON{}
	for _, value := range chirps {
		c := ChirpJSON{
			ID:        value.ID,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
			Body:      value.Body,
			UserID:    value.UserID,
		}
		chirpsJSONSlice = append(chirpsJSONSlice, c)
	}
	res.WriteHeader(200)
	resBody, _ := json.Marshal(chirpsJSONSlice)
	res.Write(resBody)
}
