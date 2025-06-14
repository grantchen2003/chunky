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

	status := <-client.UploadStatusChan()
	if upload.UploadStarted != status {
		t.Errorf("no UploadStarted status emitted")
		return
	}

	_, isOpen := <-client.UploadStatusChan()
	if isOpen {
		t.Errorf("UploadStatusChan not closed")
		return
	}

	chunkSize := 1 << 20
	for i := 0; i < len(tempFileData); i += chunkSize {
		progress := <-client.UploadProgressChan()

		expectedProgress := upload.Progress{
			UploadedBytes:      min(len(tempFileData)-i, chunkSize),
			TotalBytesToUpload: len(tempFileData),
		}

		if progress != expectedProgress {
			t.Errorf("Wrong progress emitted")
			return
		}
	}

	_, isOpen = <-client.UploadProgressChan()
	if isOpen {
		t.Errorf("UploadProgressChan not closed")
		return
	}

	result := <-client.UploadResultChan()
	if result != upload.UploadResultSuccess {
		t.Errorf("no UploadResultSuccess result emitted")
		return
	}

	_, isOpen = <-client.UploadResultChan()
	if isOpen {
		t.Errorf("UploadResultChan not closed")
		return
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Error: chunky.db not created")
			return
		}
		t.Errorf("Error getting file info for chunky.db")
		return
	}
}

func Test_UploadWithNonExistentFile(t *testing.T) {
	baseDirPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory path: %v", err)
		return
	}

	client, err := client.NewClient(
		"serverUrl",
		"non-existent-file",
		&upload.Endpoints{
			InitiateUploadSession: "/initiateUploadSession",
			ByteRangesToUpload:    "/byteRangesToUpload",
			UploadFileChunk:       "/uploadFileChunk",
		},
	)
	defer func() {
		dbPath := filepath.Join(baseDirPath, "chunky.db")
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("Error deleting chunky.db: %v", err)
		}
	}()

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

	status := <-client.UploadStatusChan()
	if status != upload.UploadStarted {
		t.Errorf("no UploadStarted status emitted")
		return
	}

	_, isOpen := <-client.UploadStatusChan()
	if isOpen {
		t.Errorf("UploadStatusChan not closed")
		return
	}

	_, isOpen = <-client.UploadProgressChan()
	if isOpen {
		t.Errorf("UploadProgressChan not closed")
		return
	}

	result := <-client.UploadResultChan()
	if result != upload.UploadResultError {
		t.Errorf("no UploadResultError result emitted")
		return
	}

	_, isOpen = <-client.UploadResultChan()
	if isOpen {
		t.Errorf("UploadResultChan not closed")
		return
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Error: chunky.db not created")
			return
		}

		t.Errorf("Error getting file info for chunky.db")
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

	status := <-client.UploadStatusChan()
	if status != upload.UploadStarted {
		t.Errorf("no UploadStarted status emitted")
		return
	}

	_, isOpen := <-client.UploadStatusChan()
	if isOpen {
		t.Errorf("UploadStatusChan not closed")
		return
	}

	_, isOpen = <-client.UploadProgressChan()
	if isOpen {
		t.Errorf("UploadProgressChan not closed")
		return
	}

	result := <-client.UploadResultChan()
	if result != upload.UploadResultSuccess {
		t.Errorf("no UploadResultSuccess result emitted")
		return
	}

	_, isOpen = <-client.UploadResultChan()
	if isOpen {
		t.Errorf("UploadResultChan not closed")
		return
	}

	dbPath := filepath.Join(baseDirPath, "chunky.db")
	_, err = os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Error: chunky.db not created")
			return
		}

		t.Errorf("Error getting file info for chunky.db")
		return
	}
}
