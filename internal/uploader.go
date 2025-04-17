package internal

import (
	"context"
	"fmt"
	"log"
	"time"
)

type UploadResult int

type FilePath string

const (
	UploadResultUnknown UploadResult = iota
	UploadResultSuccess
	UploadResultPaused
	UploadResultError
)

type Uploader struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	isUploading bool
}

func NewUploader() *Uploader {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &Uploader{
		ctx:         ctx,
		ctxCancel:   ctxCancel,
		isUploading: false,
	}
}

func (u *Uploader) PauseUpload() {
	u.ctxCancel()
}

func (u *Uploader) ValidateUpload() error {
	if u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *Uploader) ValidatePause() error {
	fmt.Println(u)
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

func (u *Uploader) Upload(url string, filePath FilePath, uploadProgressChan chan<- UploadProgress) UploadResult {
	u.isUploading = true
	defer func() { u.isUploading = false }()

	u.ctx, u.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan struct{})

	errorOccured := true

	go func() {
		log.Printf("Uploading %s to %s\n", url, filePath)
		for i := range 5 {
			time.Sleep(5 * time.Second)
			uploadProgressChan <- UploadProgress{
				PercentageUploaded: 100 * i / 5,
			}
		}
		// simulate file upload
		close(doneChan)
	}()

	for {
		select {
		case <-u.ctx.Done():
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
