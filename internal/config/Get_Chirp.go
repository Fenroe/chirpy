package config

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (C *Config) GetChirp(res http.ResponseWriter, req *http.Request) {
	chirpID, _ := uuid.Parse(req.PathValue("chirpID"))
	chirp, err := C.Queries.GetChirp(context.Background(), chirpID)
	if err != nil {
		res.WriteHeader(404)
		body, _ := json.Marshal(ErrorJSON{err.Error()})
		res.Write(body)
		return
	}
	res.WriteHeader(200)
	body, _ := json.Marshal(ChirpJSON{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
	res.Write(body)
}
