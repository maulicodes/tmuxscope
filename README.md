# tmuxscope
A lightweight terminal user interface (TUI) dashboard built in Go to capture, parse, and analyze interactive `tmux` session activity.
![tmuxscope Demo](https://github.com/user-attachments/assets/53f94c1b-be48-4222-985e-36925498c49e)
## Key Features
* **Top Commands Analytics**: Automatically aggregates your most frequently executed CLI commands (e.g., `git`, `go`, `vim`, `cargo`) across active `tmux` sessions.
* **Custom Stream Parsing (`internal/core/parser.go`)**: Strips ANSI escape sequences and extracts structured command logs from raw terminal output.
* **Embedded SQLite Storage (`internal/database/sqlite.go`)**: Persists session metrics and command counts locally in `~/.tmuxscope.db`.
* **Interactive TUI**: Built with Charm's `bubbletea` and `lipgloss` for a responsive, keyboard-driven terminal dashboard.
## How Top Commands Tracking Works
1. **Capture**: Listens to pane outputs inside active `tmux` windows.
2. **Parse**: Cleans terminal escape codes and isolates the primary executable command.
3. **Store & Aggregate**: Logs timestamps to SQLite (`~/.tmuxscope.db`) and runs analytical queries to compute total usage frequency.
4. **Display**: Visualizes top command statistics directly in the TUI dashboard.
## Tech Stack
* **Language**: Go (Golang)
* **TUI Framework**: Charm Bubble Tea & Lip Gloss
* **Database**: SQLite3 (`~/.tmuxscope.db`)
## Quickstart
```bash
# Build the binary
go build -o bin/tmuxscope cmd/tmuxscope/main.go

# Run the TUI
./bin/tmuxscope
