package internal

import (
	"fmt"

	"github.com/grantchen2003/chunky/internal/byterange"
)

// need to implement
type UploadRequester struct {
	url string
}

func newUploadRequester(url string) *UploadRequester {
	return &UploadRequester{
		url: url,
	}
}

func (ur UploadRequester) makeInitiateUploadSessionRequest(fileHash []byte, totalFileSizeBytes int) (string, error) {
	fmt.Printf("Initiating upload session for totalFileSizeBytes: %d and fileHash: %v\n", totalFileSizeBytes, fileHash)
	sessionId := "t8y3euagvkqp8fuo"
	return sessionId, nil
}

func (ur UploadRequester) makeByteRangesToUploadRequest(sessionId string, fileHash []byte) ([]byterange.ByteRange, error) {
	fmt.Printf("Getting byte ranges to upload for sessionId: %s and fileHash: %v\n", sessionId, fileHash)

	responseData := [][]int{{6, 100}, {102, 132}, {103, 104}, {133, 152}, {154, 154}, {155, 155}, {156, 159}, {161, 306}}

	var byteRangesToUpload []byterange.ByteRange

	for _, data := range responseData {
		byteRange, err := byterange.NewByteRange(data[0], data[1])
		if err != nil {
			return nil, err
		}
		byteRangesToUpload = append(byteRangesToUpload, byteRange)
	}

	return byteRangesToUpload, nil
}

func (ur UploadRequester) makeUploadFileChunkRequest(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
	fmt.Printf("Uploading to %s, sessionId: %s, fileHash: %v, startByte: %d, endByte: %d\n", ur.url, sessionId, fileHash, startByte, endByte)
	// time.Sleep(100 * time.Millisecond)
	return nil
}
