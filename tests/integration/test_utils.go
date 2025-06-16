package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
