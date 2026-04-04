# Shell Integration with History Approach

## Overview

Simplified shell integration using `history -s` instead of READLINE_LINE. When a command is selected, it's added to shell history and printed to stdout. The user can then press UP to access and edit it.

This approach is simpler, more portable, and works in ALL shells (bash, zsh, fish, etc.).

---

## Why History Instead of READLINE_LINE?

### The READLINE_LINE Problem

`READLINE_LINE` (bash) and `BUFFER` (zsh) **only work with `bind -x`**, not direct function calls:

```bash
# This DOESN'T work - READLINE_LINE has no effect
ifg() {
    local cmd=$(command ifg)
    READLINE_LINE="$cmd"
}
ifg

# This WORKS - but requires keybinding
bind -x '"\C-g": ifg-widget'
# User must press Ctrl+G, cannot type `ifg`
```

**The key insight:** READLINE_LINE only works when the function is bound to a key with `bind -x`. It does NOT work when called as a regular command.

### The History Solution

`history -s` adds a command to shell history, accessible via UP arrow:

```bash
# Works in bash, zsh, fish, etc.
ifg() {
    local cmd=$(command ifg)
    [[ -n "$cmd" ]] && history -s "$cmd"
}
```

**Advantages:**
- ✅ Works in ALL shells (bash, zsh, fish)
- ✅ No keybinding required - user types `ifg`
- ✅ UP arrow is universal muscle memory
- ✅ Command is editable before execution
- ✅ Simple - no interactive checks, no bind -x
- ✅ Still pipeable: `cmd=$(ifg)`

**Disadvantage:**
- Requires UP arrow (one extra keypress) vs direct readline injection

**Trade-off is worth it** - one UP arrow for massive simplification and portability.

---

## Implementation

### 1. Create Unified Shell Wrapper

**File: `shell/ifg.sh`**

Single wrapper that works in bash AND zsh:

```bash
# ifg - interactive command finder
# Add to ~/.bashrc OR ~/.zshrc:
#   source "$(ifg --sh)"
# Or:
#   eval "$(ifg --sh)"

ifg() {
    local cmd=$(command ifg)
    if [[ -n "$cmd" ]]; then
        history -s "$cmd"
        echo "Command: $cmd"
        echo "Press UP to access from history"
    fi
}
```

**Key features:**
- Single wrapper for bash + zsh (both support `history -s`)
- Adds command to history with `history -s`
- Prints command to stdout so user sees it
- Prints helpful message: "Press UP to access"
- No interactive checks needed

### 2. Embed in Go Binary

**File: `main.go`**

Replace `--bash` and `--zsh` with single `--sh` flag:

```go
//go:embed shell/ifg.sh
var shellWrapper string

func main() {
    // Check for shell integration flags
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "--sh":
            fmt.Print(shellWrapper)
            os.Exit(0)
        case "--help", "-h":
            // ... help text
        }
    }
    
    // ... existing TUI logic
}
```

### 3. Update Binary Behavior

**Option:** Should the binary itself print the message?

**Current:**
- Binary prints command to stdout
- Wrapper adds to history + prints message

**Alternative:**
- Binary prints: `Command: <cmd>\nAdded to history. Press UP to access.`
- Wrapper simplifies to just: `local cmd=$(command ifg); [[ -n "$cmd" ]] && history -s "$cmd"`

**Decision:** Keep message in wrapper - keeps binary output clean (useful for piping).

---

## User Experience Flow

### With Wrapper Loaded

```bash
$ ifg
# User selects "git status -s"
Command: git status -s
Press UP to access from history

$ # User presses UP
$ git status -s  # Command from history, editable
```

### Without Wrapper (Direct Use)

```bash
$ ifg
# User selects "git status -s"
git status -s  # Just outputs to stdout

$ # User can pipe/capture
$ cmd=$(ifg)
# Select command
$ echo "$cmd"
git status -s
```

---

## File Structure

```
ifg/
├── main.go           # Updated with --sh flag
├── internal/
│   ├── config/
│   ├── search/
│   └── ui/
├── shell/
│   └── ifg.sh        # Unified wrapper (bash + zsh)
├── go.mod
└── README.md
```

**Deletions:**
- `shell/ifg.bash` - replaced by unified `ifg.sh`
- `shell/ifg.zsh` - replaced by unified `ifg.sh`

---

## Usage

### Installation

```bash
# Bash - add to ~/.bashrc
source "$(ifg --sh)"

# Zsh - add to ~/.zshrc
source "$(ifg --sh)"

# Or manually add:
# ifg() {
#     local cmd=$(command ifg)
#     if [[ -n "$cmd" ]]; then
#         history -s "$cmd"
#         echo "Command: $cmd"
#         echo "Press UP to access from history"
#     fi
# }
```

