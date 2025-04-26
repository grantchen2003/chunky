package main

import (
	"fmt"
	"log"

	"github.com/grantchen2003/chunky"
)

func main() {
	port := ":8080"
	server, err := chunky.NewServer(port)
	if err != nil {
		panic(err)
	}

	server.SetInitiateUploadSessionEndpoint("/my-custom-initiate-upload-session-endpoint")
	server.SetByteRangesToUploadEndpoint("/my-custom-byte-ranges-to-upload-endpoint")
	server.SetUploadFileChunkEndpoint("/my-custom-upload-file-chunk-endpoint")

	fmt.Printf("Server started on port %s\n", port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
