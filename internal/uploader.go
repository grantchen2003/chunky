package internal

import (
	"context"
	"log"
	"time"
)

// func UploadStatus() []Range {
// 	return make([]Range, 9)
// }

type UploadResult int

type Uploader struct {
	isUploading bool
}

const (
	UploadResultUnknown UploadResult = iota
	UploadResultSuccess
	UploadResultPaused
	UploadResultError
)

func NewUploader() *Uploader {
	return &Uploader{}
}

func (u *Uploader) IsUploading() bool {
	return u.isUploading
}

func (u *Uploader) HasNoExistingupload() bool {
	return false
}

func (u *Uploader) Upload(ctx context.Context, url string, filePath FilePath) UploadResult {
	u.isUploading = true
	defer func() { u.isUploading = false }()

	doneChan := make(chan struct{})

	// byteRangesToUpload := UploadStatus()

	errorOccured := true

	go func() {
		log.Printf("Uploading %s to %s\n", url, filePath)
		// fmt.Println(byteRangesToUpload)
		time.Sleep(3 * time.Second)

		// if rand.Intn(2) == 0 {
		// 	errorOccured = true
		// }

		close(doneChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return UploadResultPaused
		case <-doneChan:
			if errorOccured {
				return UploadResultError
			} else {
				return UploadResultSuccess
			}
		}
	}
}

// func initiateUploadSession(url string, filePath string, totalBytes int, fileHash string) (string, error) {
// 	payload := map[string]any{
// 		"filePath":   filePath,
// 		"totalBytes": totalBytes,
// 		"fileHash":   fileHash,
// 	}

// 	jsonPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		return "", err
// 	}

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		return "", fmt.Errorf("response status of %d", resp.StatusCode)
// 	}

// 	var result map[string]any
// 	json.NewDecoder(resp.Body).Decode(&result)

// 	if sessionId, ok := result["sessionId"].(string); ok {
// 		return sessionId, nil
// 	} else {
// 		return "", fmt.Errorf("ID not found or not a string")
// 	}
// }

// func uploadFileChunks(sessionId string, fileChunks [][]byte, url string, uploadedChunksChannel chan<- int) {
// 	var wg sync.WaitGroup

// 	wg.Add(len(fileChunks))
// 	for i := range fileChunks {
// 		go func(i int) {
// 			defer wg.Done()
// 			for {
// 				err := uploadFileChunk(sessionId, i, fileChunks[i], url)
// 				if err != nil {
// 					fmt.Println(err)
// 					continue
// 				}

// 				uploadedChunksChannel <- i
// 				return
// 			}
// 		}(i)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(uploadedChunksChannel)
// 	}()
// }

// func uploadFileChunk(sessionId string, chunkIndex int, fileChunk []byte, url string) error {
// 	payload := map[string]any{
// 		"sessionId":  sessionId,
// 		"chunkIndex": chunkIndex,
// 		"fileChunk":  fileChunk,
// 	}

// 	jsonPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("response status of %d", resp.StatusCode)
// 	}

// 	return nil
// }
