package internal

type UploadManager interface {
	StartUpload() error
	PauseUpload() error
	ResumeUpload() error
}
