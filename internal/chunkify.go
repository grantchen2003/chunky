package internal

func Chunkify[T any](arr []T, chunkLen int) [][]T {
	var chunks [][]T
	for i := 0; i < len(arr); i += chunkLen {
		chunk := arr[i:min(i+chunkLen, len(arr))]
		chunks = append(chunks, chunk)
	}
	return chunks
}
