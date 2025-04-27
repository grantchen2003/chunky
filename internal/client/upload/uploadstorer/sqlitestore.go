package storer

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // or your DB driver
)

type SqliteUploadStore struct {
	db *sql.DB
}

func NewSqliteUploadStore() (*SqliteUploadStore, error) {
	db, err := sql.Open("sqlite3", "chunky.db")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS upload_sessions (
		url TEXT NOT NULL,
		file_path TEXT NOT NULL,
		session_id TEXT NOT NULL,
		file_hash BLOB,
		PRIMARY KEY (url, file_path)
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &SqliteUploadStore{db: db}, nil
}

func (s *SqliteUploadStore) Store(sessionId string, url string, filePath string, fileHash []byte) error {
	query := `
	INSERT INTO upload_sessions (session_id, url, file_path, file_hash)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(url, file_path) DO UPDATE SET
		session_id = excluded.session_id,
		file_hash = excluded.file_hash;
	`
	_, err := s.db.Exec(query, sessionId, url, filePath, fileHash)
	return err
}

func (s *SqliteUploadStore) GetSessionIdAndFileHash(url string, filePath string) (string, []byte, error) {
	query := `SELECT session_id, file_hash FROM upload_sessions WHERE url = ? AND file_path = ?`
	row := s.db.QueryRow(query, url, filePath)

	var sessionId string
	var fileHash []byte
	err := row.Scan(&sessionId, &fileHash)
	if err == sql.ErrNoRows {
		return "", nil, ErrNotFound
	}
	return sessionId, fileHash, err
}

func (s *SqliteUploadStore) Close() error {
	return s.db.Close()
}
