package internal

import (
	"context"
)

type UploadManager struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	isUploading bool
}

func NewUploadManager() *UploadManager {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &UploadManager{
		ctx:       ctx,
		ctxCancel: ctxCancel,

		isUploading: false,
	}
}

func (u *UploadManager) ValidateUpload() error {
	if u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *UploadManager) Upload(url string, filePath FilePath, uploadProgressChan chan<- UploadProgress) UploadResult {
	if err := u.ValidateUpload(); err != nil {
		return UploadResultError
	}

	u.isUploading = true
	defer func() { u.isUploading = false }()

	u.ctx, u.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan error)

	go func() {
		defer close(doneChan)
		doneChan <- upload(url, filePath, uploadProgressChan)
	}()

	for {
		select {
		case <-u.ctx.Done():
			return UploadResultPaused

		case err := <-doneChan:
			if err != nil {
				return UploadResultError
			}

			return UploadResultSuccess
		}
	}
}

func (u *UploadManager) ValidatePauseUpload() error {
	if !u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *UploadManager) PauseUpload() error {
	if err := u.ValidatePauseUpload(); err != nil {
		return err
	}

	u.ctxCancel()

	return nil
}

func (u *UploadManager) ValidateResumeUpload() error {
	if u.isUploading {
		return ErrResumedOnOngoingUpload
	}

	if !u.hasExistingupload() {
		return ErrResumedOnNonExistingUpload
	}

	if u.fileHasChangedSinceLastUpload() {
		return ErrResumedOnChangedFile
	}

	return nil
}

func (u *UploadManager) hasExistingupload() bool {
	return true
}

func (u *UploadManager) fileHasChangedSinceLastUpload() bool {
	return true
}
