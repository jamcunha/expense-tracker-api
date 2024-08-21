package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("error writing JSON response: %v", err)
	}
}
