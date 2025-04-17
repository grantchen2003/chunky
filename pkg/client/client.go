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

func (c *Client) Upload() error {
	if err := c.uploadManager.ValidateUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadStarted

	uploadResult := c.uploadManager.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

	c.uploadNotifier.StatusChan <- determineUploadStatus(uploadResult)

	c.uploadNotifier.Close()

	return nil
}

func (c *Client) Pause() error {
	if err := c.uploadManager.ValidatePauseUpload(); err != nil {
		return err
	}

	c.uploadManager.PauseUpload()

	return nil
}

func (c *Client) Resume() error {
	if err := c.uploadManager.ValidateResumeUpload(); err != nil {
		return err
	}

	c.uploadNotifier.StatusChan <- internal.UploadResumed

	uploadResult := c.uploadManager.Upload(c.url, c.filePath, c.uploadNotifier.ProgressChan)

	c.uploadNotifier.StatusChan <- determineUploadStatus(uploadResult)

	c.uploadNotifier.Close()

	return nil
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
