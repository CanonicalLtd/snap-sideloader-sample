package service

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONHeader is the header for JSON responses
const JSONHeader = "application/json; charset=UTF-8"

// StandardResponse is the JSON response from the web service
type StandardResponse struct {
	StatusCode int         `json:"status-code"`
	Status     string      `json:"status"`
	Type       string      `json:"type"`
	Result     interface{} `json:"result,omitempty"`
}

// formatStandardResponse returns a JSON response from an API method, indicating success or failure
func formatStandardResponse(code, message string, w http.ResponseWriter) {
	var response StandardResponse
	w.Header().Set("Content-Type", JSONHeader)

	if len(code) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		response = StandardResponse{Type: "sync", Status: "error", StatusCode: 400, Result: message}
	} else {
		response = StandardResponse{Type: "sync", Status: "OK", StatusCode: 200, Result: message}
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}

func encodeResponse(w http.ResponseWriter, response interface{}) {
	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error forming the response:", err)
	}
}
