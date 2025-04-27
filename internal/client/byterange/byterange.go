package byterange

import (
	"fmt"
)

type ByteRange struct {
	StartByte int
	EndByte   int
}

func NewByteRange(startByte int, endByte int) (ByteRange, error) {
	if startByte > endByte {
		return ByteRange{}, fmt.Errorf("endByte: %d is less than startByte: %d", endByte, startByte)
	}

	return ByteRange{
		StartByte: startByte,
		EndByte:   endByte,
	}, nil
}

func (br ByteRange) Size() int {
	return br.EndByte - br.StartByte + 1
}
