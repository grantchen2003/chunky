package filestorer

type FileStore struct {
}

func NewFileStore() *FileStore {
	return &FileStore{}
}
func (fs *FileStore) Store([]byte) (chunkId string, err error) {
	return chunkId, err
}
