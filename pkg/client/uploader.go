package client

import (
	"context"
	"fmt"
	"time"
)

func Upload(url string, filePath string, ctx context.Context) error {
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
			return fmt.Errorf("context cancelled")
		case <-doneChan:
			// log.Printf("done uploading")
			return nil
		}
	}
}
