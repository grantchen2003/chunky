package upload

import (
	"context"

	"github.com/grantchen2003/chunky/internal"
	us "github.com/grantchen2003/chunky/internal/upload/uploadstorer"
)

// NEED TO REFACTOR AND PASS CTX DOWN TO PREVENT GOROUTINE LEAKS
// figure out the interface pointer thing, what should be pointer and what shouldnt be (why cant i have pointer to interface passed as param)
type Coordinator struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	url         string
	filePath    string
	isUploading bool
	storer      *us.UploadStorer
	validator   *Validator
	requester   *Requester
}

func NewCoordinator(url string, filePath string, storer us.UploadStorer, validator *Validator, requester *Requester) *Coordinator {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &Coordinator{
		ctx:       ctx,
		ctxCancel: ctxCancel,

		url:         url,
		filePath:    filePath,
		isUploading: false,
		storer:      &storer,
		validator:   validator,
		requester:   requester,
	}
}

func (c *Coordinator) ValidateUpload() error {
	if c.isUploading {
		return internal.ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (c *Coordinator) Upload(uploadProgressChan chan<- Progress) Result {
	if err := c.ValidateUpload(); err != nil {
		return UploadResultError
	}

	uploadTask := func() error {
		uploader := NewUploader(c.url, c.filePath, uploadProgressChan, *c.storer, c.requester)

		err := uploader.Upload()

		return err
	}

	return c.runWithUploadLifeCycle(uploadTask)
}

func (c *Coordinator) ValidatePauseUpload() error {
	if !c.isUploading {
		return internal.ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (c *Coordinator) PauseUpload() error {
	if err := c.ValidatePauseUpload(); err != nil {
		return err
	}

	c.ctxCancel()

	return nil
}

func (c *Coordinator) ValidateResumeUpload() error {
	if c.isUploading {
		return internal.ErrResumedOnOngoingUpload
	}

	if !c.validator.hasExistingUpload() {
		return internal.ErrResumedOnNonExistingUpload
	}

	if c.validator.fileHasChangedSinceLastUpload() {
		return internal.ErrResumedOnChangedFile
	}

	return nil
}

func (c *Coordinator) ResumeUpload(uploadProgressChan chan<- Progress) Result {
	if err := c.ValidateResumeUpload(); err != nil {
		return UploadResultError
	}

	resumeUploadTask := func() error {
		uploader := NewUploader(c.url, c.filePath, uploadProgressChan, *c.storer, c.requester)

		err := uploader.ResumeUpload()

		return err
	}

	return c.runWithUploadLifeCycle(resumeUploadTask)
}

func (c *Coordinator) runWithUploadLifeCycle(uploadTask func() error) Result {
	c.isUploading = true
	defer func() { c.isUploading = false }()

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan error)
	go func() {
		defer close(doneChan)
		err := uploadTask()
		doneChan <- err
	}()

	for {
		select {
		case <-c.ctx.Done():
			return UploadResultPaused

		case err := <-doneChan:
			if err != nil {
				return UploadResultError
			}

			return UploadResultSuccess
		}
	}
}
