package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ByteRangesToUploadHandler struct {
}

func NewByteRangesToUploadHandler() *ByteRangesToUploadHandler {
	return &ByteRangesToUploadHandler{}
}

func (h *ByteRangesToUploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	responseData := map[string][][2]int{
		"ByteRanges": {{6, 100}, {102, 132}, {103, 104}, {133, 152}, {154, 154}, {155, 155}, {156, 159}, {161, 306}},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
