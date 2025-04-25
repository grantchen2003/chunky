package internal

import "github.com/grantchen2003/chunky/internal/database"

type UploadSessionService struct {
	db database.Database
}

func NewUploadSessionService(db database.Database) *UploadSessionService {
	return &UploadSessionService{
		db: db,
	}
}

func (uss *UploadSessionService) CreateUploadSession(fileHash []byte, totalFileSizeBytes int) (string, error) {
	sessionId, err := GenerateSessionId(16)
	if err != nil {
		return "", err
	}

	if err = uss.db.CreateUploadSession(sessionId, fileHash, totalFileSizeBytes); err != nil {
		return "", err
	}

	return sessionId, err
}
