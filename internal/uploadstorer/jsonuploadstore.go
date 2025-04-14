package uploadstorer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type uploadSession struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
	FilePath  string `json:"file_path"`
	FileHash  []byte `json:"file_hash"`
}

type JsonUploadStore struct {
	sessions map[string]uploadSession // key: url|filePath
	filePath string
	mutex    sync.Mutex
}

func NewJsonUploadStore() (*JsonUploadStore, error) {
	store := &JsonUploadStore{
		sessions: make(map[string]uploadSession),
		filePath: "chunky.json",
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *JsonUploadStore) Store(sessionId string, url string, filePath string, fileHash []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%s|%s", url, filePath)
	s.sessions[key] = uploadSession{
		SessionID: sessionId,
		URL:       url,
		FilePath:  filePath,
		FileHash:  fileHash,
	}

	return s.save()
}

func (s *JsonUploadStore) GetSessionIdAndFileHash(url string, filePath string) (string, []byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%s|%s", url, filePath)
	session, ok := s.sessions[key]
	if !ok {
		return "", nil, fmt.Errorf("no session found for file path: %s", filePath)
	}
	return session.SessionID, session.FileHash, nil
}

func (s *JsonUploadStore) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.save()
}

func (s *JsonUploadStore) load() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist; treat as empty store
		}
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&s.sessions)
}

func (s *JsonUploadStore) save() error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.sessions)
}
