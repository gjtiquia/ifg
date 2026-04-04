package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Title       string
	Description []string
	Command     string
}

func GetConfigPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "ifg", "config.sh")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	xdgPath := filepath.Join(homeDir, ".config", "ifg", "config.sh")
	if _, err := os.Stat(xdgPath); err == nil {
		return xdgPath
	}

	return filepath.Join(homeDir, ".ifg", "config.sh")
}

func LoadConfig(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var entries []Entry
	var currentEntry *Entry
	var inBlock bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			if currentEntry != nil && currentEntry.Command != "" {
				entries = append(entries, *currentEntry)
				currentEntry = nil
				inBlock = false
			}
			continue
		}

		if strings.HasPrefix(line, "#") {
			commentText := strings.TrimSpace(strings.TrimPrefix(line, "#"))

			if !inBlock {
				currentEntry = &Entry{
					Title: commentText,
				}
				inBlock = true
			} else {
				currentEntry.Description = append(currentEntry.Description, commentText)
			}
		} else {
			if !inBlock {
				currentEntry = &Entry{
					Command: strings.TrimSpace(line),
					Title:   strings.TrimSpace(line),
				}
				inBlock = true
			} else {
				if currentEntry.Command == "" {
					currentEntry.Command = strings.TrimSpace(line)
				} else {
					currentEntry.Description = append(currentEntry.Description, strings.TrimSpace(line))
				}
			}
		}
	}

	if currentEntry != nil && currentEntry.Command != "" {
		entries = append(entries, *currentEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return entries, nil
}

func CreateDefaultConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	defaultContent := `# copy to clipboard (MacOS)
# copies text to clipboard
pbcopy

# paste from clipboard (MacOS)
# pastes clipboard contents
pbpaste

# copy to clipboard (Linux)
# copies text to clipboard
xclip -selection clipboard

# git status short
# shows concise git status
git status -s

# git log oneline
# shows one commit per line
git log --oneline

# git branch list
# lists all branches
git branch -a

# docker ps
# lists running containers
docker ps

# docker images
# lists all images
docker images

# find large files
# lists files > 100M in current directory
find . -size +100M -type f
`

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(defaultContent); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
