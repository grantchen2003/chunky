package internal

type UploadEndpoints struct {
	InitiateUploadSession string
	ByteRangesToUpload    string
	UploadFileChunk       string
}

func (se *UploadEndpoints) PopulateEmptyFields() {
	if se.InitiateUploadSession == "" {
		se.InitiateUploadSession = "/initiateUploadSession"
	}
	if se.ByteRangesToUpload == "" {
		se.ByteRangesToUpload = "/byteRangesToUpload"
	}
	if se.UploadFileChunk == "" {
		se.UploadFileChunk = "/uploadFileChunk"
	}
}
