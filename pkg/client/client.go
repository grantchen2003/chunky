package client

import (
	"github.com/grantchen2003/chunky/internal"
	us "github.com/grantchen2003/chunky/internal/uploadstorer"
)

type Client struct {
	uploadNotifier *internal.UploadNotifier
	uploadManager  *internal.UploadManager
}

func NewClient(url string, filePath string) (*Client, error) {
	uploadStorer, err := us.NewSqliteUploadStore()
	if err != nil {
		return nil, err
	}

	uploadValidator := internal.NewUploadValidator(url, filePath, uploadStorer)

	return &Client{
		uploadNotifier: internal.NewUploadNotifier(),
		uploadManager:  internal.NewUploadManager(url, filePath, uploadStorer, uploadValidator),
	}, nil
}

func (c *Client) Upload() error {
	if err := c.uploadManager.ValidateUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadStarted

	uploadResult := c.uploadManager.Upload(c.uploadNotifier.ProgressChan)

	c.uploadNotifier.ResultChan <- uploadResult

	c.uploadNotifier.Close()

	return nil
}

func (c *Client) Pause() error {
	if err := c.uploadManager.ValidatePauseUpload(); err != nil {
		return err
	}

	err := c.uploadManager.PauseUpload()

	return err
}

func (c *Client) Resume() error {
	if err := c.uploadManager.ValidateResumeUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadResumed

	uploadResult := c.uploadManager.ResumeUpload(c.uploadNotifier.ProgressChan)

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
