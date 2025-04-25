package main

import (
	"fmt"
	"log"

	"github.com/grantchen2003/chunky"
)

func main() {
	server := chunky.NewServer(":8080")

	server.SetInitiateUploadSessionEndpoint("/my-custom-initiate-upload-session-endpoint")
	server.SetByteRangesToUploadEndpoint("/my-custom-byte-ranges-to-upload-endpoint")

	fmt.Println("Server started on port :8080")
	log.Fatal(server.Start())
}
