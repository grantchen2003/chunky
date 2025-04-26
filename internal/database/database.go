package database

type Database interface {
	CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error
	Exists(sessionId string, fileHash []byte) (exists bool, err error)
	AddFileChunk(sessionId string, fileHash []byte, chunkId string, startByte int, endByte int) error
	ByteRangesToUpload(sessionId string, fileHash []byte) ([][2]int, error)
}
