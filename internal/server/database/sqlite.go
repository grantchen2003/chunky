package database

// need to implement
type Sqlite struct {
}

func NewSqlite() (*Sqlite, error) {
	p := &Sqlite{}

	if err := p.createUploadSessionTableIfNotExists(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Sqlite) CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error {
	return nil
}

func (p *Sqlite) Exists(sessionId string, fileHash []byte) (exists bool, err error) {
	exists = true
	return exists, err
}

func (p *Sqlite) AddFileChunk(sessionId string, fileHash []byte, chunkId string, startByte int, endByte int) error {
	return nil
}

func (p *Sqlite) createUploadSessionTableIfNotExists() error {
	return nil
}

func (p *Sqlite) ByteRangesToUpload(sessionId string, fileHash []byte) ([][2]int, error) {
	return [][2]int{{6, 100}, {102, 132}, {103, 104}, {133, 152}, {154, 154}, {155, 155}, {156, 159}, {161, 306}}, nil
}
