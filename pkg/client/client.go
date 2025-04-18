package client

import (
	"github.com/grantchen2003/chunky/internal"
)

type Client struct {
	filePath string
	url      string

	uploadNotifier *internal.UploadNotifier
	uploadManager  *internal.UploadManager
}

func NewClient(url string, filePath string) *Client {
	return &Client{
		filePath: filePath,
		url:      url,

		uploadNotifier: internal.NewUploadNotifier(),
		uploadManager:  internal.NewUploadManager(),
	}
}

func (c *Client) Upload() error {
	if err := c.uploadManager.ValidateUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadStarted

	uploadResult := c.uploadManager.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

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

	uploadResult := c.uploadManager.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

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
