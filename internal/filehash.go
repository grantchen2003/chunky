package internal

import "crypto/sha256"

// possibly make it faster with concurrency?
func HashFile(filePath string) ([]byte, error) {
	hasher := sha256.New()

	bfr, err := NewBufferedFileReader(filePath)
	if err != nil {
		return nil, err
	}
	defer bfr.Close()

	const bufferSizeBytes = 1 << 20 // 1 MiB
	for chunk, err := range bfr.ReadChunk(bufferSizeBytes) {
		if err != nil {
			return nil, err
		}

		hasher.Write(chunk.Data)
	}

	hash := hasher.Sum(nil)

	return hash, nil
}
