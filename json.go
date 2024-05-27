package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal JSON response: %v", payload)
		writer.WriteHeader(500)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(data)
}
