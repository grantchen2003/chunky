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

func (u *Uploader) ValidateUpload() error {
	if !u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *Uploader) ValidatePause() error {
	if !u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *Uploader) ValidateResume() error {
	if u.isUploading {
		return ErrResumedOnOngoingUpload
	}

	if u.hasNoExistingupload() {
		return ErrResumedOnNonExistingUpload
	}

	if u.fileHasChangedSinceLastUpload() {
		return ErrResumedOnChangedFile
	}

	return nil
}

func (u *Uploader) Upload(ctx context.Context, url string, filePath FilePath, uploadProgressChan chan<- UploadProgress) UploadResult {
	u.isUploading = true
	defer func() { u.isUploading = false }()

	doneChan := make(chan struct{})

	errorOccured := true

	go func() {
		log.Printf("Uploading %s to %s\n", url, filePath)
		for i := range 3 {
			time.Sleep(1 * time.Second)
			uploadProgressChan <- UploadProgress{
				PercentageUploaded: 100 * i / 3,
			}
		}
		// simulate file upload
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

func (u *Uploader) hasNoExistingupload() bool {
	return false
}

func (u *Uploader) fileHasChangedSinceLastUpload() bool {
	return true
}
