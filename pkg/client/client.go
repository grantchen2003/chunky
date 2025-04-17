package client

import (
	"github.com/grantchen2003/chunky/internal"
)

type Client struct {
	filePath internal.FilePath
	url      string

	uploadNotifier *internal.UploadNotifier
	uploadManager  *internal.UploadManager
}

func NewClient(url string, filePathStr string) *Client {
	return &Client{
		filePath: internal.FilePath(filePathStr),
		url:      url,

		uploadNotifier: internal.NewUploadNotifier(),
		uploadManager:  internal.NewUploadManager(),
	}
}

func (c *Client) Upload() {
	if err := c.uploadManager.ValidateUpload(); err != nil {
		c.uploadNotifier.UserErrorChan <- err
	}

	c.upload(internal.UploadStarted)
}

func (c *Client) Pause() {
	if err := c.uploadManager.ValidatePause(); err != nil {
		c.uploadNotifier.UserErrorChan <- err
	}

	c.uploadManager.PauseUpload()
}

func (c *Client) Resume() {
	if err := c.uploadManager.ValidateResume(); err != nil {
		c.uploadNotifier.UserErrorChan <- err
	}

	c.upload(internal.UploadResumed)
}

func (c *Client) UploadProgressChan() <-chan internal.UploadProgress {
	return c.uploadNotifier.ProgressChan
}

func (c *Client) UploadErrorChan() <-chan error {
	return c.uploadNotifier.ErrorChan
}

func (c *Client) UploadStatusChan() <-chan internal.UploadStatus {
	return c.uploadNotifier.StatusChan
}

func (c *Client) UserErrorChan() <-chan error {
	return c.uploadNotifier.UserErrorChan
}

func (c *Client) upload(uploadStatus internal.UploadStatus) {
	c.uploadNotifier.StatusChan <- uploadStatus

	uploadResult := c.uploadManager.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

	c.uploadNotifier.StatusChan <- determineUploadStatus(uploadResult)

	c.uploadNotifier.Close()
}

func determineUploadStatus(uploadResult internal.UploadResult) internal.UploadStatus {
	switch uploadResult {
	case internal.UploadResultSuccess:
		return internal.UploadCompleted

	case internal.UploadResultPaused:
		return internal.UploadPaused

	case internal.UploadResultError:
		return internal.UploadFailed

	default:
		return internal.UploadFailed
	}
}
