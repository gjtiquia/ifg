package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("basic entry with title and description", func(t *testing.T) {
		content := `# git commit
# commits staged changes
# $ git commit -m "message"
git commit -m "message"
`
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		entries, err := LoadConfig(tmpFile)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}

		expected := Entry{
			Title:       "git commit",
			Description: []string{"commits staged changes", "$ git commit -m \"message\""},
			Command:     "git commit -m \"message\"",
		}

		if entries[0].Title != expected.Title {
			t.Errorf("expected title %q, got %q", expected.Title, entries[0].Title)
		}
		if len(entries[0].Description) != len(expected.Description) {
			t.Errorf("expected %d description lines, got %d", len(expected.Description), len(entries[0].Description))
		}
		if entries[0].Command != expected.Command {
			t.Errorf("expected command %q, got %q", expected.Command, entries[0].Command)
		}
	})

	t.Run("entry without comments", func(t *testing.T) {
		content := `git status
`
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		entries, err := LoadConfig(tmpFile)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}

		if entries[0].Title != "git status" {
			t.Errorf("expected title %q, got %q", "git status", entries[0].Title)
		}
		if entries[0].Command != "git status" {
			t.Errorf("expected command %q, got %q", "git status", entries[0].Command)
		}
	})

	t.Run("multiple entries", func(t *testing.T) {
		content := `# first command
echo "first"

# second command
echo "second"
`
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		entries, err := LoadConfig(tmpFile)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}

		if entries[0].Title != "first command" {
			t.Errorf("expected title %q, got %q", "first command", entries[0].Title)
		}
		if entries[1].Title != "second command" {
			t.Errorf("expected title %q, got %q", "second command", entries[1].Title)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		content := ``
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		entries, err := LoadConfig(tmpFile)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 0 {
			t.Fatalf("expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("entry with only comment lines (no command)", func(t *testing.T) {
		content := `# this has no command
# just comments

echo "actual command"
`
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		entries, err := LoadConfig(tmpFile)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}

		if entries[0].Command != "echo \"actual command\"" {
			t.Errorf("expected command %q, got %q", "echo \"actual command\"", entries[0].Command)
		}
	})
}

func TestGetConfigPath(t *testing.T) {
	t.Run("XDG_CONFIG_HOME set", func(t *testing.T) {
		origXDG := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXDG)

		os.Setenv("XDG_CONFIG_HOME", "/custom/config")
		path := GetConfigPath()

		expected := filepath.Join("/custom/config", "ifg", "config.sh")
		if path != expected {
			t.Errorf("expected %q, got %q", expected, path)
		}
	})

	t.Run("fallback to ~/.ifg", func(t *testing.T) {
		origXDG := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXDG)
		os.Unsetenv("XDG_CONFIG_HOME")

		homeDir, _ := os.UserHomeDir()
		xdgPath := filepath.Join(homeDir, ".config", "ifg", "config.sh")
		if _, err := os.Stat(xdgPath); err == nil {
			t.Skip("XDG config already exists")
		}

		path := GetConfigPath()
		expected := filepath.Join(homeDir, ".ifg", "config.sh")
		if path != expected {
			t.Errorf("expected %q, got %q", expected, path)
		}
	})
}

func TestCreateDefaultConfig(t *testing.T) {
	t.Run("creates config directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".config", "ifg", "config.sh")

		err := CreateDefaultConfig(configPath)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("config file was not created")
		}
	})

	t.Run("does not overwrite existing config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.sh")

		if err := os.WriteFile(configPath, []byte("existing content"), 0644); err != nil {
			t.Fatal(err)
		}

		err := CreateDefaultConfig(configPath)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "existing content" {
			t.Error("CreateDefaultConfig should not overwrite existing config")
		}
	})
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "ifg-config-test-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpFile.Name()
}
