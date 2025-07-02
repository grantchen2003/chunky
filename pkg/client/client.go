package client

import (
	"github.com/grantchen2003/chunky/internal/client/upload"
	us "github.com/grantchen2003/chunky/internal/client/upload/uploadstorer"
)

type Client struct {
	progressChan chan upload.Progress
	resultChan   chan upload.Result
	statusChan   chan upload.Status
	coordinator  *upload.Coordinator
}

func NewClient(url string, filePath string, uploadEndpoints *upload.Endpoints, maxChunkSizeBytes int, maxConcurrentUploads int) (*Client, error) {
	storer, err := us.NewSqliteUploadStore()
	if err != nil {
		return nil, err
	}

	requester := upload.NewRequester(url, uploadEndpoints)
	validator := upload.NewValidator(url, filePath, storer)

	return &Client{
		progressChan: make(chan upload.Progress),
		resultChan:   make(chan upload.Result),
		statusChan:   make(chan upload.Status),
		coordinator:  upload.NewCoordinator(url, filePath, maxChunkSizeBytes, maxConcurrentUploads, storer, validator, requester),
	}, nil
}

func (c *Client) Upload() error {
	if err := c.coordinator.ValidateUpload(); err != nil {
		return err
	}

	go func() {
		c.statusChan <- upload.UploadStarted
		close(c.statusChan)
	}()

	uploadResult := c.coordinator.Upload(c.progressChan)
	close(c.progressChan)

	c.resultChan <- uploadResult
	close(c.resultChan)

	return nil
}

func (c *Client) Pause() error {
	if err := c.coordinator.ValidatePauseUpload(); err != nil {
		return err
	}

	err := c.coordinator.PauseUpload()

	return err
}

func (c *Client) Resume() error {
	if err := c.coordinator.ValidateResumeUpload(); err != nil {
		return err
	}

	go func() {
		c.statusChan <- upload.UploadResumed
		close(c.statusChan)
	}()

	uploadResult := c.coordinator.ResumeUpload(c.progressChan)
	close(c.progressChan)

	c.resultChan <- uploadResult
	close(c.resultChan)

	return nil
}

func (c *Client) UploadProgressChan() <-chan upload.Progress {
	return c.progressChan
}

func (c *Client) UploadStatusChan() <-chan upload.Status {
	return c.statusChan
}

func (c *Client) UploadResultChan() <-chan upload.Result {
	return c.resultChan
}
