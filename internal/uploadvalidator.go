package internal

import (
	"bytes"

	us "github.com/grantchen2003/chunky/internal/uploadstorer"
)

type UploadValidator struct {
	url          string
	filePath     string
	uploadStorer us.UploadStorer
}

func NewUploadValidator(url string, filePath string, uploadStorer us.UploadStorer) *UploadValidator {
	return &UploadValidator{
		url:          url,
		filePath:     filePath,
		uploadStorer: uploadStorer,
	}
}

func (uv UploadValidator) hasExistingUpload() bool {
	uploadExists, err := func() (bool, error) {
		_, _, err := uv.uploadStorer.GetSessionIdAndFileHash(uv.url, uv.filePath)
		if err != nil {
			if err == us.ErrNotFound {
				return false, nil
			}
			return false, err
		}

		return true, nil
	}()

	if err != nil {
		return false
	}

	return uploadExists
}

func (uv UploadValidator) fileHasChangedSinceLastUpload() bool {
	hasChanged, err := func() (bool, error) {
		_, savedFileHash, err := uv.uploadStorer.GetSessionIdAndFileHash(uv.url, uv.filePath)
		if err != nil {
			return false, err
		}

		currFileHash, err := HashFile(uv.filePath)
		if err != nil {
			return false, err
		}

		return !bytes.Equal(currFileHash, savedFileHash), nil
	}()

	if err != nil {
		return true
	}

	return hasChanged
}
