package config

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fenroe/chirpy/internal/auth"
	"github.com/google/uuid"
)

type upgradeUserReq struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (C *Config) UpgradeUserWebhook(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		res.WriteHeader(401)
		return
	}
	if apiKey != C.PulkaKey {
		res.WriteHeader(401)
		return
	}
	reqBody := upgradeUserReq{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	if reqBody.Event != "user.upgraded" {
		res.WriteHeader(204)
		return
	}
	_, err = C.Queries.UpgradeUser(context.Background(), reqBody.Data.UserID)
	if err != nil {
		res.WriteHeader(404)
		return
	}
	res.WriteHeader(204)
}
