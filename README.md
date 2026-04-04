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

When you select a command with `Enter`, it's printed to stdout.

**With shell integration:** The command is added to your shell history. Press UP to access and edit before executing.

**Without integration:** The command prints to stdout (useful for piping/scripts).

## Shell Integration

To enable history integration (works in bash and zsh):

**Add to `~/.bashrc` or `~/.zshrc`:**
```bash
source "$(ifg --sh)"
```

**How it works:**
1. Run `ifg` to select a command
2. Command is added to shell history
3. Message: "Command: <cmd>" and "Press UP to access"
4. Press UP to retrieve command from history
5. Edit if needed, press Enter to execute

**Benefits:**
- Works in bash, zsh, and other shells with `history -s`
- No keybinding required
- UP arrow is universal and intuitive
- Command can be edited before execution

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

### Testing Shell Integration

**IMPORTANT:** Shell integration uses history, which works in both interactive and non-interactive shells.

During development, test with the built binary:

```bash
# Build
go build -o ifg

# Test flag (works everywhere)
./ifg --sh     # Outputs shell wrapper

# Test integration (interactive terminal)
# 1. Add to PATH temporarily
export PATH="$PWD:$PATH"

# 2. Load wrapper
source "$(ifg --sh)"

# 3. Test
ifg
# Select a command
# Message: Command: <cmd>
#          Press UP to access from history
# Press UP to retrieve command

# 4. Verify history
history | tail -1
# Should show the selected command

# Test without integration (no PATH needed)
./ifg
# Select command → prints to stdout only
```

**Note:** Don't use `alias ifg="go run ."` - it won't work properly. The wrapper calls `command ifg` to bypass functions, which still respects aliases and would run `go run .` each time.

## Tech Stack

- **Language**: Go 1.25+
- **Dependencies**: 
  - `github.com/gdamore/tcell/v2` - Terminal UI
  - `golang.org/x/term` - Terminal handling
  - `golang.org/x/text` - Text processing
- **Binary size**: ~3.8MB

## License

MIT