### Running

```bash
# With wrapper
ifg
# Select command
# Output: Command: <cmd>
#         Press UP to access from history
# Press UP to get command from history

# Without wrapper
./ifg
# Select command
# Output: <cmd>
```

---

## Advantages

### Why History Approach?

1. **Universal:** Works in bash, zsh, fish, dash, etc.
2. **Simple:** No shell-specific logic (READLINE_LINE vs BUFFER)
3. **No keybinding:** User types `ifg` directly
4. **No bind -x:** Works without readline magic
5. **No interactive checks:** Wrappers work everywhere
6. **Editable:** Command can be edited before execution
7. **Portable:** One wrapper for all shells

### Why Embed Instead of Generate?

1. **Easier maintenance:** Edit .sh file with syntax highlighting
2. **Single binary:** No external files
3. **Testability:** Can test shell script independently
4. **Future-proof:** Easy to extend

---

## Testing

### Unit Test Flags

```bash
# Test --sh flag
go build -o ifg
./ifg --sh
# Should output unified wrapper code
```

### Integration Test (Interactive Terminal Required)

```bash
# In an interactive terminal
go build -o ifg
export PATH="$PWD:$PATH"

# Load wrapper
source "$(ifg --sh)"

# Verify function is defined
type ifg
# Should show function definition

# Test interactive
ifg
# Select command: "git status"
# Output: Command: git status
#         Press UP to access from history

# Press UP arrow
# Should show: git status
# Edit and execute
```

### Development Testing

```bash
# Build
go build -o ifg

# Add to PATH
export PATH="$PWD:$PATH"

# Test direct invocation (no wrapper)
./ifg
# Prints command to stdout

# Test with wrapper
source "$(ifg --sh)"
ifg
# Select command
# Command added to history + message printed
```

---

## README Updates

```markdown
## Shell Integration

To use ifg commands directly with history access:

**Bash or Zsh** - Add to `~/.bashrc` or `~/.zshrc`:
```bash
source "$(ifg --sh)"
```

Now when you run `ifg`:
1. Select a command
2. Command is added to shell history
3. Message: "Command: <cmd>\nPress UP to access"
4. Press UP to get the command
5. Edit if needed, press Enter to execute

**Without integration:** `ifg` prints to stdout (useful for scripts).

**For scripts:**
```bash
cmd=$(ifg)
# Select command
echo "$cmd"
```
```

---

## Implementation Steps

1. **Delete** `shell/ifg.bash` and `shell/ifg.zsh`
2. **Create** `shell/ifg.sh` with unified wrapper
3. **Update** `main.go`:
   - Remove `bashWrapper` and `zshWrapper` embeds
   - Add single `shellWrapper` embed for `shell/ifg.sh`
   - Replace `--bash` and `--zsh` with `--sh`
   - Update help text
4. **Update** `README.md` with new integration instructions
5. **Build and test**

**Time estimate:** 15-20 minutes

---

## Questions Answered

### Why not separate bash/zsh wrappers?

Both bash and zsh support `history -s` identically. No need for separate files.

### What about fish?

Fish also supports history. To support fish, we'd add:
```fish
# shell/ifg.fish
function ifg
    set cmd (command ifg)
    if test -n "$cmd"
        history add "$cmd"
        echo "Command: $cmd"
        echo "Press UP to access from history"
    end
end
```

But fish integration is for future - not MVP.

### Should the binary print the message too?

**No** - keep binary output clean. The wrapper prints the message. This keeps `ifg` usable in scripts:
```bash
cmd=$(ifg)  # Clean output, no extra messages
```

### What if user doesn't want history?

They don't load the wrapper:
```bash
# Without wrapper - clean stdout output
ifg
# Outputs: <cmd>

# Pipe to clipboard
ifg | xclip -selection clipboard
```

---

## Success Criteria

1. ✅ Single `--sh` flag (no separate bash/zsh)
2. ✅ Works in bash and zsh with same wrapper
3. ✅ Command added to history after selection
4. ✅ Helpful message printed to user
5. ✅ Easy to install and use
6. ✅ Clean separation: binary outputs, wrapper handles history
7. ✅ Documentation updated
8. ✅ Simpler than READLINE_LINE approach

---

## Timeline

- **Phase 1:** Replace wrapper files (5 min)
- **Phase 2:** Update main.go (10 min)
- **Phase 3:** Update README (5 min)
- **Phase 4:** Build and test (5 min)

**Total:** 25 minutes