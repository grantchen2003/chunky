package client

type UploadStatus int

type UploadProgress struct {
}

const (
	UploadCompleted UploadStatus = iota
	UploadStarted
	UploadIsPaused
	UploadResumed
)

type Client struct {
	filePath         string
	url              string
	ProgressChan     chan (UploadProgress)
	UploadErrorChan  chan (error)
	UploadStatusChan chan (UploadStatus)
}

func NewClient(url string, filePath string) *Client {
	return &Client{
		filePath:         filePath,
		url:              url,
		ProgressChan:     make(chan UploadProgress),
		UploadErrorChan:  make(chan error),
		UploadStatusChan: make(chan UploadStatus),
	}
}

func (c *Client) Upload() {
}

func (c *Client) Pause() {
}

func (c *Client) Resume() {
}
