package internal

type UploadSessionStorer interface {
	Close() error
	Store(sessionId string, url string, filePath string, fileHash []byte) error
	GetSessionIdAndHash(url string, filePath string) (string, []byte, error)
}
