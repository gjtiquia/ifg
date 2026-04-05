package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	t.Run("XDG_CONFIG_HOME set", func(t *testing.T) {
		origXDG := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXDG)

		os.Setenv("XDG_CONFIG_HOME", "/custom/config")
		dir := GetConfigDir()

		expected := filepath.Join("/custom/config", "ifg")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})

	t.Run("fallback to ~/.ifg", func(t *testing.T) {
		origXDG := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXDG)
		os.Unsetenv("XDG_CONFIG_HOME")

		homeDir, _ := os.UserHomeDir()
		xdgPath := filepath.Join(homeDir, ".config", "ifg")
		if _, err := os.Stat(xdgPath); err == nil {
			t.Skip("XDG config already exists")
		}

		dir := GetConfigDir()
		expected := filepath.Join(homeDir, ".ifg")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("basic entry with title and description", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "test.sh", `# git commit
# commits staged changes
# $ git commit -m "message"
git commit -m "message"
`)

		entries, err := LoadConfig(tmpDir)
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
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "test.sh", `git status
`)

		entries, err := LoadConfig(tmpDir)
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
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "test.sh", `# first command
echo "first"

# second command
echo "second"
`)

		entries, err := LoadConfig(tmpDir)
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

	t.Run("empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		entries, err := LoadConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 0 {
			t.Fatalf("expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("entry with only comment lines (no command)", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "test.sh", `# this has no command
# just comments

echo "actual command"
`)

		entries, err := LoadConfig(tmpDir)
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

func TestLoadConfigMultipleFiles(t *testing.T) {
	t.Run("multiple files merged alphabetically", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "b.sh", `# second file
echo "from b"
`)
		writeFile(t, tmpDir, "a.sh", `# first file
echo "from a"
`)

		entries, err := LoadConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}

		if entries[0].Title != "first file" {
			t.Errorf("expected first entry to be from a.sh, got %q", entries[0].Title)
		}
		if entries[1].Title != "second file" {
			t.Errorf("expected second entry to be from b.sh, got %q", entries[1].Title)
		}
	})
}

func TestLoadConfigRecursive(t *testing.T) {
	t.Run("subdirectories included", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "root.sh", `# root file
echo "from root"
`)
		subDir := filepath.Join(tmpDir, "sub")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, subDir, "nested.sh", `# nested file
echo "from nested"
`)

		entries, err := LoadConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}

		if entries[0].Title != "root file" {
			t.Errorf("expected first entry to be from root.sh, got %q", entries[0].Title)
		}
		if entries[1].Title != "nested file" {
			t.Errorf("expected second entry to be from sub/nested.sh, got %q", entries[1].Title)
		}
	})
}

func TestLoadConfigIgnoresNonShFiles(t *testing.T) {
	t.Run("non-sh files ignored", func(t *testing.T) {
		tmpDir := t.TempDir()
		writeFile(t, tmpDir, "valid.sh", `# valid entry
echo "valid"
`)
		writeFile(t, tmpDir, "readme.txt", `# this should be ignored
echo "ignored"
`)
		writeFile(t, tmpDir, "config", `# no extension
echo "also ignored"
`)

		entries, err := LoadConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}

		if entries[0].Title != "valid entry" {
			t.Errorf("expected only valid.sh entry, got %q", entries[0].Title)
		}
	})
}

func TestLoadConfigNonexistentDirectory(t *testing.T) {
	t.Run("nonexistent directory returns error", func(t *testing.T) {
		_, err := LoadConfig("/nonexistent/path/to/config")
		if err == nil {
			t.Error("expected error for nonexistent directory")
		}
	})
}

func TestCreateDefaultConfig(t *testing.T) {
	t.Run("creates config directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".config", "ifg")

		err := CreateDefaultConfig(configDir)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			t.Error("config directory was not created")
		}

		defaultFile := filepath.Join(configDir, "config.sh")
		if _, err := os.Stat(defaultFile); os.IsNotExist(err) {
			t.Error("default config.sh was not created")
		}
	})

	t.Run("does not overwrite existing config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, "ifg")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(err)
		}
		defaultFile := filepath.Join(configDir, "config.sh")
		if err := os.WriteFile(defaultFile, []byte("existing content"), 0644); err != nil {
			t.Fatal(err)
		}

		err := CreateDefaultConfig(configDir)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		content, err := os.ReadFile(defaultFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "existing content" {
			t.Error("CreateDefaultConfig should not overwrite existing config")
		}
	})

	t.Run("does not error if directory exists but no default file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, "ifg")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(err)
		}

		err := CreateDefaultConfig(configDir)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		defaultFile := filepath.Join(configDir, "config.sh")
		if _, err := os.Stat(defaultFile); os.IsNotExist(err) {
			t.Error("default config.sh should have been created")
		}
	})
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	filePath := filepath.Join(dir, name)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
