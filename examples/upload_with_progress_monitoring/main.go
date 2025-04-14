package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grantchen2003/chunky"
)

func main() {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	// Start the upload
	go func() {
		if err := client.Upload(); err != nil {
			log.Println("Error during upload:", err)
			return
		}

		log.Println("Upload complete!")
	}()

	// Monitor upload progress (optional)
	go func() {
		for {
			select {
			case uploadedBytes := <-client.UploadedBytesChan():
				fmt.Println(uploadedBytes)

			case uploadError := <-client.UploadErrorChan():
				fmt.Println("Error:", uploadError)

			case <-client.UploadStatusChan():
				if client.UploadIsCompleted() {
					fmt.Println("Upload is in complete")
					return

				} else if client.UploadIsInProgress() {
					fmt.Println("Upload is in progress")

				} else if client.UploadIsPaused() {
					fmt.Println("Upload is in paused")

				}
			}
		}
	}()

	// Simulate user pause after 1 seconds
	time.Sleep(1 * time.Second)
	client.Pause()

	// Simulate user resume after 2 seconds
	time.Sleep(2 * time.Second)
	if err := client.Upload(); err != nil {
		log.Println("Error during upload:", err)
	}

	log.Println("Upload complete!")
}
