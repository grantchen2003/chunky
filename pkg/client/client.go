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

	filePath    string
	isUploading bool
	url         string
}

func NewClient(url string, filePath string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		uploadCtx:        ctx,
		uploadCtxCancel:  cancel,
		ProgressChan:     make(chan UploadProgress),
		UploadErrorChan:  make(chan error),
		UploadStatusChan: make(chan UploadStatus),
		filePath:         filePath,
		url:              url,
		isUploading:      false,
	}
}

func (c *Client) Upload() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot upload when upload is ongoing")
		return
	}

	c.handleUpload(UploadStarted)
}

func (c *Client) Pause() {
	if !c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot pause when no uploads are ongoing")
		return
	}

	c.uploadCtxCancel()

	c.UploadStatusChan <- UploadPaused
}

func (c *Client) Resume() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot resume when upload is already ongoing")
		return
	}

	if c.isSameFile() {
		c.handleUpload(UploadResumed)
	} else {
		c.handleUpload(UploadRestarted)
	}
}

func (c *Client) handleUpload(uploadStatus UploadStatus) {
	c.resetUploadContext()
	c.isUploading = true
	c.UploadStatusChan <- uploadStatus

	err := Upload(c.url, c.filePath, c.byteRangesToUpload(), c.uploadCtx)

	if err != nil {
		if _, ok := err.(*UploadCancelledByPauseError); !ok {
			c.UploadErrorChan <- err
		}

		c.isUploading = false
		return
	}

	c.isUploading = false
	c.UploadStatusChan <- UploadCompleted

	close(c.ProgressChan)
	close(c.UploadErrorChan)
	close(c.UploadStatusChan)
}

func (c *Client) byteRangesToUpload() []internal.Range {
	var x []internal.Range
	return x
}

func (c *Client) isSameFile() bool {
	return true
}

func (c *Client) resetUploadContext() {
	ctx, cancel := context.WithCancel(context.Background())
	c.uploadCtx = ctx
	c.uploadCtxCancel = cancel
}
