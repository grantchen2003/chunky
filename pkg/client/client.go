package client

type UploadStatus = string

const (
	UploadComplete UploadStatus = "statusComplete"
)

type Client struct {
	filePath          string
	url               string
	UploadedBytesChan chan (int)
	UploadErrorChan   chan (error)
	UploadStatusChan  chan (UploadStatus)
}

func NewClient(url string, filePath string) *Client {
	return &Client{filePath: filePath, url: url}
}

func (c *Client) Upload() error {
	return nil
}

func (c *Client) Pause() {
}

func (c *Client) Resume() {
}

// func (c *Client) UploadedBytesChan() <-chan (int) {
// 	return c.uploadedBytesChan
// }

// func (c *Client) UploadErrorChan() <-chan (error) {
// 	return c.uploadErrorChan
// }

// func (c *Client) UploadStatusChan() <-chan (struct{}) {
// 	return c.uploadStatusChan
// }

func (c *Client) UploadIsInProgress() bool {
	return false
}

func (c *Client) UploadIsCompleted() bool {
	return false
}

func (c *Client) UploadIsPaused() bool {
	return false
}
