package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grantchen2003/chunky/internal/client/byterange"
)

type Requester struct {
	baseUrl   string
	endpoints *Endpoints
}

func NewRequester(baseUrl string, endpoints *Endpoints) *Requester {
	if endpoints == nil {
		endpoints = &Endpoints{
			InitiateUploadSession: "/initiateUploadSession",
			ByteRangesToUpload:    "/byteRangesToUpload",
			UploadFileChunk:       "/uploadFileChunk",
		}
	}

	endpoints.PopulateEmptyFields()

	return &Requester{
		baseUrl:   baseUrl,
		endpoints: endpoints,
	}
}

func (r Requester) makeInitiateUploadSessionRequest(fileHash []byte, totalFileSizeBytes int) (string, error) {
	type Payload struct {
		FileHash           []byte `json:"fileHahs"`
		TotalFileSizeBytes int    `json:"TotalFileSizeBytes"`
	}

	payload := Payload{FileHash: fileHash, TotalFileSizeBytes: totalFileSizeBytes}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", r.baseUrl, r.endpoints.InitiateUploadSession), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type Response struct {
		SessionId string `json:"sessionId"`
	}

	var response Response

	json.NewDecoder(resp.Body).Decode(&response)

	return response.SessionId, nil
}

func (r Requester) makeByteRangesToUploadRequest(sessionId string, fileHash []byte) ([]byterange.ByteRange, error) {
	type Payload struct {
		SessionId string `json:"sessionId"`
		FileHash  []byte `json:"fileHash"`
	}

	payload := Payload{SessionId: sessionId, FileHash: fileHash}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", r.baseUrl, r.endpoints.ByteRangesToUpload), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Response struct {
		ByteRanges [][2]int `json:"byteRanges"`
	}

	var response Response
	json.NewDecoder(resp.Body).Decode(&response)

	var byteRangesToUpload []byterange.ByteRange
	for _, data := range response.ByteRanges {
		byteRange, err := byterange.NewByteRange(data[0], data[1])
		if err != nil {
			return nil, err
		}
		byteRangesToUpload = append(byteRangesToUpload, byteRange)
	}

	return byteRangesToUpload, nil
}

func (r Requester) makeUploadFileChunkRequest(sessionId string, fileHash []byte, chunk []byte, startByte int, endByte int) error {
	type Payload struct {
		SessionId string `json:"sessionId"`
		FileHash  []byte `json:"fileHash"`
		Chunk     []byte `json:"chunk"`
		StartByte int    `json:"startByte"`
		EndByte   int    `json:"endByte"`
	}

	payload := Payload{
		SessionId: sessionId,
		FileHash:  fileHash,
		Chunk:     chunk,
		StartByte: startByte,
		EndByte:   endByte,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", r.baseUrl, r.endpoints.UploadFileChunk), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
}
