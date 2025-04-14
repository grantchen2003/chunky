package internal

type UploadNotifier struct {
	ProgressChan chan UploadProgress
	ResultChan   chan UploadResult
	StatusChan   chan UploadStatus
}

func NewUploadNotifier() *UploadNotifier {
	return &UploadNotifier{
		ProgressChan: make(chan UploadProgress),
		ResultChan:   make(chan UploadResult),
		StatusChan:   make(chan UploadStatus),
	}
}

func (un *UploadNotifier) Close() {
	close(un.ProgressChan)
	close(un.ResultChan)
	close(un.StatusChan)
}
