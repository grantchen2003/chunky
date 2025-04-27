package byterange

import "sort"

func sortByteRanges(byteRanges []ByteRange) {
	sort.Slice(byteRanges, func(i int, j int) bool {
		if byteRanges[i].StartByte == byteRanges[j].StartByte {
			return byteRanges[i].EndByte <= byteRanges[j].EndByte
		}

		return byteRanges[i].StartByte < byteRanges[j].StartByte
	})
}
