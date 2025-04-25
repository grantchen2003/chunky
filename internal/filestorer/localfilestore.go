package filestorer

type LocalFileStore struct {
}

func NewLocalFileStore() *LocalFileStore {
	return &LocalFileStore{}
}
func (lfs *LocalFileStore) Store([]byte) (chunkId string, err error) {
	return chunkId, err
}
