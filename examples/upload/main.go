package main

import (
	"fmt"
	"log"

	"github.com/grantchen2003/chunky"
)

func main() {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	go client.Upload()

	for {
		select {
		case progress := <-client.ProgressChan:
			fmt.Println("Progress:", progress)

		case uploadError := <-client.UploadErrorChan:
			fmt.Println("Error:", uploadError)

		case status := <-client.UploadStatusChan:
			fmt.Println("Status:", status)
			if status == chunky.UploadCompleted {
				log.Println("Upload complete!")
				return
			}
		}
	}
}
