package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grantchen2003/chunky/internal/byterange"
)

// need to implement
type UploadRequester struct {
	baseUrl   string
	endpoints *UploadEndpoints
}

func NewUploadRequester(baseUrl string, endpoints *UploadEndpoints) *UploadRequester {
	if endpoints == nil {
		endpoints = &UploadEndpoints{
			InitiateUploadSession: "/initiateUploadSession",
			ByteRangesToUpload:    "/byteRangesToUpload",
			UploadFileChunk:       "/uploadFileChunk",
		}
	}

	endpoints.PopulateEmptyFields()

	return &UploadRequester{
		baseUrl:   baseUrl,
		endpoints: endpoints,
	}
}

type InitiateUploadSessionPayload struct {
	FileHash           []byte `json:"fileHahs"`
	TotalFileSizeBytes int    `json:"TotalFileSizeBytes"`
}

type InitiateUploadSessionResponse struct {
	SessionId string
}

func (ur UploadRequester) makeInitiateUploadSessionRequest(fileHash []byte, totalFileSizeBytes int) (string, error) {
	payload := InitiateUploadSessionPayload{FileHash: fileHash, TotalFileSizeBytes: totalFileSizeBytes}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", ur.baseUrl, ur.endpoints.InitiateUploadSession), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response InitiateUploadSessionResponse
	json.NewDecoder(resp.Body).Decode(&response)

	return response.SessionId, nil
}

func (ur UploadRequester) makeByteRangesToUploadRequest(sessionId string, fileHash []byte) ([]byterange.ByteRange, error) {
	fmt.Printf("Getting byte ranges to upload for sessionId: %s and fileHash: %v\n", sessionId, fileHash)

	responseData := [][]int{{6, 100}, {102, 132}, {103, 104}, {133, 152}, {154, 154}, {155, 155}, {156, 159}, {161, 306}}

	var byteRangesToUpload []byterange.ByteRange

	for _, data := range responseData {
		byteRange, err := byterange.NewByteRange(data[0], data[1])
		if err != nil {
			return nil, err
		}
		byteRangesToUpload = append(byteRangesToUpload, byteRange)
	}

	return byteRangesToUpload, nil
}

func (ur UploadRequester) makeUploadFileChunkRequest(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
	// fmt.Printf("Uploading to %s, sessionId: %s, fileHash: %v, startByte: %d, endByte: %d\n", ur.endpoints.UploadFileChunk, sessionId, fileHash, startByte, endByte)
	// time.Sleep(100 * time.Millisecond)
	return nil
}
