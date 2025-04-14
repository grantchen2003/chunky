package internal

import "github.com/grantchen2003/chunky/internal/byterange"

type FileChunk struct {
	byteRange byterange.ByteRange
	data      []byte
}
