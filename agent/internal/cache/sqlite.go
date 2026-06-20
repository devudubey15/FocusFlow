package cache

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type OfflineLog struct {
	ID        int
	AppName   string
	Title     string
	URL       string
	StartedAt time.Time
	EndedAt   time.Time
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS activity_buffer (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_name TEXT,
			title TEXT,
			url TEXT,
			started_at DATETIME,
			ended_at DATETIME
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) SaveLog(appName, title, url string, started, ended time.Time) error {
	_, err := s.db.Exec(
		"INSERT INTO activity_buffer (app_name, title, url, started_at, ended_at) VALUES (?, ?, ?, ?, ?)",
		appName, title, url, started, ended,
	)
	return err
}

func (s *SQLiteStore) GetPendingLogs() ([]OfflineLog, error) {
	rows, err := s.db.Query("SELECT id, app_name, title, url, started_at, ended_at FROM activity_buffer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []OfflineLog
	for rows.Next() {
		var l OfflineLog
		if err := rows.Scan(&l.ID, &l.AppName, &l.Title, &l.URL, &l.StartedAt, &l.EndedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (s *SQLiteStore) DeleteLog(id int) error {
	_, err := s.db.Exec("DELETE FROM activity_buffer WHERE id = ?", id)
	return err
}
