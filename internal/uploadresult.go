package internal

type UploadResult int

const (
	UploadResultSuccess = iota
	UploadResultPaused
	UploadResultError
)
