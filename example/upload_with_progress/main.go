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

	// Monitor upload progress (optional)
	go func() {
		for {
			select {
			case progress := <-client.ProgressChan:
				fmt.Println(progress)
			case err := <-client.ErrorChan:
				fmt.Println("Error:", err)
			case status := <-client.StatusChan:
				fmt.Println("Status:", status)
				if status == client.Complete {
					return
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