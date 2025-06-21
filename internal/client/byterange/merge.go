package byterange

import (
	"fmt"
)

func mergePair(a ByteRange, b ByteRange) (ByteRange, error) {
	if !intersects(a, b) {
		return ByteRange{}, fmt.Errorf("byte ranges %v and %v don't intersect", a, b)
	}

	byteRange, err := NewByteRange(
		min(a.startByte, b.startByte),
		max(a.endByte, b.endByte),
	)

	return byteRange, err
}

func Merge(byteRanges []ByteRange) []ByteRange {
	sortByteRanges(byteRanges)

	var mergedByteRanges []ByteRange

	for _, br := range byteRanges {
		if len(mergedByteRanges) != 0 && intersects(mergedByteRanges[len(mergedByteRanges)-1], br) {
			lastMergedByteRange := mergedByteRanges[len(mergedByteRanges)-1]
			mergedByteRanges = mergedByteRanges[:len(mergedByteRanges)-1]
			var err error
			br, err = mergePair(lastMergedByteRange, br)
			if err != nil {
				panic(err)
			}
		}
		mergedByteRanges = append(mergedByteRanges, br)
	}

	return mergedByteRanges
}
