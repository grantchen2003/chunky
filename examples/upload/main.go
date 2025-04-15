package main

import (
	"fmt"
	"log"

	"github.com/grantchen2003/chunky"
)

func main() {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	client.Upload()

	for {
		select {
		case uploadedBytes := <-client.UploadedBytesChan:
			fmt.Println("UploadedBytes:", uploadedBytes)

		case uploadError := <-client.UploadErrorChan:
			fmt.Println("Error:", uploadError)

		case status := <-client.UploadStatusChan:
			fmt.Println("Status:", status)
			if status == chunky.UploadComplete {
				log.Println("Upload complete!")
				return
			}
		}
	}
}
