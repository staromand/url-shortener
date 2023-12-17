package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func MigrateNew(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
		    id INTEGER PRIMARY KEY,
		    alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias IN url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while preparing a migration command: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while executing a migration command: %s", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("an error occurred while preparing an insertion command: %s", err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqlErr sqlite3.Error
		sqlErrorMsg := "an error occurred while preparing an insertion command: %s"

		if errors.As(err, &sqlErr) && errors.Is(sqlErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf(sqlErrorMsg, storage.ErrURLExists)
		}

		return 0, fmt.Errorf(sqlErrorMsg, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted id: %s", err)
	}

	return id, nil
}

func (s *Storage) GetUrlByAlias(alias string) (string, error) {
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("an error occurred while selecting url: %s", err)
	}

	var resultUrl string
	err = stmt.QueryRow(alias).Scan(&resultUrl)
	if err != nil {
		var sqlErr sqlite3.Error
		sqlErrorMsg := "an error occurred while selecting url: %s"

		if errors.As(err, &sqlErr) && errors.Is(sqlErr.ExtendedCode, sqlite3.ErrNotFound) {
			return "", fmt.Errorf(sqlErrorMsg, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf(sqlErrorMsg, err)
	}

	return resultUrl, nil
}
