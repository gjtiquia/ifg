# [i] [f]or[g]ot cli

A CLI tool to help you remember commands for tasks you're trying to accomplish.

Dead-simple config format with fuzzy search and vim-style modal editing.

## Features

- **Fuzzy search**: Order-agnostic keyword matching across titles, descriptions, and commands
- **Modal editing**: Insert mode (default) for typing, Normal mode for navigation (vim-style)
- **Simple config**: Plain text config file with comment-based metadata
- **Smart sorting**: Results ranked by relevance (command > title > description)
- **Terminal resize**: Adapts to terminal window changes

## Installation

```bash
# Prerequisites - Go 1.25+ installed
go install github.com/gjtiquia/ifg@latest
```

## Usage

Run `ifg` to enter interactive mode:

```bash
ifg
```

### Insert Mode (default)
- Type keywords to filter commands
- `Backspace` - delete character
- `↑`/`↓` - navigate results
- `Enter` - select command
- `Esc` - switch to Normal mode
- `Ctrl+C` - exit

### Normal Mode
- `j`/`k` or `↑`/`↓` - navigate results
- `i`/`I` - switch to Insert mode (cursor at start)
- `a`/`A` - switch to Insert mode (cursor at end)
- `Enter` - select command
- `Esc` or `Ctrl+C` - exit

### Output

When you select a command with `Enter`, it's printed to stdout. You can then copy-paste or use it as needed.

## Configuration

Config location (checked in order):
1. `$XDG_CONFIG_HOME/ifg/config.sh`
2. `~/.config/ifg/config.sh` (if XDG_CONFIG_HOME not set)
3. `~/.ifg/config.sh` (fallback)

If no config exists, a default one is created automatically.

### Format

```bash
# git commit with message
# Commits staged changes with a message
# $ git commit -m "message"
git commit -m "message"

# copy to clipboard (MacOS)
# Copies text to clipboard
pbcopy

# paste from clipboard (MacOS)
# Pastes clipboard contents
pbpaste
```

**Rules:**
- Entries separated by blank lines
- First `#` line = title
- Subsequent `#` lines = description
- Lines starting with `# $` = usage examples (part of description)
- Last non-comment line = command
- Commands without comments use the command itself as title

### Examples

```bash
# show running containers
docker ps

# find large files
# Lists files > 100M in current directory
find . -size +100M -type f

# View compact git log
# One commit per line with graph
git log --oneline --graph --all
```

## Development

```bash
# Clone the repository
git clone https://github.com/gjtiquia/ifg
cd ifg

# Build
go build -o ifg

# Run tests
go test ./...

# Install locally
go install
```

## Tech Stack

- **Language**: Go 1.25+
- **Dependencies**: `golang.org/x/term` (terminal raw mode)
- **Binary size**: ~2.6MB

## License

MIT


