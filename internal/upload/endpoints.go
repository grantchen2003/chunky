package upload

type Endpoints struct {
	InitiateUploadSession string
	ByteRangesToUpload    string
	UploadFileChunk       string
}

func (e *Endpoints) PopulateEmptyFields() {
	if e.InitiateUploadSession == "" {
		e.InitiateUploadSession = "/initiateUploadSession"
	}

	if e.ByteRangesToUpload == "" {
		e.ByteRangesToUpload = "/byteRangesToUpload"
	}

	if e.UploadFileChunk == "" {
		e.UploadFileChunk = "/uploadFileChunk"
	}
}
