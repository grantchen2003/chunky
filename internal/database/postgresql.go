package database

type Postgresql struct {
}

func NewPostgresql() (*Postgresql, error) {
	return &Postgresql{}, nil
}

func (p *Postgresql) CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error {
	return nil
}

func (p *Postgresql) Exists(sessionId string, fileHash []byte) (exists bool, err error) {
	return exists, err
}

func (p *Postgresql) AddFileChunk(sessionId string, fileHash []byte, chunkId string, startByte int, endByte int) error {
	return nil
}
