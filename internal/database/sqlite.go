package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// LogEntry represents a single captured command record.
type LogEntry struct {
	ID        int
	Session   string
	Command   string
	Timestamp string
}

// DB wraps the sql.DB connection pool.
type DB struct {
	conn *sql.DB
}

// NewDB opens an SQLite file and initializes the table schema.
func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	query := `
	CREATE TABLE IF NOT EXISTS command_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session TEXT NOT NULL,
		command TEXT NOT NULL,
		captured_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := conn.Exec(query); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	return &DB{conn: conn}, nil
}

// SaveCommands inserts a slice of sanitized commands using a single transaction.
func (d *DB) SaveCommands(session string, commands []string) error {
	if len(commands) == 0 {
		return nil
	}
	tx, err := d.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	stmt, err := tx.Prepare("INSERT INTO command_logs (session, command) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for _, cmd := range commands {
		if _, err := stmt.Exec(session, cmd); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert line: %w", err)
		}
	}
	return tx.Commit()
}

// FetchHistory retrieves the most recent N commands captured.
func (d *DB) FetchHistory(limit int) ([]LogEntry, error) {
	query := `
	SELECT id, session, command, captured_at 
	FROM command_logs 
	ORDER BY id DESC 
	LIMIT ?;`
	rows, err := d.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	var entries []LogEntry
	for rows.Next() {
		var e LogEntry
		if err := rows.Scan(&e.ID, &e.Session, &e.Command, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Close gracefully closes the database handle.
func (d *DB) Close() error {
	return d.conn.Close()
}

// conn returns the underlying *sql.DB handle for analytical packages
func (d *DB) Conn() *sql.DB {
	return d.conn
}
