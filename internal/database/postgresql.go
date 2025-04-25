package database

type Postgresql struct {
}

func NewPostgresql() *Postgresql {
	return &Postgresql{}
}

func (p *Postgresql) CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error {
	return nil
}
