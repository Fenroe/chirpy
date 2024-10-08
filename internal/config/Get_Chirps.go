package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ChirpJSON struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ErrorJSON struct {
	Error string `json:"error"`
}

func (C *Config) GetChirps(res http.ResponseWriter, _ *http.Request) {
	chirps, err := C.Queries.GetChirps(context.Background())
	if err != nil {
		res.WriteHeader(400)
		resBody, _ := json.Marshal(ErrorJSON{Error: fmt.Sprint(err)})
		res.Write(resBody)
		return
	}
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
