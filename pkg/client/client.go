package client

import (
	"context"

	"github.com/grantchen2003/chunky/internal"
)

type Client struct {
	uploadCtx       context.Context
	uploadCtxCancel context.CancelFunc

	UploadProgressChan chan (internal.UploadProgress)
	UploadErrorChan    chan (error)
	UploadStatusChan   chan (internal.UploadStatus)
	UserErrorChan      chan (error)

	uploader *internal.Uploader

	filePath internal.FilePath
	url      string
}

func NewClient(url string, filePathStr string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	filePath := internal.FilePath(filePathStr)

	return &Client{
		uploadCtx:       ctx,
		uploadCtxCancel: cancel,

		UploadProgressChan: make(chan internal.UploadProgress),
		UploadErrorChan:    make(chan error),
		UploadStatusChan:   make(chan internal.UploadStatus),
		UserErrorChan:      make(chan error),

		uploader: internal.NewUploader(),

		filePath: filePath,
		url:      url,
	}
}

func (c *Client) Upload() {
	if c.uploader.IsUploading() {
		c.UserErrorChan <- internal.ErrStartedOnOngoingUpload
		return
	}

	c.handleUpload(internal.UploadStarted)
}

func (c *Client) Pause() {
	if !c.uploader.IsUploading() {
		c.UserErrorChan <- internal.ErrPausedOnNoOngoingUpload
		return
	}

	c.uploadCtxCancel()
}

func (c *Client) Resume() {
	if c.uploader.IsUploading() {
		c.UserErrorChan <- internal.ErrResumedOnOngoingUpload
		return
	}

	if c.uploader.HasNoExistingupload() {
		c.UserErrorChan <- internal.ErrResumedOnNonExistingUpload
		return
	}

	if c.uploader.FileHasChangedSinceLastUpload() {
		c.UserErrorChan <- internal.ErrResumedOnChangedFile
		return
	}

	c.handleUpload(internal.UploadResumed)
}

func (c *Client) handleUpload(uploadStatus internal.UploadStatus) {
	c.uploadCtx, c.uploadCtxCancel = context.WithCancel(context.Background())

	c.UploadStatusChan <- uploadStatus

	uploadResult := c.uploader.Upload(c.uploadCtx, c.url, c.filePath, c.UploadProgressChan)

	switch uploadResult {
	case internal.UploadResultSuccess:
		c.UploadStatusChan <- internal.UploadCompleted
		close(c.UploadProgressChan)
		close(c.UploadErrorChan)
		close(c.UploadStatusChan)

	case internal.UploadResultPaused:
		c.UploadStatusChan <- internal.UploadPaused

	case internal.UploadResultError:
		c.UploadStatusChan <- internal.UploadFailed
	}
}
