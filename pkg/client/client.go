package client

import (
	"context"
	"fmt"
)

type UploadStatus = string

type UploadProgress struct {
}

const (
	UploadCompleted UploadStatus = "upload completed"
	UploadStarted   UploadStatus = "upload started"
	UploadPaused    UploadStatus = "upload paused"
	UploadRestarted UploadStatus = "upload restarted"
	UploadResumed   UploadStatus = "upload resumed"
)

type Client struct {
	filePath         string
	url              string
	ctx              context.Context
	ctxCancel        context.CancelFunc
	isUploading      bool
	ProgressChan     chan (UploadProgress)
	UploadErrorChan  chan (error)
	UploadStatusChan chan (UploadStatus)
}

func NewClient(url string, filePath string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		filePath:         filePath,
		url:              url,
		isUploading:      false,
		ctx:              ctx,
		ctxCancel:        cancel,
		ProgressChan:     make(chan UploadProgress),
		UploadErrorChan:  make(chan error),
		UploadStatusChan: make(chan UploadStatus),
	}
}

func (c *Client) Upload() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot resume upload when upload is ongoing")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.ctxCancel = cancel
	c.isUploading = true

	c.UploadStatusChan <- UploadStarted
	err := Upload(c.url, c.filePath, c.ctx)
	if err != nil {
		c.isUploading = false
		return
	}
	c.isUploading = false
	c.UploadStatusChan <- UploadCompleted
	close(c.ProgressChan)
	close(c.UploadErrorChan)
	close(c.UploadStatusChan)

}

func (c *Client) Pause() {
	if !c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot pause upload when no uploads are ongoing")
		return
	}
	c.ctxCancel()
	c.UploadStatusChan <- UploadPaused
}

func (c *Client) Resume() {
	if c.isUploading {
		c.UploadErrorChan <- fmt.Errorf("cannot resume upload when upload is ongoing")
		return
	}

	c.UploadStatusChan <- UploadResumed
	changesSincePause := true
	if changesSincePause {
		c.UploadStatusChan <- UploadRestarted
		c.Upload()
		return
	} else {

	}
}
