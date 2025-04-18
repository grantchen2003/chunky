package internal

import (
	"bufio"
	"io"
	"iter"
	"os"
)

type BufferedFileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func NewBufferedFileReader(filePath string) (*BufferedFileReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &BufferedFileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (bfr *BufferedFileReader) ReadChunk(bufferSizeBytes int) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for {
			buffer := make([]byte, bufferSizeBytes)

			bytesRead, err := bfr.reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}

				yield(nil, err)
				return
			}

			if !yield(buffer[:bytesRead], nil) {
				return
			}
		}
	}
}

func (bfr *BufferedFileReader) Close() {
	bfr.file.Close()
}
