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

	for chunk, err := range bfr.ReadChunk(1024) {
		if err != nil {
			return nil, err
		}

		hasher.Write(chunk.data)
	}

	hash := hasher.Sum(nil)

	return hash, nil
}
