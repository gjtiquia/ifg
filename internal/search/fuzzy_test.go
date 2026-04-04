package search

import (
	"testing"

	"github.com/gjtiquia/ifg/internal/config"
)

func TestMatch(t *testing.T) {
	entries := []config.Entry{
		{
			Title:   "copy to clipboard (MacOS)",
			Command: "pbcopy",
		},
		{
			Title:       "git commit with message",
			Description: []string{"commits staged changes", "$ git commit -m \"message\""},
			Command:     "git commit -m \"message\"",
		},
		{
			Title:   "docker ps",
			Command: "docker ps -a",
		},
	}

	t.Run("empty query returns all entries", func(t *testing.T) {
		result := Match(entries, "")
		if len(result) != len(entries) {
			t.Errorf("expected %d entries, got %d", len(entries), len(result))
		}
	})

	t.Run("single token match", func(t *testing.T) {
		result := Match(entries, "git")
		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}
		if result[0].Title != "git commit with message" {
			t.Errorf("expected title %q, got %q", "git commit with message", result[0].Title)
		}
	})

	t.Run("multiple tokens order-agnostic", func(t *testing.T) {
		result1 := Match(entries, "copy macos")
		if len(result1) != 1 {
			t.Fatalf("expected 1 entry for 'copy macos', got %d", len(result1))
		}

		result2 := Match(entries, "macos copy")
		if len(result2) != 1 {
			t.Fatalf("expected 1 entry for 'macos copy', got %d", len(result2))
		}
	})

	t.Run("case-insensitive match", func(t *testing.T) {
		result := Match(entries, "DOCKER")
		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}
		if result[0].Command != "docker ps -a" {
			t.Errorf("expected command %q, got %q", "docker ps -a", result[0].Command)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		result := Match(entries, "nonexistent")
		if len(result) != 0 {
			t.Errorf("expected 0 entries, got %d", len(result))
		}
	})

	t.Run("description match", func(t *testing.T) {
		result := Match(entries, "staged")
		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}
		if result[0].Title != "git commit with message" {
			t.Errorf("expected title %q, got %q", "git commit with message", result[0].Title)
		}
	})
}

func TestScoring(t *testing.T) {
	entries := []config.Entry{
		{
			Title:   "copy file",
			Command: "cp source dest",
		},
		{
			Title:   "copy command",
			Command: "copy",
		},
		{
			Title:   "file operations",
			Command: "file manager",
		},
	}

	t.Run("command match scores higher than title", func(t *testing.T) {
		result := Match(entries, "copy")
		if len(result) < 2 {
			t.Fatalf("expected at least 2 matches, got %d", len(result))
		}

		if result[0].Command != "copy" {
			t.Errorf("expected first result to have 'copy' in command, got %q", result[0].Command)
		}
	})

	t.Run("scores are sorted descending", func(t *testing.T) {
		result := Match(entries, "copy")
		for i := 1; i < len(result); i++ {
			prevScore := matchToken(result[i-1], "copy")
			currScore := matchToken(result[i], "copy")
			if prevScore < currScore {
				t.Errorf("results not properly sorted by score")
			}
		}
	})
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"git commit", []string{"git", "commit"}},
		{"  multiple   spaces  ", []string{"multiple", "spaces"}},
		{"UPPER CASE", []string{"upper", "case"}},
		{"", []string{}},
		{"   ", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tokenize(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(result))
				return
			}
			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("expected token %q, got %q", tt.expected[i], token)
				}
			}
		})
	}
}
