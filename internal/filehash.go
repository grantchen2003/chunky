package internal

import "crypto/sha256"

func hashFile(filePath string) ([]byte, error) {
	hasher := sha256.New()

	bfr, err := NewBufferedFileReader(filePath)
	if err != nil {
		return nil, err
	}
	defer bfr.Close()

	for chunk := range bfr.ReadChunk(1024) {
		hasher.Write(chunk)
	}

	hash := hasher.Sum(nil)

	return hash, nil
}
