# Shell Integration with Embedded Scripts

## Overview

Add `--bash` and `--zsh` flags that output shell wrapper code. Wrapper scripts are embedded in the binary using Go's `embed` package.

---

## Implementation

### 1. Create Shell Wrapper Files

**File: `shell/ifg.bash`**
```bash
# ifg - interactive command finder
# Source this file or run: eval "$(ifg --bash)"

ifg() {
    local cmd=$(command ifg)
    if [[ -n "$cmd" ]]; then
        READLINE_LINE="$cmd"
        READLINE_POINT=${#READLINE_LINE}
    fi
}
```

**File: `shell/ifg.zsh`**
```zsh
# ifg - interactive command finder
# Source this file or run: eval "$(ifg --zsh)"

ifg() {
    local cmd=$(command ifg)
    if [[ -n "$cmd" ]]; then
        BUFFER="$cmd"
        CURSOR=${#BUFFER}
    fi
}
```

### 2. Embed in Go Binary

**File: `main.go`**

Add flag handling:

```go
//go:embed shell/ifg.bash
var bashWrapper string

//go:embed shell/ifg.zsh
var zshWrapper string

func main() {
    // Check for shell integration flags
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "--bash":
            fmt.Print(bashWrapper)
            os.Exit(0)
        case "--zsh":
            fmt.Print(zshWrapper)
            os.Exit(0)
        case "--help", "-h":
            fmt.Println("Usage: ifg [--bash|--zsh]")
            fmt.Println()
            fmt.Println("  --bash    Print bash integration code")
            fmt.Println("  --zsh     Print zsh integration code")
            fmt.Println()
            fmt.Println("Add to .bashrc:")
            fmt.Println("  eval \"$(ifg --bash)\"")
            fmt.Println()
            fmt.Println("Add to .zshrc:")
            fmt.Println("  eval \"$(ifg --zsh)\"")
            os.Exit(0)
        }
    }
    
    // ... rest of main() - current TUI logic
}
```

**Note:** Need `import "embed"` at top

### 3. Go Module Updates

The `embed` package requires Go 1.16+ (we have Go 1.25). No additional dependencies needed.

The `//go:embed` directive tells Go to embed the files at compile time.

---

## File Structure

```
ifg/
├── main.go              # Updated with --bash/--zsh flags
├── internal/
│   ├── config/
│   ├── search/
│   └── ui/
├── shell/                # NEW
│   ├── ifg.bash          # Bash wrapper
│   └── ifg.zsh           # Zsh wrapper
├── go.mod
└── README.md
```

---

## Usage

### Installation

```bash
# In .bashrc
eval "$(ifg --bash)"

# In .zshrc
eval "$(ifg --zsh)"
```

### Running

```bash
# Direct (no wrapper) - outputs to stdout
ifg
# Command prints, user can pipe/copy

# With wrapper loaded
ifg
# Command appears in readline, ready to edit/execute
```

---

## Advantages

**Why embed instead of generate strings?**
1. **Easier maintenance** - Edit .sh files with syntax highlighting
2. **Shell-specific logic** - Complex wrappers are easier to read/write as .sh files
3. **Single binary** - No external files, works anywhere
4. **Testability** - Can test shell scripts independently
5. **Future-proof** - Easy to add Fish, PowerShell, etc.

**Why shell wrapper?**
- fzf proven pattern
- Only way to inject into readline
- User controls when to enable (via `eval`)

---

## Testing

```bash
# Test flag output
go build -o ifg
./ifg --bash
# Should output bash wrapper code

./ifg --zsh
# Should output zsh wrapper code

# Test integration
eval "$(./ifg --bash)"
# Should define ifg function

type ifg
# Should show the function definition

# Test full flow (interactive)
./ifg
# Select command → should appear in readline (if wrapper loaded)
```

---

## Development Testing

During development, use the `go build -o` flow to test shell integration:

### 1. Build and Test Flags

```bash
# Build binary
go build -o ifg

# Test flags
./ifg --bash    # Should output bash wrapper code
./ifg --zsh     # Should output zsh wrapper code
```

### 2. Test Integration (Bash)

```bash
# Start clean subshell for testing
bash

# Load wrapper
eval "$(./ifg --bash)"

# Verify function is defined
type ifg
# Should output: ifg is a function
# ifg () 
# { 
#     local cmd=$(command ifg);
#     if [[ -n "$cmd" ]]; then
#         READLINE_LINE="$cmd";
#         READLINE_POINT=${#READLINE_LINE};
#     fi
# }

# Test full flow
ifg
# Select command with Enter
# Command should appear in readline, ready to edit/execute
```

### 3. Test Integration (Zsh)

```bash
# Start clean subshell for testing
zsh

# Load wrapper
eval "$(./ifg --zsh)"

# Verify function is defined
type ifg
# Should show function definition

# Test full flow
ifg
# Select command with Enter
# Command should appear in buffer, ready to edit/execute
```

### 4. Test Without Integration

```bash
# Run binary directly (no wrapper)
./ifg
# Select command → should print to stdout
# NOT injected into readline
```

### Why Not `go run .`?

**DON'T use `alias ifg="go run ."`** - it won't work properly with shell integration because:

- Shell wrapper calls `command ifg` to bypass the wrapper function
- `command` bypasses **functions**, but still respects **aliases**
- So `command ifg` would still run `go run .`
- Each invocation would rebuild and restart the Go program
- Slower and doesn't accurately simulate the real binary

**DO use `go build -o ifg`** for testing:

- Builds once, fast iterations
- Binary exists on disk for `command ifg` to find
- Accurately simulates installed binary behavior

---

## README Updates

Add to README.md:

```markdown
## Shell Integration

To enable command editing (like fzf's Ctrl+R):

**Bash** - Add to `~/.bashrc`:
```bash
eval "$(ifg --bash)"
```

**Zsh** - Add to `~/.zshrc`:
```zsh
eval "$(ifg --zsh)"
```

Now when you run `ifg`, the selected command appears in your prompt for editing before execution.

**Without integration:** `ifg` prints the command to stdout (useful for piping/scripts).
```

---

## Implementation Steps

1. Create `shell/` directory
2. Create `shell/ifg.bash` with wrapper
3. Create `shell/ifg.zsh` with wrapper
4. Update `main.go`:
   - Add `embed` import
   - Add `//go:embed` directives
   - Add flag parsing before TUI logic
5. Update README with integration instructions
6. Build and test

**Time estimate:** 30 minutes

---

## Questions

None - approach is straightforward. Ready to implement when you give the go-ahead.