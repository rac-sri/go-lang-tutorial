package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON( w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload);

	if err!=nil {
		log.Println("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code>499 {
		log.Println("Responding with 5XX error:",msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	// {
	// 	error: somethign went wrong
	// }

	respondWithJSON(w, code, errResponse{
		Error: msg,
	});
}