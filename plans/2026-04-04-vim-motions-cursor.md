# Vim Motions and Visual Cursor - Implementation Plan

## Overview

Add visual cursor support and vim-style motions for navigating within the search buffer in NORMAL mode.

---

## Requirements Summary

### Visual Cursor
- **Insert mode:** Show cursor (`_` or `| `) at current position in search buffer
- **Normal mode:** Keep current `>` prefix (no additional cursor)

### Vim Motions (Normal Mode Only)
All motions operate on the **search buffer**, not on the entry list:

- `h` - move cursor left (character)
- `l` - move cursor right (character)  
- `w` - move to next word start
- `W` - move to next WORD start (whitespace-delimited)
- `b` - move to previous word start
- `B` - move to previous WORD start
- `e` - move to next word end
- `E` - move to next WORD end

**Word vs WORD:**
- **word:** Alphanumeric characters only (`[a-zA-Z0-9_]`)
- **WORD:** Non-whitespace characters (delimited by whitespace)

---

## Current Architecture

### State Management (`internal/ui/state.go`)
```go
type State struct {
    Mode        Mode
    SearchBuf   string     // Current search input
    CursorIdx   int        // Position in SearchBuf
    SelectedIdx int        // Which entry is selected
    Entries     []config.Entry
    Filtered    []config.Entry
    // ...
}
```

**Already has `CursorIdx`** - tracks cursor position in search buffer!

### Input Handling (`internal/ui/input.go`)
- Insert mode: Currently appends all printable chars to `SearchBuf`
- No cursor movement or word navigation

---

## Implementation Details

### 1. Visual Cursor Rendering

**File: `internal/ui/tui.go`**

Current search prompt rendering:
```go
prompt := "type to search: "
if state.Mode == ModeNormal {
    prompt = "normal mode: "
}
fmt.Printf("%s%s\n", prompt, state.SearchBuf)
```

**Updated rendering:**
```go
prompt := "type to search: "
if state.Mode == ModeNormal {
    prompt = "normal mode: "
}

// Draw search buffer with cursor
fmt.Print(prompt)
if state.Mode == ModeInsert && len(state.SearchBuf) > 0 {
    before := state.SearchBuf[:state.CursorIdx]
    at := ""
    if state.CursorIdx < len(state.SearchBuf) {
        at = string(state.SearchBuf[state.CursorIdx])
    }
    after := ""
    if state.CursorIdx < len(state.SearchBuf) {
        after = state.SearchBuf[state.CursorIdx+1:]
    }
    
    // Print with cursor highlight
    fmt.Print(before)
    if at != "" {
        // Draw character under cursor with reverse video
        fmt.Printf("\x1b[7m%s\x1b[0m", at)  // Reverse video
    }
    fmt.Print(after)
} else if state.Mode == ModeInsert && len(state.SearchBuf) == 0 {
    // Empty buffer - show cursor at position 0
    // Just draw prompt, cursor is at start
}
fmt.Println()
```

**ANSI escape codes:**
- `\x1b[7m` - reverse video (or use `\x1b[4m` for underline)
- `\x1b[0m` - reset

**Cursor styles:**
- Block cursor: `\x1b[7m%s\x1b[0m` (reverse video of character)
- Underline: `\x1b[4m%s\x1b[0m` 
- Bar cursor: Print `|` at position

**Recommendation:** Block cursor (reverse video) - matches vim behavior.

---

### 2. Cursor Movement Functions

**File: `internal/ui/state.go`**

Add methods for cursor movement:

