package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal/server/service"
)

type InitiateUploadSessionHandler struct {
	uploadService *service.UploadService
}

func NewInitiateUploadSessionHandler(uploadService *service.UploadService) *InitiateUploadSessionHandler {
	return &InitiateUploadSessionHandler{
		uploadService: uploadService,
	}
}

func (h *InitiateUploadSessionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	type Payload struct {
		FileHash           []byte `json:"fileHash"`
		TotalFileSizeBytes int    `json:"TotalFileSizeBytes"`
	}

	var payload Payload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sessionId, err := h.uploadService.CreateUploadSession(payload.FileHash, payload.TotalFileSizeBytes)
	if err != nil {
		http.Error(w, "Failed to create upload session", http.StatusInternalServerError)
		log.Printf("Error creating uploading session response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{"sessionId": sessionId}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
