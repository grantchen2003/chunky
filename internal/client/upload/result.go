package upload

type Result string

const (
	UploadResultSuccess Result = "Success"
	UploadResultPaused  Result = "Paused"
	UploadResultError   Result = "Error"
)
