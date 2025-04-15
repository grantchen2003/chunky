package main

import (
	"fmt"
	"time"

	"github.com/grantchen2003/chunky"
)

func main() {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	// start upload
	go client.Upload()

	// pause upload after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		client.Pause()
	}()

	// resume upload after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		client.Resume()
	}()

	for {
		select {
		case uploadedBytes := <-client.UploadedBytesChan:
			fmt.Println("UploadedBytes:", uploadedBytes)

		case uploadError := <-client.UploadErrorChan:
			fmt.Println("Error:", uploadError)

		case status := <-client.UploadStatusChan:
			fmt.Println("Status:", status)
			if status == chunky.UploadComplete {
				return
			}
		}
	}
}
