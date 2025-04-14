package client

type Client struct {
	filePath          string
	url               string
	uploadedBytesChan chan (int)
	uploadErrorChan   chan (error)
	uploadStatusChan  chan (struct{})
}

func NewClient(url string, filePath string) *Client {
	return &Client{filePath: filePath, url: url}
}

func (c *Client) Upload() error {
	return nil
}

func (c *Client) Pause() error {
	return nil
}

func (c *Client) UploadedBytesChan() <-chan (int) {
	return c.uploadedBytesChan
}

func (c *Client) UploadErrorChan() <-chan (error) {
	return c.uploadErrorChan
}

func (c *Client) UploadStatusChan() <-chan (struct{}) {
	return c.uploadStatusChan
}

func (c *Client) UploadIsInProgress() bool {
	return false
}

func (c *Client) UploadIsCompleted() bool {
	return false
}

func (c *Client) UploadIsPaused() bool {
	return false
}
