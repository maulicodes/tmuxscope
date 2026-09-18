package core

import (
	"regexp"
	"strings"
)

var ansiRegex = regexp.MustCompile(`(?i)\x1b\[[0-9;]*[a-z]|[\x80-\x9A\x9C-\x9F]|\x1b[@-Z\\-_]`)

// StripANSI strips away ANSI escape codes from byte slices
func StripANSI(in []byte) []byte {
	return ansiRegex.ReplaceAllLiteral(in, []byte(""))
}

// StripANSIString removes ANSI escape codes from strings
func StripANSIString(in string) string {
	return ansiRegex.ReplaceAllString(in, "")
}

// CleanLine strips ANSI formatting and trims surrounding whitespace of strings
func CleanLine(raw string) string {
	return strings.TrimSpace(StripANSIString(raw))
}

// ParseCommands cleans a slice of raw terminal output lines and removes blank lines
func ParseCommands(rawLines []string) []string {
	var cleaned []string
	for _, line := range rawLines {
		if c := CleanLine(line); c != "" {
			cleaned = append(cleaned, c)
		}
	}
	return cleaned
}
