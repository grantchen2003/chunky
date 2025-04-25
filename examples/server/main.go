package main

import (
	"fmt"
	"log"

	"github.com/grantchen2003/chunky"
)

func main() {
	server, err := chunky.NewServer(":8080")
	if err != nil {
		panic(err)
	}

	server.SetInitiateUploadSessionEndpoint("/my-custom-initiate-upload-session-endpoint")
	server.SetByteRangesToUploadEndpoint("/my-custom-byte-ranges-to-upload-endpoint")
	server.SetUploadFileChunkEndpoint("/my-custom-upload-file-chunk-endpoint")

	fmt.Println("Server started on port :8080")
	log.Fatal(server.Start())
}
