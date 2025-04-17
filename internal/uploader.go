package internal

import (
	"context"
	"log"
	"time"
)

type UploadResult int

type Uploader struct {
	isUploading bool
}

const (
	UploadResultUnknown UploadResult = iota
	UploadResultSuccess
	UploadResultPaused
	UploadResultError
)

func NewUploader() *Uploader {
	return &Uploader{}
}

func (u *Uploader) IsUploading() bool {
	return u.isUploading
}

func (u *Uploader) HasNoExistingupload() bool {
	return false
}

func (u *Uploader) Upload(ctx context.Context, url string, filePath FilePath) UploadResult {
	u.isUploading = true
	defer func() { u.isUploading = false }()

	doneChan := make(chan struct{})

	errorOccured := true

	go func() {
		log.Printf("Uploading %s to %s\n", url, filePath)
		time.Sleep(3 * time.Second) // simulate file upload
		close(doneChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return UploadResultPaused
		case <-doneChan:
			if errorOccured {
				return UploadResultError
			} else {
				return UploadResultSuccess
			}
		}
	}
}
