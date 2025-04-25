package internal

import "github.com/grantchen2003/chunky/internal/byterange"

type FileChunk struct {
	ByteRange byterange.ByteRange
	Data      []byte
}
