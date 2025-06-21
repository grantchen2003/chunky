package byterange

import (
	"fmt"
)

type ByteRange struct {
	startByte int
	endByte   int
}

func NewByteRange(startByte int, endByte int) (ByteRange, error) {
	if startByte > endByte {
		return ByteRange{}, fmt.Errorf("endByte: %d is less than startByte: %d", endByte, startByte)
	}

	return ByteRange{
		startByte: startByte,
		endByte:   endByte,
	}, nil
}

func (br *ByteRange) StartByte() int {
	return br.startByte
}

func (br *ByteRange) EndByte() int {
	return br.endByte
}

func (br *ByteRange) Size() int {
	return br.endByte - br.startByte + 1
}
