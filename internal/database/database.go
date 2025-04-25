package database

type Database interface {
	CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error
}
