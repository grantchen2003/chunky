package client

type UploadStatus string

const (
	UploadComplete UploadStatus = "upload complete"
)

type Client struct {
	file string
	url                   string
	uploadedBytesChan chan(int)
	uploadErrorChan chan (error)
	uploadStatusChan chan (UploadStatus)
}

func NewClient(url string, file []byte) *Client {
	return &Client{file: file, url: url}
}

func (c *Client) Upload(file []byte) error {
	return nil
}

func (c *Client) Pause() error {
	return nil
}

func (c *Client) UploadedBytesChan() <-chan(int){
	return c.uploadedBytesChan
}

func (c *Client) UploadErrorChan() <-chan(error){
	return c.uploadErrorChan
}

func (c *Client) UploadStatusChan() <- chan (UploadStatus) {
	return c
}