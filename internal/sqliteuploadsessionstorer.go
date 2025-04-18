package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteUploadSessionStorer struct {
}

func NewSqliteUploadSessionStorer() *SqliteUploadSessionStorer {
	return &SqliteUploadSessionStorer{}
}

func (s *SqliteUploadSessionStorer) Store(sessionId string, filePath string, fileHash []byte) error {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return nil
}
