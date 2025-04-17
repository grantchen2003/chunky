package client

import (
	"context"
	"errors"

	"github.com/grantchen2003/chunky/internal"
)

var (
	ErrStartedOnOngoingUpload     = errors.New("upload start when an upload is already ongoing")
	ErrPausedOnNoOngoingUpload    = errors.New("paused when no upload is ongoing")
	ErrResumedOnNonExistingUpload = errors.New("resumed on non existing upload error")
	ErrResumedOnOngoingUpload     = errors.New("resumed on ongoing upload error")
)

var (
	UploadCompleted UploadStatus = UploadStatus{Message: "upload completed", IsTerminating: true}
	UploadFailed    UploadStatus = UploadStatus{Message: "upload failed", IsTerminating: true}
	UploadStarted   UploadStatus = UploadStatus{Message: "upload started", IsTerminating: false}
	UploadPaused    UploadStatus = UploadStatus{Message: "upload paused", IsTerminating: true}
	UploadResumed   UploadStatus = UploadStatus{Message: "upload resumed", IsTerminating: false}
)

type UploadProgress struct{}

type UploadStatus = struct {
	Message       string
	IsTerminating bool
}

type Client struct {
	uploadCtx       context.Context
	uploadCtxCancel context.CancelFunc

	UploadProgressChan chan (UploadProgress)
	UploadErrorChan    chan (error)
	UploadStatusChan   chan (UploadStatus)
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

		UploadProgressChan: make(chan UploadProgress),
		UploadErrorChan:    make(chan error),
		UploadStatusChan:   make(chan UploadStatus),
		UserErrorChan:      make(chan error),

		uploader: internal.NewUploader(),

		filePath: filePath,
		url:      url,
	}
}

func (c *Client) Upload() {
	if c.uploader.IsUploading() {
		c.UserErrorChan <- ErrStartedOnOngoingUpload
		return
	}

	c.handleUpload(UploadStarted)
}

func (c *Client) Pause() {
	if !c.uploader.IsUploading() {
		c.UserErrorChan <- ErrPausedOnNoOngoingUpload
		return
	}

	c.uploadCtxCancel()
}

func (c *Client) Resume() {
	if c.uploader.IsUploading() {
		c.UserErrorChan <- ErrResumedOnOngoingUpload
		return
	}

	if c.uploader.HasNoExistingupload() {
		c.UserErrorChan <- ErrResumedOnNonExistingUpload
		return
	}

	c.handleUpload(UploadResumed)
}

func (c *Client) handleUpload(uploadStatus UploadStatus) {
	c.uploadCtx, c.uploadCtxCancel = context.WithCancel(context.Background())

	c.UploadStatusChan <- uploadStatus

	uploadResult := c.uploader.Upload(c.uploadCtx, c.url, c.filePath)

	switch uploadResult {
	case internal.UploadResultSuccess:
		c.UploadStatusChan <- UploadCompleted
		close(c.UploadProgressChan)
		close(c.UploadErrorChan)
		close(c.UploadStatusChan)

	case internal.UploadResultPaused:
		c.UploadStatusChan <- UploadPaused

	case internal.UploadResultError:
		c.UploadStatusChan <- UploadFailed
	}
}
