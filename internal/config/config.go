package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var defaultConfig string

func SetDefaultConfig(content string) {
	defaultConfig = content
}

type Entry struct {
	Title       string
	Description []string
	Command     string
}

func GetConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "ifg")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".ifg")
}

func LoadConfig(dir string) ([]Entry, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory does not exist: %w", err)
	}

	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".sh") {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk config directory: %w", err)
	}

	sort.Strings(files)

	var entries []Entry
	for _, file := range files {
		fileEntries, err := parseFile(filepath.Join(dir, file))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}
		entries = append(entries, fileEntries...)
	}

	return entries, nil
}

func parseFile(path string) ([]Entry, error) {
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

func CreateDefaultConfig(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	defaultFile := filepath.Join(dir, "config.sh")
	if _, err := os.Stat(defaultFile); err == nil {
		return nil
	}

	file, err := os.Create(defaultFile)
	if err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(defaultConfig); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
