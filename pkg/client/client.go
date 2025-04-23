package client

import (
	"github.com/grantchen2003/chunky/internal"
	us "github.com/grantchen2003/chunky/internal/uploadstorer"
)

type Client struct {
	uploadNotifier    *internal.UploadNotifier
	uploadCoordinator *internal.UploadCoordinator
}

func NewClient(url string, filePath string) (*Client, error) {
	uploadStorer, err := us.NewSqliteUploadStore()
	if err != nil {
		return nil, err
	}

	uploadValidator := internal.NewUploadValidator(url, filePath, uploadStorer)

	return &Client{
		uploadNotifier:    internal.NewUploadNotifier(),
		uploadCoordinator: internal.NewUploadCoordinator(url, filePath, uploadStorer, uploadValidator),
	}, nil
}

func (c *Client) Upload() error {
	if err := c.uploadCoordinator.ValidateUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadStarted

	uploadResult := c.uploadCoordinator.Upload(c.uploadNotifier.ProgressChan)

	c.uploadNotifier.ResultChan <- uploadResult

	c.uploadNotifier.Close()

	return nil
}

func (c *Client) Pause() error {
	if err := c.uploadCoordinator.ValidatePauseUpload(); err != nil {
		return err
	}

	err := c.uploadCoordinator.PauseUpload()

	return err
}

func (c *Client) Resume() error {
	if err := c.uploadCoordinator.ValidateResumeUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadResumed

	uploadResult := c.uploadCoordinator.ResumeUpload(c.uploadNotifier.ProgressChan)

	c.uploadNotifier.ResultChan <- uploadResult

	c.uploadNotifier.Close()

	return nil
}

func (c *Client) UploadProgressChan() <-chan internal.UploadProgress {
	return c.uploadNotifier.ProgressChan
}

func (c *Client) UploadStatusChan() <-chan internal.UploadStatus {
	return c.uploadNotifier.StatusChan
}

func (c *Client) UploadResultChan() <-chan internal.UploadResult {
	return c.uploadNotifier.ResultChan
}
