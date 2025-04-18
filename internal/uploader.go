package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// NEED TO REFACTOR
type Uploader struct {
	url                 string
	filePath            string
	uploadProgressChan  chan<- UploadProgress
	uploadSessionStorer UploadSessionStorer
}

func NewUploader(url string, filePath string, uploadProgressChan chan<- UploadProgress, uploadSessionStorer UploadSessionStorer) *Uploader {
	return &Uploader{
		url:                 url,
		filePath:            filePath,
		uploadProgressChan:  uploadProgressChan,
		uploadSessionStorer: uploadSessionStorer,
	}
}

// simulate file upload
func (u *Uploader) Upload() error {
	fileHash, err := hashFile(u.filePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(u.filePath)
	if err != nil {
		return err
	}

	totalFileSizeBytes := fileInfo.Size()

	sessionId, err := initiateUploadSession(fileHash, int(totalFileSizeBytes))
	if err != nil {
		return err
	}

	bfr, err := NewBufferedFileReader(u.filePath, 3)
	if err != nil {
		return err
	}
	defer bfr.Close()

	var startByte int
	for chunk := range bfr.ReadChunk() {
		chunkSize := len(chunk)
		endByte := startByte + chunkSize - 1
		err := u.uploadFileChunk(sessionId, fileHash, chunk, startByte, endByte)
		if err != nil {
			return err
		}
		startByte = endByte + 1

		u.uploadProgressChan <- UploadProgress{
			UploadedBytes: chunkSize,
		}

		u.uploadSessionStorer.Store(sessionId, u.filePath, fileHash)
	}

	return nil
}

func initiateUploadSession(fileHash []byte, totalFileSizeBytes int) (string, error) {
	fmt.Printf("Initiating upload session for totalFileSizeBytes: %d and fileHash: %v\n", totalFileSizeBytes, fileHash)
	sessionId := "t8y3euagvkqp8fuo"
	return sessionId, nil
}

// need to acutally implement
func (u *Uploader) uploadFileChunk(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
	fmt.Printf("Uploading to %s, sessionId: %s, fileHash: %v, chunk: %s, startByte: %d, endByte: %d\n", u.url, sessionId, fileHash, chunk, startByte, endByte)
	time.Sleep(1 * time.Second)

	errPercentage := 50
	randomNumber := rand.Intn(100)
	if randomNumber < errPercentage {
		return errors.New("Failed to upload")
	} else {
		return nil
	}
}
