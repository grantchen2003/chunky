package internal

func Chunkify[T any](arr []T, chunkLen int) [][]T {
	var chunks [][]T
	for i := 0; i < len(arr); i += chunkLen {
		endIndex := i + chunkLen
		if endIndex > len(arr) {
			endIndex = len(arr)
		}
		chunk := arr[i:endIndex]
		chunks = append(chunks, chunk)
	}
	return chunks
}
