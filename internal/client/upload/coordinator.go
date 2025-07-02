package upload

import (
	"context"
	"errors"

	us "github.com/grantchen2003/chunky/internal/client/upload/uploadstorer"
)

// figure out the interface pointer thing, what should be pointer and what shouldnt be (why cant i have pointer to interface passed as param)
type Coordinator struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	url                  string
	filePath             string
	maxChunkSizeBytes    int
	maxConcurrentUploads int
	isUploading          bool
	storer               *us.UploadStorer
	validator            *Validator
	requester            *Requester
}

func NewCoordinator(url string, filePath string, maxChunkSizeBytes int, maxConcurrentUploads int, storer us.UploadStorer, validator *Validator, requester *Requester) *Coordinator {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &Coordinator{
		ctx:       ctx,
		ctxCancel: ctxCancel,

		url:                  url,
		filePath:             filePath,
		maxChunkSizeBytes:    maxChunkSizeBytes,
		maxConcurrentUploads: maxConcurrentUploads,
		isUploading:          false,
		storer:               &storer,
		validator:            validator,
		requester:            requester,
	}
}

func (c *Coordinator) ValidateUpload() error {
	if c.isUploading {
		return ErrPausedOnNoOngoingUpload
	}

	return nil
}

func (c *Coordinator) Upload(uploadProgressChan chan<- Progress) Result {
	if err := c.ValidateUpload(); err != nil {
		return UploadResultError
	}

	uploadTask := func(ctx context.Context) error {
		uploader := NewUploader(c.url, c.filePath, c.maxChunkSizeBytes, c.maxConcurrentUploads, uploadProgressChan, *c.storer, c.requester)

		err := uploader.Upload(ctx)

		return err
	}

	return c.runWithUploadLifeCycle(uploadTask)
}

func (c *Coordinator) ValidatePauseUpload() error {
	if !c.isUploading {
		return ErrPausedOnNoOngoingUpload
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
		return ErrResumedOnOngoingUpload
	}

	if !c.validator.hasExistingUpload() {
		return ErrResumedOnNonExistingUpload
	}

	if c.validator.fileHasChangedSinceLastUpload() {
		return ErrResumedOnChangedFile
	}

	return nil
}

func (c *Coordinator) ResumeUpload(uploadProgressChan chan<- Progress) Result {
	if err := c.ValidateResumeUpload(); err != nil {
		return UploadResultError
	}

	resumeUploadTask := func(ctx context.Context) error {
		uploader := NewUploader(c.url, c.filePath, c.maxChunkSizeBytes, c.maxConcurrentUploads, uploadProgressChan, *c.storer, c.requester)

		err := uploader.ResumeUpload(ctx)

		return err
	}

	return c.runWithUploadLifeCycle(resumeUploadTask)
}

func (c *Coordinator) runWithUploadLifeCycle(uploadTask func(ctx context.Context) error) Result {
	c.isUploading = true
	defer func() { c.isUploading = false }()

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	doneChan := make(chan error)
	go func() {
		defer close(doneChan)
		err := uploadTask(c.ctx)
		doneChan <- err
	}()

	err := <-doneChan

	switch {
	case errors.Is(err, context.Canceled):
		return UploadResultPaused

	case err != nil:
		return UploadResultError

	default:
		return UploadResultSuccess
	}
}
