# Implementation Plan: ifg CLI

## Overview

Interactive CLI tool for remembering shell commands. Modal interface (Insert/Normal modes) similar to vim, with fuzzy search across a simple config file.

---

## Architecture

```
ifg/
├── cmd/
│   └── ifg/
│       └── main.go              # Entry point, wires everything together
├── internal/
│   ├── config/
│   │   ├── config.go           # Entry struct, parser
│   │   └── config_test.go      # Unit tests for parser
│   ├── search/
│   │   ├── fuzzy.go            # Fuzzy matching logic
│   │   └── fuzzy_test.go       # Unit tests for search
│   └── ui/
│       ├── tui.go              # Terminal UI rendering
│       ├── input.go            # Keyboard input handling
│       └── state.go            # State management (mode, cursor, etc.)
├── go.mod
├── go.sum
└── README.md
```

---

## Technical Details

### 1. Config Parser (`internal/config/config.go`)

**Config file location:**
- Primary: `$XDG_CONFIG_HOME/ifg/config.sh` (or `~/.config/ifg/config.sh` if XDG_CONFIG_HOME not set)
- Fallback: `~/.ifg/config.sh`
- Check XDG location first, fall back to `~/.ifg/` if not found

**Entry struct:**
```go
type Entry struct {
    Title       string
    Description []string
    Command     string
}
```

**Parsing rules:**
- Entries separated by empty newline (`\n\n`)
- First `#` comment line → Title (strip leading `#` and whitespace)
- Subsequent `#` comment lines → Description (strip leading `#` and whitespace)
- Last non-comment line in a block → Command
- Lines starting with `#` but containing `$` (e.g., `# $ command example`) are usage examples, also part of description

**Example parsing:**

Input:
```bash
# git commit with message
# Commits staged changes with a message
# $ git commit -m "message"
git commit -m "message"
```

Output:
```go
Entry{
    Title: "git commit with message",
    Description: []string{
        "Commits staged changes with a message",
        "$ git commit -m \"message\"",
    },
    Command: "git commit -m \"message\"",
}
```

**Default config creation:**
If `~/.ifg/config.sh` doesn't exist, create it with example entries:
- copy to clipboard (MacOS)
- paste from clipboard (MacOS)
- copy to clipboard (Linux)
- git status short
- git log oneline
- etc.

---

### 2. Fuzzy Search (`internal/search/fuzzy.go`)

**Matching algorithm:**
1. Tokenize search query by whitespace (e.g., "macos copy" → ["macos", "copy"])
2. For each entry, check if ALL tokens match (order-agnostic)
3. Match is case-insensitive
4. Search across: Title, Description lines, Command
5. Token matches if it's a substring of any field

**Scoring (for sorting):**
- Preferred order: Command match > Title match > Description match
- Earlier matches rank higher
- Exact word matches rank higher than substring matches

**Return:** `[]Entry` sorted by relevance score

**Implementation notes:**
- Keep it simple for MVP - substring matching is sufficient
- Scoring can be basic: first match wins, no complex ranking needed

---

### 3. State Management (`internal/ui/state.go`)

**State struct:**
```go
type Mode int

const (
    ModeInsert Mode = iota
    ModeNormal
)

type State struct {
    Mode        Mode
    SearchBuf   string     // User's search input
    CursorIdx   int        // Cursor position in SearchBuf
    SelectedIdx int        // Index in filtered list
    Entries     []config.Entry
    Filtered    []config.Entry  // Entries matching SearchBuf
    TerminalHeight int
    TerminalWidth  int
}
```

**State transitions:**
- Initial state: ModeInsert, empty SearchBuf, SelectedIdx=0
- Typing in Insert mode: append to SearchBuf, update Filtered
- Switching to Normal: Mode = ModeNormal, CursorIdx preserved
- Switching to Insert: 
  - `i/I` → ModeInsert, CursorIdx = 0
  - `a/A` → ModeInsert, CursorIdx = len(SearchBuf)

---

### 4. Terminal UI (`internal/ui/tui.go`)

**Dependencies:**
- `golang.org/x/term` for raw mode and terminal size detection

**Terminal setup:**
1. Save original terminal state
2. Enter raw mode (disable echo, line buffering)
3. Hide cursor (optional)
4. Clear screen and move cursor to top-left

**Rendering loop:**
1. Clear screen
2. Draw search prompt at top: `search: <SearchBuf>`
3. Draw filtered entries (or all entries if SearchBuf empty)
4. Highlight selected entry with `> ` prefix
5. Draw scroll indicator if list > terminal height
6. Flush to stdout

**Drawing format:**
```
search: macos copy
> copy to clipboard (MacOS)
  paste from clipboard (MacOS)
  copy to clipboard (Linux)
```

**Scrolling:**
- If list exceeds terminal height, only render visible portion
- Track scroll offset in state
- Adjust offset when navigating near edges

**Terminal restoration:**
- On exit (or crash), restore original terminal state

---

### 5. Input Handling (`internal/ui/input.go`)

**Keyboard input (raw mode):**

