package main

import "github.com/grantchen2003/chunky"

func main () {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	// Start the upload
	go func() {
		if err := client.Upload(); err != nil {
			log.Println("Error during upload:", err)
			return
		}

		log.Println("Upload complete!")
	}()

	// Simulate user pause after 1 second
	time.Sleep(1 * time.Second)
	client.Pause()

	// Simulate user resume after 2 seconds
	time.Sleep(2 * time.Second)
	client.Resume()
}