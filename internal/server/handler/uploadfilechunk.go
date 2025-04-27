package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal/server/service"
)

// add constraints to make sure payload's chunk size isn't too big (aka can fit all in memory)
// also have to consider when there are many concurrent requests, each goroutine handler's payload
// is valid, but the total memory used by all goroutine's might be too much
type UploadFileChunkHandler struct {
	uploadService *service.UploadService
}

func NewUploadFileChunkHandler(uploadService *service.UploadService) *UploadFileChunkHandler {
	return &UploadFileChunkHandler{
		uploadService: uploadService,
	}
}

func (h *UploadFileChunkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	type Payload struct {
		SessionId string `json:"sessionId"`
		FileHash  []byte `json:"fileHash"`
		Chunk     []byte `json:"chunk"`
		StartByte int    `json:"startByte"`
		EndByte   int    `json:"endByte"`
	}

	var payload Payload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.uploadService.AddFileChunk(
		payload.SessionId,
		payload.FileHash,
		payload.Chunk,
		payload.StartByte,
		payload.EndByte,
	)

	if err != nil {
		http.Error(w, "Failed to upload fike chunk", http.StatusInternalServerError)
		log.Printf("Error uploading file chunk response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
