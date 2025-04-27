package database

import (
	"database/sql"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite() (*Sqlite, error) {
	db, err := sql.Open("sqlite3", "chunky.db")
	if err != nil {
		return nil, err
	}

	s := &Sqlite{db: db}

	if err := s.createUploadSessionTableIfNotExists(); err != nil {
		return nil, err
	}

	if err := s.createUploadedFileChunksIfNotExists(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sqlite) CreateUploadSession(sessionId string, fileHash []byte, totalFileSizeBytes int) error {
	query := `
	INSERT INTO upload_sessions (session_id, file_hash, total_file_size_bytes)
	VALUES (?, ?, ?)
	`
	_, err := s.db.Exec(query, sessionId, fileHash, totalFileSizeBytes)

	return err
}

func (s *Sqlite) Exists(sessionId string, fileHash []byte) (exists bool, err error) {
	query := `
	SELECT COUNT(*) FROM upload_sessions
	WHERE session_id = ? AND file_hash = ?
	`
	var count int

	if err = s.db.QueryRow(query, sessionId, fileHash).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Sqlite) AddFileChunk(sessionId string, fileHash []byte, chunkId string, startByte int, endByte int) error {
	query := `
	INSERT INTO uploaded_file_chunks (chunk_id, session_id, file_hash, start_byte, end_byte)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, chunkId, sessionId, fileHash, startByte, endByte)
	return err
}

func (s *Sqlite) ByteRangesToUpload(sessionId string, fileHash []byte) ([][2]int, error) {
	query := `
		SELECT start_byte, end_byte 
		FROM uploaded_file_chunks
		WHERE session_id = ? AND file_hash = ?
		ORDER BY start_byte
	`
	rows, err := s.db.Query(query, sessionId, fileHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploadedChunks [][2]int
	for rows.Next() {
		var start, end int
		if err := rows.Scan(&start, &end); err != nil {
			return nil, err
		}
		uploadedChunks = append(uploadedChunks, [2]int{start, end})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var missingRanges [][2]int
	var currentByte int

	var totalSize int
	sizeQuery := `
		SELECT total_file_size_bytes 
		FROM upload_sessions 
		WHERE session_id = ? AND file_hash = ?
	`
	if err := s.db.QueryRow(sizeQuery, sessionId, fileHash).Scan(&totalSize); err != nil {
		return nil, err
	}

	for _, chunk := range uploadedChunks {
		if currentByte < chunk[0] {
			missingRanges = append(missingRanges, [2]int{currentByte, chunk[0] - 1})
		}
		currentByte = chunk[1] + 1
	}

	if currentByte < totalSize {
		missingRanges = append(missingRanges, [2]int{currentByte, totalSize - 1})
	}

	return missingRanges, nil
}

func (s *Sqlite) createUploadSessionTableIfNotExists() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS upload_sessions (
		session_id TEXT NOT NULL,
		file_hash BLOB NOT NULL,
		total_file_size_bytes INTEGER NOT NULL,
		PRIMARY KEY (session_id, file_hash)
	);`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) createUploadedFileChunksIfNotExists() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS uploaded_file_chunks (
		chunk_id TEXT NOT NULL PRIMARY KEY,
		session_id TEXT NOT NULL,
		file_hash BLOB NOT NULL,
		start_byte INTEGER NOT NULL,
		end_byte INTEGER NOT NULL
	);`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
