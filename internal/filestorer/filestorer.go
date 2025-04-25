package filestorer

type FileStorer interface {
	Store([]byte) (chunkId string, err error)
}
