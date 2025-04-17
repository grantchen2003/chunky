package client

import (
	"github.com/grantchen2003/chunky/internal"
)

type Client struct {
	filePath internal.FilePath
	url      string

	uploadNotifier *internal.UploadNotifier
	uploader       *internal.Uploader
}

func NewClient(url string, filePathStr string) *Client {
	return &Client{
		filePath: internal.FilePath(filePathStr),
		url:      url,

		uploadNotifier: internal.NewUploadNotifier(),
		uploader:       internal.NewUploader(),
	}
}

func (c *Client) Upload() {
	if err := c.uploader.ValidateUpload(); err != nil {
		c.uploadNotifier.UserErrorChan <- err
	}

	c.upload(internal.UploadStarted)
}

func (c *Client) Pause() {
	if err := c.uploader.ValidatePause(); err != nil {
		c.uploadNotifier.UserErrorChan <- err
	}

	c.uploader.PauseUpload()
}

func (c *Client) Resume() {
	if err := c.uploader.ValidateResume(); err != nil {
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

	uploadResult := c.uploader.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

	c.uploadNotifier.StatusChan <- uploadResultToUploadStatus(uploadResult)

	c.uploadNotifier.Close()
}

func uploadResultToUploadStatus(uploadResult internal.UploadResult) internal.UploadStatus {
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
