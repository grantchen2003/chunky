package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grantchen2003/chunky"
)

func handleUserCommands(client *chunky.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter a command (upload, pause, resume):")

	for {
		input, _ := reader.ReadString('\n')

		input = strings.TrimSpace(input)

		switch input {
		case "upload":

			start := time.Now()
			go func() {
				if err := client.Upload(); err != nil {
					fmt.Printf("User error: %v\n", err)
				}
				fmt.Printf("Upload took %s\n", time.Since(start))
			}()

		case "pause":
			go func() {
				if err := client.Pause(); err != nil {
					fmt.Printf("User error: %v\n", err)
				}
			}()

		case "resume":
			go func() {
				if err := client.Resume(); err != nil {
					fmt.Printf("User error: %v\n", err)
				}
			}()

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

func main() {
	client, err := chunky.NewClient(
		"http://localhost:8080",
		"file.txt",
		&chunky.UploadEndpoints{
			InitiateUploadSession: "/my-custom-initiate-upload-session-endpoint",
			ByteRangesToUpload:    "/my-custom-byte-ranges-to-upload-endpoint",
			UploadFileChunk:       "/my-custom-upload-file-chunk-endpoint",
		},
	)

	if err != nil {
		panic(err)
	}

	go handleUserCommands(client)

	var totalUploadedBytes int

	for {
		select {
		case result := <-client.UploadResultChan():
			fmt.Println("Result:", result)
			return

		case status := <-client.UploadStatusChan():
			fmt.Println("Status:", status.Message)

		case uploadProgress := <-client.UploadProgressChan():
			totalUploadedBytes += uploadProgress.UploadedBytes
			fmt.Println("Upload progress:", 100*totalUploadedBytes/uploadProgress.TotalBytesToUpload)
		}
	}
}
