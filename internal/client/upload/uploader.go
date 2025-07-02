package upload

import (
	"context"
	"os"

	"github.com/grantchen2003/chunky/internal/client/byterange"
	"github.com/grantchen2003/chunky/internal/client/file"
	us "github.com/grantchen2003/chunky/internal/client/upload/uploadstorer"
	"github.com/grantchen2003/chunky/internal/util/workerpool"
)

type Uploader struct {
	url               string
	filePath          string
	maxChunkSizeBytes int
	progressChan      chan<- Progress
	uploadStorer      us.UploadStorer
	uploadRequester   *Requester
	workerPool        *workerpool.WorkerPool
}

func NewUploader(
	url string,
	filePath string,
	maxChunkSizeBytes int,
	maxConcurrentUploads int,
	progressChan chan<- Progress,
	uploadStorer us.UploadStorer,
	uploadRequester *Requester,
) *Uploader {
	return &Uploader{
		url:               url,
		filePath:          filePath,
		maxChunkSizeBytes: maxChunkSizeBytes,
		progressChan:      progressChan,
		uploadStorer:      uploadStorer,
		uploadRequester:   uploadRequester,
		workerPool:        workerpool.NewWorkerPool(maxConcurrentUploads),
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

	fileSizeBytes, err := u.getFileSizeBytes()
	if err != nil {
		return err
	}

	if fileSizeBytes == 0 {
		return nil
	}

	byteRange, err := byterange.NewByteRange(0, fileSizeBytes-1)
	if err != nil {
		return err
	}

	byteRangesToUpload := []byterange.ByteRange{byteRange}

	err = u.streamFileUpload(ctx, sessionId, fileHash, byteRangesToUpload, u.maxChunkSizeBytes)
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

	err = u.streamFileUpload(ctx, sessionId, fileHash, byteRangesToUpload, u.maxChunkSizeBytes)
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
	byteRanges []byterange.ByteRange,
	maxChunkSizeBytes int,
) error {
	bfr, err := file.NewBufferedFileReader(u.filePath)
	if err != nil {
		return err
	}
	defer bfr.Close()

	totalBytesToUpload := byterange.TotalByteCount(byteRanges)

	streamFileUploadErrorChannel := make(chan error)

	go func() {
		for result := range u.workerPool.ResultsChannel() {
			if result.Err != nil {
				streamFileUploadErrorChannel <- result.Err
				return
			}
		}
		streamFileUploadErrorChannel <- nil
	}()

	for fileChunk, err := range bfr.ReadChunkWithRange(maxChunkSizeBytes, byteRanges) {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		u.workerPool.AddJob(*workerpool.NewJob(func() (any, error) {
			return nil, u.uploadFileChunkWithProgress(
				sessionId,
				fileHash,
				fileChunk,
				totalBytesToUpload,
			)
		}))
	}

	u.workerPool.Wait()

	err = <-streamFileUploadErrorChannel
	return err
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
		fileChunk.ByteRange.StartByte(),
		fileChunk.ByteRange.EndByte(),
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
