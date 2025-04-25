package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal"
)

type UploadFileChunkHandler struct {
	uploadSessionService *internal.UploadSessionService
}

func NewUploadFileChunkHandler(uploadSessionService *internal.UploadSessionService) *UploadFileChunkHandler {
	return &UploadFileChunkHandler{
		uploadSessionService: uploadSessionService,
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

	err := h.uploadSessionService.AddFileChunk(
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
