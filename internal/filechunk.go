package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func InitiateUploadSession() string {
	sessionId := "123"
	return sessionId
}

func UploadFileChunks(sessionId string, fileChunks [][]byte, url string, uploadedChunksChannel chan<- int) {
	var wg sync.WaitGroup

	wg.Add(len(fileChunks))
	for i := range fileChunks {
		go func(i int) {
			defer wg.Done()
			for {
				err := uploadFileChunk(sessionId, i, fileChunks[i], url)
				if err != nil {
					fmt.Println(err)
					continue
				}

				uploadedChunksChannel <- i
				return
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(uploadedChunksChannel)
	}()
}

func uploadFileChunk(sessionId string, chunkIndex int, fileChunk []byte, url string) error {
	payload := map[string]any{
		"sessionId":  sessionId,
		"chunkIndex": chunkIndex,
		"fileChunk":  fileChunk,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("response status of %d", resp.StatusCode)
	}

	return nil
}
