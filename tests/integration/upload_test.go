package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantchen2003/chunky/internal/client/upload"
	"github.com/grantchen2003/chunky/pkg/client"
)

func Test_Upload(t *testing.T) {
	tempFileData := strings.Repeat("dummy data\n", 100000)
	mockServer, baseDirPath, tempFile, cleanUp, err := SetupTest(t, tempFileData)
	if err != nil {
		panic(err)
	}
	defer cleanUp()

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
	tempFileData := ""
	mockServer, baseDirPath, tempFile, cleanUp, err := SetupTest(t, tempFileData)
	if err != nil {
		panic(err)
	}
	defer cleanUp()

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

func Test_UploadIsNotBlockedByStatusChannelRead(t *testing.T) {
	tempFileData := strings.Repeat("dummy data\n", 100000)
	mockServer, baseDirPath, tempFile, cleanUp, err := SetupTest(t, tempFileData)
	if err != nil {
		panic(err)
	}
	defer cleanUp()

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

	_, isOpen := <-client.UploadProgressChan()
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

	status := <-client.UploadStatusChan()
	if upload.UploadStarted != status {
		t.Errorf("no UploadStarted status emitted")
		return
	}

	_, isOpen = <-client.UploadStatusChan()
	if isOpen {
		t.Errorf("UploadStatusChan not closed")
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
