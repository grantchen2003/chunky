package byterange

import "sort"

func sortByteRanges(byteRanges []ByteRange) {
	sort.Slice(byteRanges, func(i int, j int) bool {
		if byteRanges[i].startByte == byteRanges[j].startByte {
			return byteRanges[i].endByte <= byteRanges[j].endByte
		}

		return byteRanges[i].startByte < byteRanges[j].startByte
	})
}
