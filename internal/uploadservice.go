package internal

import (
	"fmt"

	"github.com/grantchen2003/chunky/internal/database"
	"github.com/grantchen2003/chunky/internal/filestorer"
	"github.com/grantchen2003/chunky/internal/util"
)

type UploadService struct {
	db         database.Database
	fileStorer filestorer.FileStorer
}

func NewUploadService(db database.Database, fileStorer filestorer.FileStorer) *UploadService {
	return &UploadService{
		db:         db,
		fileStorer: fileStorer,
	}
}

func (s *UploadService) CreateUploadSession(fileHash []byte, totalFileSizeBytes int) (string, error) {
	sessionId, err := util.GenerateRandomHexString(16)
	if err != nil {
		return "", err
	}

	if err = s.db.CreateUploadSession(sessionId, fileHash, totalFileSizeBytes); err != nil {
		return "", err
	}

	return sessionId, err
}

// name this better
func (s *UploadService) AddFileChunk(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
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

func (s *UploadService) ByteRangesToUpload(sessionId string, fileHash []byte) ([][2]int, error) {
	uploadExists, err := s.db.Exists(sessionId, fileHash)
	if err != nil {
		return nil, err
	}

	if !uploadExists {
		return nil, fmt.Errorf("cannot add file chunk to non-existent upload with session id: %s and fileHash: %s", sessionId, fileHash)
	}

	byteRangesToUpload, err := s.db.ByteRangesToUpload(sessionId, fileHash)
	if err != nil {
		return nil, err
	}

	return byteRangesToUpload, nil
}
