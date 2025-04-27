package file

import "github.com/grantchen2003/chunky/internal/client/byterange"

type FileChunk struct {
	ByteRange byterange.ByteRange
	Data      []byte
}
