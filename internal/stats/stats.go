package stats
import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)
//commandstat holds freq for specific commands
type Summary struct {
	TotalCommands int
	TotalSessions int
	TopCommands []CommandStat
	TypeDist map[string]int
}
type Analyzer struct {
	db *sq.DB
}
func NewAnalyzer(db *sql.dB) {
	return &Analyzer{db:db}
}
//categorize command assigns a command string to a functional banter
func CategorizeCommand(cmd string) string {
	trimmed := strings.TrimSpace(cmd)
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return "other"
	}
	root := fields[0]
	switch root {
		case "git"
		return "Git"
	case "go", "cargo" , "make" , "cmake" , "ninja" , "gcc" , "g++" , "clang" :
		return "Build/Dev"
	case "vim", "nvim","vi", "nano","emacs":
		return "editing"
  case "ls", "cd", "cat", "pwd", "mkdir", "rm", "cp", "cp","mv" :
	  return "core utilities"
  case "docker", "podman", "kubectl":
	  return "containers/cloud"
  default: 
  return "other"
	}
}
//getsummary collects metrics and computes the type distribution
func (a *Analyzer) GetSummary(topLimit int) (*Summary, error) {
	summary := &Summary{
		TypeDist: make(map[string]int),
	}
	// 1 total commands
	err := a.db.QueryRow("SELECT COUNT(*) FROM command_logs").Scan(&summary.TotalCommands)
	if err != nil {
		return nil, fmt.Errorf("failed to count commands: %w", err)
	}
// 2 unique sessions
	err = a.db.QueryRow("SELECT COUNT(DISTINCT session) FROM command_logs").Scan(&summary.TotalSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}
	// 3. top commands
	rows, err := a.db.Query(`
		SELECT command, COUNT(*) as frequency
		FROM command_logs
	    	GROUP BY command
		 ORDER BY frequency DESC
		LIMIT ?`, topLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top commands: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cs CommandStat
	if err := rows.Scan(&cs.Command, &cs.Count); err != nil {
			return nil, fmt.Errorf("failed to scan command stat: %w", err)
		}
	summary.TopCommands = append(summary.TopCommands, cs)
	// 4. calculate type distribution across all recorded commands
	allRows, err := a.db.Query("SELECT command FROM command_logs")
	if err != nil {
		return nil, fmt.Errorf("failed to query all commands for typeDist: %w", err)
	}
	defer allRows.Close()
	for allRows.Next() {
		var cmd string
		if err := allRows.Scan(&cmd); err == nil {
			category := CategorizeCommand(cmd)
			summary.TypeDist[category]++
		}
}
	return summary, nil
}
