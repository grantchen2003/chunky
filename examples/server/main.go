package main

import (
	"net/http"

	"github.com/grantchen2003/chunky"
)

func main() {
	chunky.StartServer(
		":8080",
		map[string]func(w http.ResponseWriter, r *http.Request){
			"/my-custom-initiate-upload-session-endpoint": chunky.HandleInitiateUploadSession,
			"/my-custom-byte-ranges-to-upload-endpoint":   chunky.HandleByteRangesToUpload,
			"/my-custom-upload-file-chunk-endpoint":       chunky.HandleUploadFileChunk,
		},
	)
}
