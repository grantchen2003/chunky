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

func SetupTest(t *testing.T, fileData string) (*httptest.Server, string, string, *os.File) {
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

	baseDirPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory path: %v", err)
		return nil, "", "", nil
	}

	tempDirPath, err := os.MkdirTemp(baseDirPath, "tmpdir_*")
	if err != nil {
		t.Errorf("Error creating temp directory: %v", err)
		return nil, "", "", nil
	}

	tempFile, err := os.CreateTemp(tempDirPath, "temporary file")
	if err != nil {
		t.Errorf("Error creating temp file")
		return nil, "", "", nil
	}

	_, err = tempFile.WriteString(fileData)
	tempFile.Close()
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
		return nil, "", "", nil
	}

	return mockServer, baseDirPath, tempDirPath, tempFile
}

func Test_Upload(t *testing.T) {
	tempFileData := strings.Repeat("dummy data\n", 100000)
	mockServer, baseDirPath, tempDirPath, tempFile := SetupTest(t, tempFileData)
	defer func() {
		mockServer.Close()

		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("Error removing temp file: %v", err)
		}

		if err := os.Remove(tempDirPath); err != nil {
			t.Errorf("Error removing temp directory: %v", err)
		}

		dbPath := filepath.Join(baseDirPath, "chunky.db")
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("Error deleting chunky.db: %v", err)
		}
	}()

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
		return
	}

	go func() {
		err := client.Upload()
		if err != nil {
			t.Errorf("error uploading")
			return
		}
	}()

	var progressCount int

outer:
	for {
		select {
		case result := <-client.UploadResultChan():
			if result != upload.UploadResultSuccess {
				t.Errorf("no UploadResultSuccess result emitted")
				return
			}

			break outer

		case status := <-client.UploadStatusChan():
			if upload.UploadStarted != status {
				t.Errorf("no UploadStarted status emitted")
				return
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
				return
			}
		}
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Error: chunky.db not created")
			return
		}
		t.Errorf("Error getting file infor for chunky.db")
		return
	}
}
func Test_UploadWithNonExistentFile(t *testing.T) {
	client, err := client.NewClient(
		"serverUrl",
		"non-existent-file",
		&upload.Endpoints{
			InitiateUploadSession: "/initiateUploadSession",
			ByteRangesToUpload:    "/byteRangesToUpload",
			UploadFileChunk:       "/uploadFileChunk",
		},
	)

	if err != nil {
		t.Errorf("error with NewClient emitted")
		return
	}

	go func() {
		err := client.Upload()
		if err != nil {
			t.Errorf("error uploading")
			return
		}
	}()

outer:
	for {
		select {
		case result := <-client.UploadResultChan():
			if result != upload.UploadResultError {
				t.Errorf("no UploadResultError result emitted")
				return
			}

			break outer

		case status := <-client.UploadStatusChan():
			if status != upload.UploadStarted {
				t.Errorf("no UploadStarted status emitted")
				return
			}
		}
	}

	baseDirPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory path: %v", err)
		return
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error: chunky.db created")
		return
	}
}

func Test_UploadWithEmptyFile(t *testing.T) {
	mockServer, baseDirPath, tempDirPath, tempFile := SetupTest(t, "")
	defer func() {
		mockServer.Close()

		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("Error removing temp file: %v", err)
		}

		if err := os.Remove(tempDirPath); err != nil {
			t.Errorf("Error removing temp directory: %v", err)
		}
	}()

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
		return
	}

	go func() {
		err := client.Upload()
		if err != nil {
			t.Errorf("error uploading")
			return
		}
	}()

outer:
	for {
		select {
		case result := <-client.UploadResultChan():
			if result != upload.UploadResultSuccess {
				t.Errorf("no UploadResultSuccess result emitted")
				return
			}

			break outer

		case status := <-client.UploadStatusChan():
			if status != upload.UploadStarted {
				t.Errorf("no UploadStarted status emitted")
				return
			}
		}
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error: chunky.db created")
		return
	}
}
