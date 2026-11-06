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

type PriceRecord struct {
	Ticker    string  `json:"ticker"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
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

func (s *Storage) GetPriceHistory(ticker string) ([]PriceRecord, error) {
	query := `SELECT ticker, price, created_at FROM rates WHERE ticker = ? ORDER BY created_at DESC LIMIT 50`
	rows, err := s.db.Query(query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	defer rows.Close()

	var records []PriceRecord

	for rows.Next() {
		var record PriceRecord

		err := rows.Scan(&record.Ticker, &record.Price, &record.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return records, nil
}
