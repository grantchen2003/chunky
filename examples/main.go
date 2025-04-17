package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			go client.Upload()

		case "pause":
			go client.Pause()

		case "resume":
			go client.Resume()

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

func main() {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	go handleUserCommands(client)

	for {
		select {
		case status := <-client.UploadStatusChan:
			fmt.Println("Status:", status)

			if status == chunky.UploadCompleted {
				return
			}

		case uploadProgress := <-client.UploadProgressChan:
			fmt.Println("Upload progress:", uploadProgress)

		case uploadError := <-client.UploadErrorChan:
			fmt.Println("Error:", uploadError)
		}
	}
}
