package byterange

func TotalByteCount(byteRanges []ByteRange) int {
	var count int

	for _, br := range Merge(byteRanges) {
		count += br.Size()
	}

	return count
}
