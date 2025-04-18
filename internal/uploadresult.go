package internal

type UploadResult string

const (
	UploadResultSuccess UploadResult = "Success"
	UploadResultPaused  UploadResult = "Paused"
	UploadResultError   UploadResult = "Error"
)
