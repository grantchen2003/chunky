package internal

import (
	"bytes"
	"context"
	"fmt"
)

// NEED TO REFACTOR
type UploadManager struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	url         string
	filePath    string
	isUploading bool
}

func NewUploadManager(url string, filePath string) *UploadManager {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &UploadManager{
		ctx:       ctx,
		ctxCancel: ctxCancel,

		url:         url,
		filePath:    filePath,
		isUploading: false,
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

	u.isUploading = true
	defer func() { u.isUploading = false }()

	u.ctx, u.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan error)

	go func() {
		defer close(doneChan)

		err := func() error {
			sqliteSessionStore, err := NewSqliteUploadSessionStore()
			if err != nil {
				return err
			}
			defer sqliteSessionStore.Close()

			uploader := NewUploader(u.url, u.filePath, uploadProgressChan, sqliteSessionStore)
			err = uploader.Upload()
			return err
		}()

		doneChan <- err
	}()

	for {
		select {
		case <-u.ctx.Done():
			return UploadResultPaused

		case err := <-doneChan:
			if err != nil {
				fmt.Println(err)
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

	if !u.hasExistingUpload() {
		return ErrResumedOnNonExistingUpload
	}

	if u.fileHasChangedSinceLastUpload() {
		return ErrResumedOnChangedFile
	}

	return nil
}

func (u *UploadManager) hasExistingUpload() bool {
	sqliteSessionStore, err := NewSqliteUploadSessionStore()
	if err != nil {
		return false // default to false if program fuck up
	}
	defer sqliteSessionStore.Close()

	_, _, err = sqliteSessionStore.GetSessionIdAndHash(u.filePath, u.url)

	return err != nil
}

func (u *UploadManager) fileHasChangedSinceLastUpload() bool {
	sqliteSessionStore, err := NewSqliteUploadSessionStore()
	if err != nil {
		return true // default to true if program fuck up
	}
	defer sqliteSessionStore.Close()

	_, savedFileHash, err := sqliteSessionStore.GetSessionIdAndHash(u.filePath, u.url)
	if err != nil {
		return true
	}

	currFileHash, err := hashFile(u.filePath)
	if err != nil {
		return true
	}

	if !bytes.Equal(currFileHash, savedFileHash) {
		return true
	}

	return false
}
