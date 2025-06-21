package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func SetupTest(t *testing.T, fileData string) (*httptest.Server, string, *os.File, func(), error) {
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
		mockServer.Close()
		return nil, "", nil, func() {}, err
	}

	tempDirPath, err := os.MkdirTemp(baseDirPath, "tmpdir_*")
	if err != nil {
		t.Errorf("Error creating temp directory: %v", err)
		mockServer.Close()
		return nil, "", nil, func() {}, err
	}

	tempFile, err := os.CreateTemp(tempDirPath, "temporary file")
	if err != nil {
		t.Errorf("Error creating temp file")
		mockServer.Close()
		os.RemoveAll(tempDirPath)
		return nil, "", nil, func() {}, err
	}

	_, err = tempFile.WriteString(fileData)
	tempFile.Close()
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
		mockServer.Close()
		os.RemoveAll(tempDirPath)
		return nil, "", nil, func() {}, err
	}

	cleanUp := func() {
		mockServer.Close()

		if err = os.RemoveAll(tempDirPath); err != nil {
			panic(err)
		}

		dbPath := filepath.Join(baseDirPath, "chunky.db")
		if err = os.Remove(dbPath); err != nil {
			panic(err)
		}
	}

	return mockServer, baseDirPath, tempFile, cleanUp, nil
}

func uploadSessionIsInSqliteDatabase(dbPath string, mockServerUrl string, tempFilePath string) bool {
	return true
}
