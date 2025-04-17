package client

import (
	"context"
	"fmt"

	"github.com/grantchen2003/chunky/internal"
)

type UploadProgress struct{}

type UploadStatus = string

const (
	UploadCompleted UploadStatus = "upload completed"
	UploadStarted   UploadStatus = "upload started"
	UploadPaused    UploadStatus = "upload paused"
	UploadRestarted UploadStatus = "upload restarted"
	UploadResumed   UploadStatus = "upload resumed"
)

type Client struct {
	uploadCtx       context.Context
	uploadCtxCancel context.CancelFunc

	ProgressChan     chan (UploadProgress)
	UploadErrorChan  chan (error)
	UploadStatusChan chan (UploadStatus)

	filePath    internal.FilePath
	fileTracker *internal.FileTracker
	isUploading bool
	url         string
}

func NewClient(url string, filePathStr string) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	filePath := internal.FilePath(filePathStr)

	return &Client{
		uploadCtx:        ctx,
		uploadCtxCancel:  cancel,
		ProgressChan:     make(chan UploadProgress),
		UploadErrorChan:  make(chan error),
		UploadStatusChan: make(chan UploadStatus),
		filePath:         filePath,
		fileTracker:      internal.NewFileTracker(filePath),
		url:              url,
		isUploading:      false,
	}
}

func (c *Client) Upload() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot upload when an upload is already ongoing")
		return
	}

	c.handleUpload(UploadStarted)
}

func (c *Client) Pause() {
	if !c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot pause when no upload is ongoing")
		return
	}

	c.uploadCtxCancel()
}

func (c *Client) Resume() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot resume when an upload is already ongoing")
		return
	}

	if c.fileTracker.FileHasChangedSincePause() {
		c.handleUpload(UploadRestarted)
	} else {
		c.handleUpload(UploadResumed)
	}
}

func (c *Client) handleUpload(uploadStatus UploadStatus) {
	c.isUploading = true
	defer func() { c.isUploading = false }()

	c.uploadCtx, c.uploadCtxCancel = context.WithCancel(context.Background())

	c.UploadStatusChan <- uploadStatus

	uploadResult := internal.Upload(c.uploadCtx, c.url, c.filePath)

	switch uploadResult {
	case internal.UploadResultSuccess:
		close(c.ProgressChan)
		close(c.UploadErrorChan)
		close(c.UploadStatusChan)
		c.UploadStatusChan <- UploadCompleted

	case internal.UploadResultPaused:
		c.UploadStatusChan <- UploadPaused

	case internal.UploadResultError:
		c.UploadErrorChan <- fmt.Errorf("%s", string(fmt.Sprint(uploadResult)))
	}
}
