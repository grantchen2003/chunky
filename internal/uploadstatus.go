package internal

type UploadStatus = struct {
	Message       string
	IsTerminating bool
}

var (
	// UploadCompleted UploadStatus = UploadStatus{Message: "upload completed", IsTerminating: true}
	// UploadFailed    UploadStatus = UploadStatus{Message: "upload failed", IsTerminating: true}
	UploadStarted UploadStatus = UploadStatus{Message: "upload started", IsTerminating: false}
	// UploadPaused    UploadStatus = UploadStatus{Message: "upload paused", IsTerminating: true}
	UploadResumed UploadStatus = UploadStatus{Message: "upload resumed", IsTerminating: false}
)
