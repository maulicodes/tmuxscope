package core

import (
	"reflect"
	"testing"
)

func TestStripANSIString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard ansi color code",
			input:    "go build",
			expected: "go build",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripANSIString(tt.input); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
func TestParseCommands(t *testing.T) {
	rawLines := []string{
		"runnnnn!!!",
	}
	expected := []string{
		"runnnnn!!!",
	}
	result := ParseCommands(rawLines)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
