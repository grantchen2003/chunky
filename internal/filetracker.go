package internal

type FilePath string

type FileTracker struct {
	filePath FilePath
}

// create singleton instance?
func NewFileTracker(filePath FilePath) *FileTracker {
	return &FileTracker{
		filePath: filePath,
	}
}

func (ft *FileTracker) FileHasChangedSincePause() bool {
	return false
}

func (ft *FileTracker) ByteRangesToUpload() []Range {
	var x []Range
	return x
}
