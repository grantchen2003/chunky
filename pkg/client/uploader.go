package client

import (
	"context"
	"log"
	"time"

	"github.com/grantchen2003/chunky/internal"
)

type UploadPausedError struct{}

func (e *UploadPausedError) Error() string {
	return "Upload paused"
}

func Upload(url string, filePath string, byteRanges []internal.Range, ctx context.Context) error {
	doneChan := make(chan struct{})

	go func() {
		log.Printf("Uploading %s to %s\n", url, filePath)
		time.Sleep(10 * time.Second)
		close(doneChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return &UploadPausedError{}
		case <-doneChan:
			return nil
		}
	}
}
