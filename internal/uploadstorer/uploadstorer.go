package uploadstorer

import "errors"

var ErrNotFound = errors.New("not found")

type UploadStorer interface {
	Close() error
	Store(sessionId string, url string, filePath string, fileHash []byte) error
	GetSessionIdAndFileHash(url string, filePath string) (string, []byte, error)
}
