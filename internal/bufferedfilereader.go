package internal

import (
	"bufio"
	"io"
	"iter"
	"os"

	"github.com/grantchen2003/chunky/internal/byterange"
)

// implement ReadChunkWithRange
// maybe remove reader field, and just initialize it inside the ReadChunk method and close in this method too
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

func (bfr *BufferedFileReader) ReadChunk(bufferSizeBytes int) iter.Seq2[FileChunk, error] {
	return func(yield func(FileChunk, error) bool) {
		var startByte int
		for {
			buffer := make([]byte, bufferSizeBytes)

			bytesRead, err := bfr.reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}

				yield(FileChunk{}, err)
				return
			}

			byteRange, err := byterange.NewByteRange(startByte, startByte+bytesRead-1)
			if err != nil {
				yield(FileChunk{}, err)
				return
			}

			if !yield(FileChunk{ByteRange: byteRange, Data: buffer[:bytesRead]}, nil) {
				return
			}

			startByte += bufferSizeBytes
		}
	}
}

func (bfr *BufferedFileReader) ReadChunkWithRange(bufferSizeBytes int, byteRanges []byterange.ByteRange) iter.Seq2[FileChunk, error] {
	groupedByteRanges := byterange.GroupByteRanges(byteRanges, bufferSizeBytes)

	return func(yield func(FileChunk, error) bool) {
		for _, byteRangeGroup := range groupedByteRanges {
			offset := byteRangeGroup[0].StartByte

			_, err := bfr.file.Seek(int64(offset), io.SeekStart)
			if err != nil {
				yield(FileChunk{}, err)
				return
			}

			bufferSizeBytes = byteRangeGroup[len(byteRangeGroup)-1].EndByte - byteRangeGroup[0].StartByte + 1
			buffer := make([]byte, bufferSizeBytes)
			_, err = bfr.reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}

				yield(FileChunk{}, err)
				return
			}

			for _, byteRange := range byteRangeGroup {
				fc := FileChunk{
					ByteRange: byteRange,
					Data:      buffer[byteRange.StartByte-offset : byteRange.EndByte-offset],
				}

				if !yield(fc, nil) {
					return
				}
			}
		}
	}
}

func (bfr *BufferedFileReader) Close() {
	bfr.file.Close()
}
