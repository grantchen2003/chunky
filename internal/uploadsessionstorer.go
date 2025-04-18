package internal

type UploadSessionStorer interface {
	Store(sessionId string, filePath string, fileHash []byte) error
}
