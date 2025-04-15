package client

import (
	"context"
	"time"

	"github.com/grantchen2003/chunky/internal"
)

type UploadCancelledByPauseError struct {
}

func (e *UploadCancelledByPauseError) Error() string {
	return "Upload cancelled by pause"
}

func Upload(url string, filePath string, byteRanges []internal.Range, ctx context.Context) error {
	doneChan := make(chan struct{})

	go func() {
		// log.Printf("Uploading %s to %s\n", url, filePath)
		time.Sleep(10 * time.Second)
		close(doneChan)
	}()

	for {
		select {
		case <-ctx.Done():
			// log.Printf("context cancelled")
			return &UploadCancelledByPauseError{}
		case <-doneChan:
			// log.Printf("done uploading")
			return nil
		}
	}
}

func ResumeUpload() {

}
