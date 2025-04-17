package internal

type UploadResult int

const (
	UploadResultUnknown UploadResult = iota
	UploadResultSuccess
	UploadResultPaused
	UploadResultError
)
