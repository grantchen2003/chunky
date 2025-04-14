package byterange

func intersects(a ByteRange, b ByteRange) bool {
	return !(a.EndByte < b.StartByte || b.EndByte < a.StartByte)
}
