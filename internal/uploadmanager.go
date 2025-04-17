package internal

import (
	"context"
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

type UploadManager struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	isUploading bool
}

func NewUploadManager() *UploadManager {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &UploadManager{
		ctx:         ctx,
		ctxCancel:   ctxCancel,
		isUploading: false,
	}
}

func (u *UploadManager) PauseUpload() {
	u.ctxCancel()
}

func (u *UploadManager) ValidateUpload() error {
	if u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *UploadManager) ValidatePause() error {
	if !u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *UploadManager) ValidateResume() error {
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

func (u *UploadManager) Upload(url string, filePath FilePath, uploadProgressChan chan<- UploadProgress) UploadResult {
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

func (u *UploadManager) hasNoExistingupload() bool {
	return false
}

func (u *UploadManager) fileHasChangedSinceLastUpload() bool {
	return true
}
