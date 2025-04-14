package internal

type UploadStatus = struct {
	Message       string
	IsTerminating bool
}

var (
	UploadStarted UploadStatus = UploadStatus{Message: "upload started"}
	UploadResumed UploadStatus = UploadStatus{Message: "upload resumed"}
)
