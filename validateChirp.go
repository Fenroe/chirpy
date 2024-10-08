package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type reqVals struct {
	Body string `json:"body"`
}

type errorVals struct {
	Error error `json:"error"`
}

type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirp(res http.ResponseWriter, req *http.Request) {
	defer res.Header().Set("Content-Type", "application/json")
	body := reqVals{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		errorRes, _ := json.Marshal(errorVals{Error: err})
		res.Write(errorRes)
		return
	}
	isValid := len(body.Body) <= 140
	if !isValid {
		errorJson := errorVals{
			Error: errors.New("chirp is too long"),
		}
		responseJson, _ := json.Marshal(errorJson)
		res.WriteHeader(400)
		res.Write(responseJson)
	} else {
		words := strings.Fields(body.Body)
		for i, value := range words {
			lower := strings.ToLower(value)
			if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
				words[i] = "****"
			}
		}
		cleanedBody := strings.Join(words, " ")
		validJson := returnVals{
			CleanedBody: cleanedBody,
		}
		responseJson, _ := json.Marshal(validJson)
		res.WriteHeader(200)
		res.Write(responseJson)
	}
}
