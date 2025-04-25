package handler

import "net/http"

type UploadFileChunkHandler struct {
}

func NewUploadFileChunkHandler() *UploadFileChunkHandler {
	return &UploadFileChunkHandler{}
}

func (h *UploadFileChunkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
