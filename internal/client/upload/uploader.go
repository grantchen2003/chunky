package upload

import (
	"context"
	"os"

	"github.com/grantchen2003/chunky/internal/client/byterange"
	"github.com/grantchen2003/chunky/internal/client/file"
	us "github.com/grantchen2003/chunky/internal/client/upload/uploadstorer"
)

type Uploader struct {
	url             string
	filePath        string
	progressChan    chan<- Progress
	uploadStorer    us.UploadStorer
	uploadRequester *Requester
}

func NewUploader(
	url string,
	filePath string,
	progressChan chan<- Progress,
	uploadStorer us.UploadStorer,
	uploadRequester *Requester,
) *Uploader {
	return &Uploader{
		url:             url,
		filePath:        filePath,
		progressChan:    progressChan,
		uploadStorer:    uploadStorer,
		uploadRequester: uploadRequester,
	}
}

func (u *Uploader) Upload(ctx context.Context) error {
	fileHash, err := file.HashFile(u.filePath)
	if err != nil {
		return err
	}

	sessionId, err := u.initiateUploadSession(fileHash)
	if err != nil {
		return err
	}

	err = u.streamFileUpload(ctx, sessionId, fileHash)
	return err
}

func (u *Uploader) ResumeUpload(ctx context.Context) error {
	sessionId, fileHash, err := u.uploadStorer.GetSessionIdAndFileHash(
		u.url, u.filePath,
	)
	if err != nil {
		return err
	}

	byteRangesToUpload, err := u.byteRangesToUpload(sessionId, fileHash)
	if err != nil {
		return err
	}

	err = u.streamFileResumeUpload(ctx, sessionId, fileHash, byteRangesToUpload)
	return err
}

func (u *Uploader) initiateUploadSession(fileHash []byte) (string, error) {
	totalFileSizeBytes, err := u.getFileSizeBytes()
	if err != nil {
		return "", err
	}

	sessionId, err := u.uploadRequester.makeInitiateUploadSessionRequest(
		fileHash, totalFileSizeBytes,
	)
	if err != nil {
		return "", err
	}

	err = u.uploadStorer.Store(sessionId, u.url, u.filePath, fileHash)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (u *Uploader) getFileSizeBytes() (int, error) {
	fileInfo, err := os.Stat(u.filePath)
	if err != nil {
		return 0, err
	}

	totalFileSizeBytes := fileInfo.Size()
	return int(totalFileSizeBytes), nil
}

func (u *Uploader) byteRangesToUpload(
	sessionId string,
	fileHash []byte,
) ([]byterange.ByteRange, error) {
	byteRanges, err := u.uploadRequester.makeByteRangesToUploadRequest(
		sessionId, fileHash,
	)

	return byteRanges, err
}

func (u *Uploader) streamFileUpload(
	ctx context.Context,
	sessionId string,
	fileHash []byte,
) error {
	bfr, err := file.NewBufferedFileReader(u.filePath)
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

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Do not make this concurrent: uploading chunks in parallel would
		// bypass the buffered reader's memory management, potentially loading
		// the entire file into memory. Sequential uploads preserve the
		// intended low memory footprint.
		if err := u.uploadFileChunkWithProgress(
			sessionId,
			fileHash,
			fileChunk,
			fileSizeBytes,
		); err != nil {
			return err
		}
	}

	return nil
}

func (u *Uploader) streamFileResumeUpload(
	ctx context.Context,
	sessionId string,
	fileHash []byte,
	byteRanges []byterange.ByteRange,
) error {
	bfr, err := file.NewBufferedFileReader(u.filePath)
	if err != nil {
		return err
	}
	defer bfr.Close()

	totalBytesToUpload := byterange.TotalByteCount(byteRanges)

	const bufferSizeBytes = 1 << 20 // 1 MiB
	for fileChunk, err := range bfr.ReadChunkWithRange(bufferSizeBytes, byteRanges) {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Do not make this concurrent: uploading chunks in parallel would
		// bypass the buffered reader's memory management, potentially loading
		// the entire file into memory. Sequential uploads preserve the
		// intended low memory footprint.
		if err := u.uploadFileChunkWithProgress(
			sessionId,
			fileHash,
			fileChunk,
			totalBytesToUpload,
		); err != nil {
			return err
		}
	}

	return nil
}

func (u *Uploader) uploadFileChunkWithProgress(
	sessionId string,
	fileHash []byte,
	fileChunk file.FileChunk,
	totalBytesToUpload int,
) error {
	err := u.uploadRequester.makeUploadFileChunkRequest(
		sessionId,
		fileHash,
		fileChunk.Data,
		fileChunk.ByteRange.StartByte,
		fileChunk.ByteRange.EndByte,
	)
	if err != nil {
		return err
	}

	go func() {
		u.progressChan <- Progress{
			UploadedBytes:      fileChunk.ByteRange.Size(),
			TotalBytesToUpload: totalBytesToUpload,
		}
	}()

	return nil
}
