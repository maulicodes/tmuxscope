package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tmux "github.com/jubnzv/go-tmux"
	_ "github.com/mattn/go-sqlite3"
	"tmuxscope/internal/stats"
)

func main() {
	server := tmux.Server{}
	sessions, err := server.GetSessions()
	if err != nil {
		fmt.Println("[tmuxscope] Running outside tmux or server offline.")
	} else {
		fmt.Printf("[tmuxscope] Active tmux sessions: %d\n", len(sessions))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to locate home directory: %v", err)
	}

	dbPath := filepath.Join(home, ".tmuxscope.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	// Ensure schema exists before querying
	schema := `
	CREATE TABLE IF NOT EXISTS command_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT,
		command TEXT,
		category TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}
	analyzer := stats.NewAnalyzer(db)
	summary, err := analyzer.GetSummary(5)
	if err != nil {
		log.Fatalf("failed to fetch stats: %v", err)
	}
	fmt.Println("\n=== tmuxscope Stats Readout ===")
	fmt.Printf("Total Commands: %d\n", summary.TotalCommands)
	fmt.Printf("Total Sessions: %d\n", summary.TotalSessions)
}
