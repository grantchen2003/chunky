package internal

type UploadNotifier struct {
	ProgressChan  chan UploadProgress
	ErrorChan     chan error
	StatusChan    chan UploadStatus
	UserErrorChan chan error
}

func NewUploadNotifier() *UploadNotifier {
	return &UploadNotifier{
		ProgressChan:  make(chan UploadProgress),
		ErrorChan:     make(chan error),
		StatusChan:    make(chan UploadStatus),
		UserErrorChan: make(chan error),
	}
}

func (un *UploadNotifier) Close() {
	close(un.ProgressChan)
	close(un.ErrorChan)
	close(un.StatusChan)
	close(un.UserErrorChan)
}
