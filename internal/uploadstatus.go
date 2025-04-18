package internal

type UploadStatus = struct {
	Message       string
	IsTerminating bool
}

var (
	UploadStarted UploadStatus = UploadStatus{Message: "upload started", IsTerminating: false}
	UploadResumed UploadStatus = UploadStatus{Message: "upload resumed", IsTerminating: false}
)
