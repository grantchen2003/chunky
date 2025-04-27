package byterange

import "fmt"

func GroupByteRanges(byteRanges []ByteRange, maxGroupSizeBytes int) [][]ByteRange {
	if maxGroupSizeBytes <= 0 {
		panic(fmt.Sprintf("maxGroupSizeBytes has to be greater than 0, received %d", maxGroupSizeBytes))
	}

	sortByteRanges(byteRanges)
	mergedAndSortedByteRanges := Merge(byteRanges)

	var groups [][]ByteRange
	var group []ByteRange

	for _, br := range mergedAndSortedByteRanges {
		for {
			groupCanAddEntireByteRange :=
				len(group) == 0 && br.Size() <= maxGroupSizeBytes ||
					len(group) != 0 && br.EndByte-group[0].StartByte+1 <= maxGroupSizeBytes

			if groupCanAddEntireByteRange {
				group = append(group, br)
				break
			}

			groupCanAddPartialByteRange :=
				len(group) == 0 ||
					br.StartByte-group[0].StartByte+1 <= maxGroupSizeBytes

			if groupCanAddPartialByteRange {
				groupStartByte := br.StartByte
				if len(group) != 0 {
					groupStartByte = group[0].StartByte
				}

				brToAddEndByte := groupStartByte + maxGroupSizeBytes - 1
				brToAdd, err := NewByteRange(br.StartByte, brToAddEndByte)
				if err != nil {
					panic(err)
				}

				group = append(group, brToAdd)
				br.StartByte = brToAddEndByte + 1
				continue
			}

			groups = append(groups, group)
			group = []ByteRange{}
		}
	}

	if len(group) > 0 {
		groups = append(groups, group)
	}

	return groups
}
