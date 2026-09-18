package core

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type CaptureOptions struct {
	target    string
	startline int
	endline   int
}

func CapturePane(opts CaptureOptions) ([]string, error) {
	args := []string{"capture-pane", "-p"}

	if opts.target != "" {
		args = append(args, "-t", opts.target)
	}
	if opts.startline != 0 {
		args = append(args, "-S", strconv.Itoa(opts.startline))
	}
	if opts.endline != 0 {
		args = append(args, "-E", strconv.Itoa(opts.endline))
	}

	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to capture tmux: %w", err)
	}
	raw := strings.Split(string(out), "\n")
	return ParseCommands(raw), nil
}
