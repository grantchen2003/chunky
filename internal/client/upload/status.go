package upload

type Status = struct {
	Message string
}

var (
	UploadStarted Status = Status{Message: "upload started"}
	UploadResumed Status = Status{Message: "upload resumed"}
)
