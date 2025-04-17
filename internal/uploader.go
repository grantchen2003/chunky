package internal

import (
	"log"
	"time"
)

// simulate file upload
func upload(url string, filePath FilePath, uploadProgressChan chan<- UploadProgress) error {
	log.Printf("Uploading %s to %s\n", url, filePath)
	for i := range 5 {
		time.Sleep(5 * time.Second)
		uploadProgressChan <- UploadProgress{
			PercentageUploaded: 100 * i / 5,
		}
	}
	return nil
}
