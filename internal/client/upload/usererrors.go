package upload

import "errors"

var (
	ErrStartedOnOngoingUpload     = errors.New("upload start when an upload is already ongoing")
	ErrPausedOnNoOngoingUpload    = errors.New("paused when no upload is ongoing")
	ErrResumedOnNonExistingUpload = errors.New("resumed on non existing upload error")
	ErrResumedOnOngoingUpload     = errors.New("resumed on ongoing upload error")
	ErrResumedOnChangedFile       = errors.New("resumed on changed file")
)
