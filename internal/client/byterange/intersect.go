package byterange

func intersects(a ByteRange, b ByteRange) bool {
	return !(a.endByte < b.startByte || b.endByte < a.startByte)
}