```go
// Move cursor left one character
func (s *State) MoveCursorLeft() {
    if s.CursorIdx > 0 {
        s.CursorIdx--
    }
}

// Move cursor right one character  
func (s *State) MoveCursorRight() {
    if s.CursorIdx < len(s.SearchBuf) {
        s.CursorIdx++
    }
}

// Move cursor to next word start
func (s *State) MoveWordForward() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Skip current word (alphanumeric)
    for i < len(buf) && isWordChar(buf[i]) {
        i++
    }
    
    // Skip whitespace/punctuation
    for i < len(buf) && !isWordChar(buf[i]) {
        i++
    }
    
    s.CursorIdx = i
}

// Move cursor to next WORD start (whitespace-delimited)
func (s *State) MoveWORDForward() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Skip non-whitespace
    for i < len(buf) && !isSpace(buf[i]) {
        i++
    }
    
    // Skip whitespace
    for i < len(buf) && isSpace(buf[i]) {
        i++
    }
    
    s.CursorIdx = i
}

// Move cursor to previous word start
func (s *State) MoveWordBackward() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Skip whitespace/punctuation before cursor
    for i > 0 && !isWordChar(buf[i-1]) {
        i--
    }
    
    // Skip word characters
    for i > 0 && isWordChar(buf[i-1]) {
        i--
    }
    
    s.CursorIdx = i
}

// Move cursor to previous WORD start
func (s *State) MoveWORDBackward() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Skip whitespace before cursor
    for i > 0 && isSpace(buf[i-1]) {
        i--
    }
    
    // Skip non-whitespace
    for i > 0 && !isSpace(buf[i-1]) {
        i--
    }
    
    s.CursorIdx = i
}

// Move cursor to next word end
func (s *State) MoveWordEnd() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Move right at least once if not already at end
    if i < len(buf) {
        i++
    }
    
    // Skip whitespace/punctuation
    for i < len(buf) && !isWordChar(buf[i]) {
        i++
    }
    
    // Find end of word
    for i < len(buf) && isWordChar(buf[i]) {
        i++
    }
    
    s.CursorIdx = i - 1
    if s.CursorIdx < 0 {
        s.CursorIdx = 0
    }
}

// Move cursor to next WORD end
func (s *State) MoveWORDEnd() {
    buf := []rune(s.SearchBuf)
    i := s.CursorIdx
    
    // Move right at least once
    if i < len(buf) {
        i++
    }
    
    // Skip whitespace
    for i < len(buf) && isSpace(buf[i]) {
        i++
    }
    
    // Find end of WORD
    for i < len(buf) && !isSpace(buf[i]) {
        i++
    }
    
    s.CursorIdx = i - 1
    if s.CursorIdx < 0 {
        s.CursorIdx = 0
    }
}

// Helper: check if rune is a word character (alphanumeric + underscore)
func isWordChar(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// Helper: check if rune is whitespace
func isSpace(r rune) bool {
    return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
```

---

### 3. Input Handling Updates

**Key Design Decision:**

In INSERT mode, regular characters (`h`, `l`, `w`, `b`, `e`) should INSERT, not navigate.

Two approaches:

**Approach A:** Use Ctrl modifier in INSERT mode
- `Ctrl+h` = cursor left
- `Ctrl+l` = cursor right
- `Ctrl+w` = next word
- etc.
- Regular characters insert

**Approach B:** Use NORMAL mode for motions
- NORMAL mode already exists
- Add `h`/`l`/`w`/`b`/`e` motions in NORMAL mode
- These motions navigate the search buffer cursor
- `j`/`k` navigate entries (existing)

**Recommendation:** Approach B - motions in NORMAL mode.

This is more vim-like:
- INSERT mode: insert text
- NORMAL mode: navigate (both cursor and entries)

---

### 4. Key Binding Summary

**INSERT mode:**
- Printable chars → insert
- Backspace → delete char before cursor
- `↑`/`↓` → navigate entries
- Enter → select
- Esc → NORMAL mode
- **No motions** (all printable chars insert)

**NORMAL mode:**
- `j`/`↓` → next entry
- `k`/`↑` → previous entry
- `h` → cursor left in search buffer
- `l` → cursor right in search buffer
- `w` → next word
- `W` → next WORD
- `b` → previous word
- `B` → previous WORD
- `e` → word end
- `E` → WORD end
- `i`/`I` → INSERT mode (cursor at start)
- `a`/`A` → INSERT mode (cursor at end)
- Enter → select
- Esc → exit

---

### 5. Input Handling Implementation

**File: `internal/ui/input.go`**

Add Ctrl key detection:

