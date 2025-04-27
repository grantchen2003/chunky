package upload

import (
	"bytes"

	"github.com/grantchen2003/chunky/internal/client/file"
	us "github.com/grantchen2003/chunky/internal/client/upload/uploadstorer"
)

type Validator struct {
	url          string
	filePath     string
	uploadStorer us.UploadStorer
}

func NewValidator(url string, filePath string, uploadStorer us.UploadStorer) *Validator {
	return &Validator{
		url:          url,
		filePath:     filePath,
		uploadStorer: uploadStorer,
	}
}

func (v Validator) hasExistingUpload() bool {
	uploadExists, err := func() (bool, error) {
		_, _, err := v.uploadStorer.GetSessionIdAndFileHash(v.url, v.filePath)
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

func (v Validator) fileHasChangedSinceLastUpload() bool {
	hasChanged, err := func() (bool, error) {
		_, savedFileHash, err := v.uploadStorer.GetSessionIdAndFileHash(v.url, v.filePath)
		if err != nil {
			return false, err
		}

		currFileHash, err := file.HashFile(v.filePath)
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
