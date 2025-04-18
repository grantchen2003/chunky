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
			go func() {
				if err := client.Upload(); err != nil {
					fmt.Printf("User error: %v\n", err)
				}
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
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	go handleUserCommands(client)

	for {
		select {
		case result := <-client.UploadResultChan():
			fmt.Println("Result:", result)
			return

		case status := <-client.UploadStatusChan():
			fmt.Println("Status:", status.Message)

		case uploadProgress := <-client.UploadProgressChan():
			fmt.Println("Upload progress:", uploadProgress)
		}
	}
}