```go
type KeyEvent struct {
    Type  Key
    Char  rune
    Ctrl  bool
    Alt   bool
}

func ReadKey(screen tcell.Screen) KeyEvent {
    ev := screen.PollEvent()
    
    switch ev := ev.(type) {
    case *tcell.EventKey:
        switch ev.Key() {
        case tcell.KeyUp:
            return KeyEvent{Type: KeyUp}
        case tcell.KeyDown:
            return KeyEvent{Type: KeyDown}
        case tcell.KeyEnter:
            return KeyEvent{Type: KeyEnter}
        case tcell.KeyEscape:
            return KeyEvent{Type: KeyEscape}
        case tcell.KeyBackspace, tcell.KeyBackspace2:
            return KeyEvent{Type: KeyBackspace}
        case tcell.KeyCtrlC:
            return KeyEvent{Type: KeyCtrlC}
        default:
            if ev.Rune() != 0 {
                return KeyEvent{Type: KeyChar, Char: ev.Rune()}
            }
        }
    case *tcell.EventResize:
        // Handled elsewhere
    }
    
    return KeyEvent{Type: KeyUnknown}
}
```

Note: tcell doesn't directly support Ctrl+letter detection. We'd need to check `ev.Modifiers()` for Ctrl.

Actually, looking at tcell docs, Ctrl key combinations are detected via `ev.Modifiers()`.

But for simplicity, let's use NORMAL mode for all motions instead.

---

### 6. Updated Input Loop

**File: `main.go`**

```go
func runInputLoop(state *ui.State, term *ui.Terminal) string {
    ui.Render(state, term.Screen())
    
    for {
        key := ui.ReadKey(term.Screen())
        
        switch state.Mode {
        case ui.ModeInsert:
            switch key.Type {
            case ui.KeyChar:
                state.AppendChar(key.Char)
            case ui.KeyBackspace:
                state.DeleteChar()
            case ui.KeyUp:
                state.NavigateUp()
            case ui.KeyDown:
                state.NavigateDown()
            case ui.KeyEnter:
                return state.GetSelectedCommand()
            case ui.KeyEscape:
                state.SwitchToNormal()
            case ui.KeyCtrlC:
                return ""
            }
            
        case ui.ModeNormal:
            switch key.Type {
            case ui.KeyChar:
                switch key.Char {
                case 'j':
                    state.NavigateDown()
                case 'k':
                    state.NavigateUp()
                case 'h':
                    state.MoveCursorLeft()
                case 'l':
                    state.MoveCursorRight()
                case 'w':
                    state.MoveWordForward()
                case 'W':
                    state.MoveWORDForward()
                case 'b':
                    state.MoveWordBackward()
                case 'B':
                    state.MoveWORDBackward()
                case 'e':
                    state.MoveWordEnd()
                case 'E':
                    state.MoveWORDEnd()
                case 'i', 'I':
                    state.SwitchToInsert("start")
                case 'a', 'A':
                    state.SwitchToInsert("end")
                }
            case ui.KeyUp:
                state.NavigateUp()
            case ui.KeyDown:
                state.NavigateDown()
            case ui.KeyEnter:
                return state.GetSelectedCommand()
            case ui.KeyEscape:
                return ""
            case ui.KeyCtrlC:
                return ""
            }
        }
        
        ui.Render(state, term.Screen())
    }
}
```

---

## Testing Strategy

1. **Visual cursor:**
   - Type in INSERT mode, verify cursor moves right
   - Switch to NORMAL mode, move cursor with h/l, verify position
   - Verify cursor doesn't go negative or past end of buffer

2. **Word motions:**
   - Type "foo bar baz", test w/b/e
   - Type "foo-bar_baz", test word vs WORD
   - Type "  spaces  ", test whitespace handling

3. **Edge cases:**
   - Empty search buffer (no cursor shown)
   - Cursor at start (h doesn't go negative)
   - Cursor at end (l doesn't go past end)
   - Single character
   - Unicode characters (use `[]rune` throughout)

---

## Implementation Steps

1. Add cursor movement methods to `internal/ui/state.go`
2. Add word character detection helpers
3. Update `main.go` input loop for NORMAL mode motions
4. Update `internal/ui/tui.go` for cursor rendering
5. Test all motions
6. Update README with key bindings

**Estimated time:** 2-3 hours
