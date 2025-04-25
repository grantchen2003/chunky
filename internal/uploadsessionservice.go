package internal

import (
	"fmt"

	"github.com/grantchen2003/chunky/internal/database"
	"github.com/grantchen2003/chunky/internal/filestorer"
)

type UploadSessionService struct {
	db         database.Database
	fileStorer filestorer.FileStorer
}

func NewUploadSessionService(db database.Database, fileStorer filestorer.FileStorer) *UploadSessionService {
	return &UploadSessionService{
		db:         db,
		fileStorer: fileStorer,
	}
}

func (s *UploadSessionService) CreateUploadSession(fileHash []byte, totalFileSizeBytes int) (string, error) {
	sessionId, err := GenerateRandomHexString(16)
	if err != nil {
		return "", err
	}

	if err = s.db.CreateUploadSession(sessionId, fileHash, totalFileSizeBytes); err != nil {
		return "", err
	}

	return sessionId, err
}

// name this better
func (s *UploadSessionService) AddFileChunk(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
	uploadExists, err := s.db.Exists(sessionId, fileHash)
	if err != nil {
		return err
	}

	if !uploadExists {
		return fmt.Errorf("cannot add file chunk to non-existent upload with session id: %s and fileHash: %s", sessionId, fileHash)
	}

	chunkId, err := s.fileStorer.Store(chunk)
	if err != nil {
		return err
	}

	err = s.db.AddFileChunk(sessionId, fileHash, chunkId, startByte, endByte)
	if err != nil {
		return err
	}

	return nil
}
