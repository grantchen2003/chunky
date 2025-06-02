package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantchen2003/chunky/internal/client/upload"
	"github.com/grantchen2003/chunky/pkg/client"
)

// test upload with empty file, file doesn't exist

// no db alr
func Test_Upload(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/initiateUploadSession" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{"sessionId": "dummySessionId"}); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				log.Printf("Error encoding response: %v", err)
			}
			return
		}

		if r.URL.Path == "/uploadFileChunk" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
	defer mockServer.Close()

	baseDirPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory path: %v", err)
		return
	}

	tempDirPath, err := os.MkdirTemp(baseDirPath, "tmpdir_*")
	if err != nil {
		t.Errorf("Error creating temp directory: %v", err)
		return
	}
	defer func() {
		if err = os.RemoveAll(tempDirPath); err != nil {
			t.Errorf("Error removing temp directory: %v", err)
		}
	}()

	tempFile, err := os.CreateTemp(tempDirPath, "temporary file")
	if err != nil {
		t.Errorf("Error creating temp file")
	}

	tempFileData := strings.Repeat("dummy data\n", 100000)
	_, err = tempFile.WriteString(tempFileData)
	defer tempFile.Close()
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
		return
	}

	client, err := client.NewClient(
		mockServer.URL,
		tempFile.Name(),
		&upload.Endpoints{
			InitiateUploadSession: "/initiateUploadSession",
			ByteRangesToUpload:    "/byteRangesToUpload",
			UploadFileChunk:       "/uploadFileChunk",
		},
	)

	if err != nil {
		t.Errorf("error with NewClient emitted")
	}

	go func() {
		err := client.Upload()
		if err != nil {
			t.Errorf("error uploading")
		}
	}()

	var progressCount int

	for {
		select {
		case result := <-client.UploadResultChan():
			if result != upload.UploadResultSuccess {
				t.Errorf("no UploadResultSuccess result emitted")
			}
			dbPath := filepath.Join(baseDirPath, "chunky.db")
			if err := os.Remove(dbPath); err != nil {
				t.Errorf("Error deleting chunky.db: %v", err)
			}

			return

		case status := <-client.UploadStatusChan():
			if upload.UploadStarted != status {
				t.Errorf("no UploadStarted status emitted")
			}

		case progress := <-client.UploadProgressChan():
			progressCount++
			chunkSize := 1 << 20

			expectedProgress := upload.Progress{
				UploadedBytes:      chunkSize,
				TotalBytesToUpload: len(tempFileData),
			}

			if chunkSize*progressCount > len(tempFileData) {
				expectedProgress.UploadedBytes = len(tempFileData) % chunkSize
			}

			if progress != expectedProgress {
				t.Errorf("Wrong progress emitted")
			}
		}
	}
}
