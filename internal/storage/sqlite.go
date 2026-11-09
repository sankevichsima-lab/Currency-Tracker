package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Storage manages database operations for price records.
// Storage управляет операциями с базой данных для записей цен.
type Storage struct {
	db *sql.DB
}

// PriceRecord represents a single price entry with timestamp.
// PriceRecord представляет одну запись котировки с временной меткой.
type PriceRecord struct {
	Ticker    string  `json:"ticker"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

// NewStorage opens SQLite database and creates the rates table if it does not exist.
// NewStorage открывает базу SQLite и создает таблицу rates, если она отсутствует.
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

	users := `CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create rates table: %w", err)
	}

	if _, err = db.Exec(users); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &Storage{db: db}, nil
}

// SavePrice inserts a new currency price record into the database.
// SavePrice сохраняет новую запись курса валюты в базу данных.
func (s *Storage) SavePrice(ticker string, price float64) error {
	query := `INSERT INTO rates	(ticker, price, created_at) VALUES(?, ?, ?)`
	_, err := s.db.Exec(query, ticker, price, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save price: %w", err)
	}
	return nil
}

// Close gracefully closes the database connection pool.
// Close безопасно закрывает соединение с базой данных.
func (s *Storage) Close() error {
	return s.db.Close()
}

// GetPriceHistory retrieves the last 50 price records for a given ticker.
// GetPriceHistory возвращает последние 50 записей котировок для указанного тикера.
func (s *Storage) GetPriceHistory(ticker string) ([]PriceRecord, error) {
	query := `SELECT ticker, price, created_at FROM rates WHERE ticker = ? ORDER BY created_at DESC LIMIT 50`
	rows, err := s.db.Query(query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	defer rows.Close()

	var records []PriceRecord

	// Scan query results into slice of records.
	// Сканируем строки из результата запроса в срез записей.
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

func (s *Storage) CreateUser(username, passwordHash string) error {
	query := `INSERT INTO users(username, password_hash) VALUES(?, ?)`
	if _, err := s.db.Exec(query, username, passwordHash); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Storage) GetUserHash(username string) (string, error) {
	query := `SELECT password_hash FROM users WHERE username = ?`
	var hash string
	if err := s.db.QueryRow(query, username).Scan(&hash); err != nil {
		return "", fmt.Errorf("username not found: %w", err)
	}
	return hash, nil
}
