package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal"
)

type InitiateUploadSessionHandler struct {
	uploadSessionService *internal.UploadSessionService
}

func NewInitiateUploadSessionHandler(uploadSessionService *internal.UploadSessionService) *InitiateUploadSessionHandler {
	fmt.Println("bro", uploadSessionService)

	return &InitiateUploadSessionHandler{
		uploadSessionService: uploadSessionService,
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

	fmt.Println("yooo", h.uploadSessionService)

	sessionId, err := h.uploadSessionService.CreateUploadSession(payload.FileHash, payload.TotalFileSizeBytes)
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
