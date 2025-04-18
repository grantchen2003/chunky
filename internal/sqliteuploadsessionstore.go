package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteUploadSessionStore struct {
	db *sql.DB
}

func NewSqliteUploadSessionStore() (*SqliteUploadSessionStore, error) {
	dbPath := "chunky.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	createStmt := `
	CREATE TABLE IF NOT EXISTS uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		url TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_hash BLOB NOT NULL
	);
	`
	_, err = db.Exec(createStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SqliteUploadSessionStore{db: db}, nil
}

func (s *SqliteUploadSessionStore) Store(sessionId string, url string, filePath string, fileHash []byte) error {
	insertStmt := `
	INSERT INTO uploads (session_id, url, file_path, file_hash)
	VALUES (?, ?, ?, ?);
	`
	_, err := s.db.Exec(insertStmt, sessionId, url, filePath, fileHash)
	if err != nil {
		return fmt.Errorf("failed to insert upload session: %w", err)
	}
	return nil
}

func (s *SqliteUploadSessionStore) GetSessionIdAndHash(url string, filePath string) (sessionId string, fileHash []byte, err error) {
	queryStmt := `SELECT session_id, file_hash FROM uploads WHERE url = ? AND file_path = ?`
	row := s.db.QueryRow(queryStmt, url, filePath)

	err = row.Scan(&sessionId, &fileHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, fmt.Errorf("no session found for file path: %s", filePath)
		}
		return "", nil, fmt.Errorf("failed to retrieve session for file path %s: %w", filePath, err)
	}

	return sessionId, fileHash, nil
}

func (s *SqliteUploadSessionStore) Close() error {
	return s.db.Close()
}
