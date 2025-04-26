package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal"
)

type ByteRangesToUploadHandler struct {
	uploadService *internal.UploadService
}

func NewByteRangesToUploadHandler(uploadService *internal.UploadService) *ByteRangesToUploadHandler {
	return &ByteRangesToUploadHandler{
		uploadService: uploadService,
	}
}

func (h *ByteRangesToUploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	type Payload struct {
		SessionId string `json:"sessionId"`
		FileHash  []byte `json:"fileHash"`
	}

	var payload Payload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	byteRangesToUpload, err := h.uploadService.ByteRangesToUpload(payload.SessionId, payload.FileHash)
	if err != nil {
		http.Error(w, "Failed to create upload session", http.StatusInternalServerError)
		log.Printf("Error creating uploading session response: %v", err)
		return
	}

	responseData := map[string][][2]int{
		"ByteRanges": byteRangesToUpload,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
