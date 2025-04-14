package main

import "github.com/grantchen2003/chunky"

func main () {
	client := chunky.NewClient("http://localhost:8080", "bigfile.txt")

	if err := client.Upload(); err != nil {
		log.Println("Error during upload:", err)
		return
	}
	
	log.Println("Upload complete!")
}