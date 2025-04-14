package internal

import (
	"context"

	us "github.com/grantchen2003/chunky/internal/uploadstorer"
)

// NEED TO REFACTOR AND PASS CTX DOWN TO PREVENT GOROUTINE LEAKS
// figure out the interface pointer thing, what should be pointer and what shouldnt be (why cant i have pointer to interface passed as param)
type UploadManager struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	url             string
	filePath        string
	isUploading     bool
	uploadStorer    *us.UploadStorer
	uploadValidator *UploadValidator
}

func NewUploadManager(url string, filePath string, uploadStorer us.UploadStorer, uploadValidator *UploadValidator) *UploadManager {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &UploadManager{
		ctx:       ctx,
		ctxCancel: ctxCancel,

		url:             url,
		filePath:        filePath,
		isUploading:     false,
		uploadStorer:    &uploadStorer,
		uploadValidator: uploadValidator,
	}
}

func (u *UploadManager) ValidateUpload() error {
	if u.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (u *UploadManager) Upload(uploadProgressChan chan<- UploadProgress) UploadResult {
	if err := u.ValidateUpload(); err != nil {
		return UploadResultError
	}

	uploadTask := func() error {
		uploader := NewUploader(u.url, u.filePath, uploadProgressChan, *u.uploadStorer)

		err := uploader.Upload()

		return err
	}

	return u.runWithUploadLifeCycle(uploadTask)
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

	if !u.uploadValidator.hasExistingUpload() {
		return ErrResumedOnNonExistingUpload
	}

	if u.uploadValidator.fileHasChangedSinceLastUpload() {
		return ErrResumedOnChangedFile
	}

	return nil
}

func (u *UploadManager) ResumeUpload(uploadProgressChan chan<- UploadProgress) UploadResult {
	if err := u.ValidateResumeUpload(); err != nil {
		return UploadResultError
	}

	resumeUploadTask := func() error {
		uploader := NewUploader(u.url, u.filePath, uploadProgressChan, *u.uploadStorer)

		err := uploader.ResumeUpload()

		return err
	}

	return u.runWithUploadLifeCycle(resumeUploadTask)
}

func (u *UploadManager) runWithUploadLifeCycle(uploadTask func() error) UploadResult {
	u.isUploading = true
	defer func() { u.isUploading = false }()

	u.ctx, u.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan error)
	go func() {
		defer close(doneChan)
		err := uploadTask()
		doneChan <- err
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
