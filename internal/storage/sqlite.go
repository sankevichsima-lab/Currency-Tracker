package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("problem with openig DB: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS rates(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ticker TEXT NOT NULL,
	price REAL NOT NULL,
	created_at DATETIME NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create rates table: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SavePrice(ticker string, price float64) error {
	query := `INSERT INTO rates	(ticker, price, created_at) VALUES(?, ?, ?)`
	_, err := s.db.Exec(query, ticker, price, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save price: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
