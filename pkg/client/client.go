package client

import (
	"github.com/grantchen2003/chunky/internal/upload"
	us "github.com/grantchen2003/chunky/internal/upload/uploadstorer"
)

type Client struct {
	notifier    *upload.Notifier
	coordinator *upload.Coordinator
}

func NewClient(url string, filePath string, uploadEndpoints *upload.Endpoints) (*Client, error) {
	storer, err := us.NewSqliteUploadStore()
	if err != nil {
		return nil, err
	}

	requester := upload.NewRequester(url, uploadEndpoints)
	validator := upload.NewValidator(url, filePath, storer)

	return &Client{
		notifier:    upload.NewNotifier(),
		coordinator: upload.NewCoordinator(url, filePath, storer, validator, requester),
	}, nil
}

func (c *Client) Upload() error {
	if err := c.coordinator.ValidateUpload(); err != nil {
		return err
	}

	c.notifier.StatusChan <- upload.UploadStarted

	uploadResult := c.coordinator.Upload(c.notifier.ProgressChan)

	c.notifier.ResultChan <- uploadResult

	c.notifier.Close()

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

	c.notifier.StatusChan <- upload.UploadResumed

	uploadResult := c.coordinator.ResumeUpload(c.notifier.ProgressChan)

	c.notifier.ResultChan <- uploadResult

	c.notifier.Close()

	return nil
}

func (c *Client) UploadProgressChan() <-chan upload.Progress {
	return c.notifier.ProgressChan
}

func (c *Client) UploadStatusChan() <-chan upload.Status {
	return c.notifier.StatusChan
}

func (c *Client) UploadResultChan() <-chan upload.Result {
	return c.notifier.ResultChan
}