**Insert mode:**
- Printable chars → append to SearchBuf, re-filter
- `Backspace` (127 or `\b`) → delete char before CursorIdx
- `↑`/`↓` arrow keys → navigate Filtered list
- `Enter` → select Filtered[SelectedIdx], exit
- `Esc` → switch to ModeNormal
- `Ctrl+C` → exit without selection

**Normal mode:**
- `j`/`↓` → move selection down
- `k`/`↑` → move selection up
- `Enter` → select Filtered[SelectedIdx], exit
- `i`/`I` → switch to Insert, CursorIdx = 0
- `a`/`A` → switch to Insert, CursorIdx = len(SearchBuf)
- `Esc` → exit without selection
- `Ctrl+C` → exit without selection

**Arrow key detection:**
- Arrow keys send escape sequences: `Esc [ A` (up), `Esc [ B` (down), etc.
- Need to parse these sequences properly

---

### 6. Main Entry Point (`cmd/ifg/main.go`)

**Flow:**
1. Load config from `~/.ifg/config.sh`
   - If missing, create default config
2. Initialize UI state with all entries
3. Enter input loop
4. On selection, print command to stdout
5. Restore terminal and exit

**Exit codes:**
- 0: Success (command selected and printed)
- 1: No selection (user cancelled)
- 2: Error (config error, terminal error, etc.)

**Output:**
- Selected command printed to stdout
- User can then edit and execute in shell:

---

## Implementation Phases

### Phase 1: Foundation
1. Initialize Go module: `go mod init github.com/gjtiquia/ifg`
2. Create project structure
3. Implement `internal/config/config.go`:
   - `Entry` struct
   - `LoadConfig(path string) ([]Entry, error)`
   - `GetConfigPath() string` - returns XDG path or fallback
   - `CreateDefaultConfig(path string) error`
4. Write tests for config parsing

### Phase 2: Search
5. Implement `internal/search/fuzzy.go`:
   - `Match(entries []config.Entry, query string) []config.Entry`
   - Case-insensitive substring matching
   - Order-agnostic token matching
6. Write tests for fuzzy matching

### Phase 3: State Management
7. Implement `internal/ui/state.go`:
   - `State` struct
   - `NewState(entries []config.Entry) *State`
   - `UpdateSearch(input string)`
   - `NavigateUp()`, `NavigateDown()`
   - `SwitchMode()`

### Phase 4: Terminal UI
8. Implement `internal/ui/tui.go`:
   - `SetupTerminal() error`
   - `RestoreTerminal()`
   - `Render(state *State) error`
   - `ClearScreen()`
9. Implement `internal/ui/input.go`:
   - `ReadKey() (Key, error)`
   - Key type enum for special keys
   - Parse arrow key escape sequences

### Phase 5: Main Integration
10. Implement `cmd/ifg/main.go`:
    - Load config or create default
    - Setup terminal
    - Input loop with mode switching
    - Print selected command
    - Handle exit codes
11. Add signal handling (SIGINT, SIGTERM) for clean exit

### Phase 6: Polish
12. Add scroll handling for long lists
13. Handle terminal resize (SIGWINCH)
14. Add usage documentation to README
15. Add shell integration examples

---

## Edge Cases

### Config Parsing
- Empty file → show default entries
- Malformed entries (no command line) → skip and warn
- Duplicate entries → keep both (allow user to have same command with different titles)
- Comments without command → skip
- Entries with only command (no comments) → use command as title, empty description

### Search
- Empty search → show all entries
- No matches → display "No results found" message
- Search with special characters → escape or ignore special regex chars (keep simple)

### UI
- Terminal too small → display "Terminal too small" and exit gracefully
- Binary data in config → validate and filter
- Unicode in commands/titles → handle properly (golang.org/x/term handles this)

### Input
- Invalid UTF-8 input → replace with � or ignore
- Very long search query → truncate display, but keep searching
- Rapid input → handle properly (no flickering)

---

## Testing Strategy

### Unit Tests
- Config parser tests (various formats, edge cases)
- Fuzzy search tests (matching, scoring, order-agnostic)
- State transition tests

---

## Dependencies

**Required:**
- `golang.org/x/term` - Terminal raw mode, size detection

**No external dependencies for:**
- Config parsing (stdlib `bufio`, `os`, `strings`)
- Fuzzy search (stdlib `strings`, `sort`)
- UI rendering (raw escape sequences via `fmt`)

**Why minimal dependencies:**
- Keep binary small
- Reduce attack surface
- Easier to build and distribute

---

## Binary Size Target

- Aim for < 5MB binary size
- Static linking for easy distribution
- Support: Linux (amd64, arm64), MacOS (amd64, arm64), Windows (amd64)

---

## Success Criteria

1. **Works smoothly with 100+ entries** - no lag
2. **Intuitive for vim users** - modal editing feels natural
3. **Easy for non-vim users** - arrow keys work, discoverable
4. **Reliable config parsing** - handles edge cases gracefully
5. **Clean exit** - terminal always restored on exit/crash
6. **Simple installation** - `go install` or download binary

---

## Notes

- Keep it simple for MVP - avoid feature creep
- Focus on core UX: fast, responsive, intuitive
