package internal

import (
	"os"

	"github.com/grantchen2003/chunky/internal/byterange"
	us "github.com/grantchen2003/chunky/internal/uploadstorer"
)

// NEED TO REFACTOR
type Uploader struct {
	url                string
	filePath           string
	uploadProgressChan chan<- UploadProgress
	uploadStorer       us.UploadStorer
	uploadRequester    *UploadRequester
}

func NewUploader(url string, filePath string, uploadProgressChan chan<- UploadProgress, uploadStorer us.UploadStorer, uploadRequester *UploadRequester) *Uploader {
	return &Uploader{
		url:                url,
		filePath:           filePath,
		uploadProgressChan: uploadProgressChan,
		uploadStorer:       uploadStorer,
		uploadRequester:    uploadRequester,
	}
}

func (u *Uploader) Upload() error {
	fileHash, err := HashFile(u.filePath)
	if err != nil {
		return err
	}

	sessionId, err := u.initiateUploadSession(fileHash)
	if err != nil {
		return err
	}

	err = u.streamFileUpload(sessionId, fileHash)
	return err
}

func (u *Uploader) ResumeUpload() error {
	sessionId, fileHash, err := u.uploadStorer.GetSessionIdAndFileHash(u.url, u.filePath)
	if err != nil {
		return err
	}

	byteRangesToUpload, err := u.byteRangesToUpload(sessionId, fileHash)
	if err != nil {
		return err
	}

	err = u.streamFileResumeUpload(sessionId, fileHash, byteRangesToUpload)
	return err
}

func (u *Uploader) initiateUploadSession(fileHash []byte) (string, error) {
	totalFileSizeBytes, err := u.getFileSizeBytes()
	if err != nil {
		return "", err
	}

	sessionId, err := u.uploadRequester.makeInitiateUploadSessionRequest(fileHash, totalFileSizeBytes)
	return sessionId, err
}

func (u *Uploader) getFileSizeBytes() (int, error) {
	fileInfo, err := os.Stat(u.filePath)
	if err != nil {
		return 0, err
	}

	totalFileSizeBytes := fileInfo.Size()
	return int(totalFileSizeBytes), nil
}

func (u *Uploader) byteRangesToUpload(sessionId string, fileHash []byte) ([]byterange.ByteRange, error) {
	byteRanges, err := u.uploadRequester.makeByteRangesToUploadRequest(sessionId, fileHash)
	return byteRanges, err
}

func (u *Uploader) streamFileUpload(sessionId string, fileHash []byte) error {
	bfr, err := NewBufferedFileReader(u.filePath)
	if err != nil {
		return err
	}
	defer bfr.Close()

	fileSizeBytes, err := u.getFileSizeBytes()
	if err != nil {
		return err
	}

	const bufferSizeBytes = 1 << 20 // 1 MiB
	for fileChunk, err := range bfr.ReadChunk(bufferSizeBytes) {
		if err != nil {
			return err
		}

		// Do not make this concurrent: uploading chunks in parallel would bypass
		// the buffered reader's memory management, potentially loading the entire
		// file into memory. Sequential uploads preserve the intended low memory footprint.
		u.uploadFileChunkWithProgress(sessionId, fileHash, fileChunk, fileSizeBytes)
	}

	return nil
}

func (u *Uploader) streamFileResumeUpload(sessionId string, fileHash []byte, byteRanges []byterange.ByteRange) error {
	bfr, err := NewBufferedFileReader(u.filePath)
	if err != nil {
		return err
	}
	defer bfr.Close()

	var totalBytesToUpload int
	for _, br := range byterange.Merge(byteRanges) {
		totalBytesToUpload += br.Size()
	}

	const bufferSizeBytes = 1 << 20 // 1 MiB
	for fileChunk, err := range bfr.ReadChunkWithRange(bufferSizeBytes, byteRanges) {
		if err != nil {
			return err
		}

		// Do not make this concurrent: uploading chunks in parallel would bypass
		// the buffered reader's memory management, potentially loading the entire
		// file into memory. Sequential uploads preserve the intended low memory footprint.
		u.uploadFileChunkWithProgress(sessionId, fileHash, fileChunk, totalBytesToUpload)
	}

	return nil
}

func (u *Uploader) uploadFileChunkWithProgress(sessionId string, fileHash []byte, fileChunk FileChunk, totalBytesToUpload int) error {
	err := u.uploadRequester.makeUploadFileChunkRequest(sessionId, fileHash, fileChunk.data, fileChunk.byteRange.StartByte, fileChunk.byteRange.EndByte)
	if err != nil {
		return err
	}

	err = u.uploadStorer.Store(sessionId, u.url, u.filePath, fileHash)
	if err != nil {
		return err
	}

	u.uploadProgressChan <- UploadProgress{UploadedBytes: fileChunk.byteRange.Size(), TotalBytesToUpload: totalBytesToUpload}

	return nil
}
