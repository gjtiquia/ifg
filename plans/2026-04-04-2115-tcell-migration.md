# Migration to tcell - Plan

## Problem

Current implementation has "staggered" output - likely issues with:
- Manual ANSI escape sequence handling
- Timing issues with terminal updates
- Buffer flushing problems
- Rendering artifacts from incremental updates
- Complex escape sequence logic prone to errors

## Solution: Use tcell Library

`tcell` is a battle-tested terminal UI library that handles:
- Terminal capabilities detection (different terminal types)
- Proper escape sequence handling
- Screen buffering and double buffering
- Event handling (keyboard, mouse, resize)
- Color and style support
- Terminal restoration

### Why tcell over stdlib?

**Current approach (stdlib + manual escapes):**
- ✅ No dependencies
- ❌ Manual handling of escape sequences
- ❌ No screen buffering
- ❌ Terminal-specific quirks not handled
- ❌ Easy to introduce rendering bugs
- ❌ Timing/flickering issues

**tcell approach:**
- ✅ Handles all terminal types (xterm, vt100, screen, tmux, Windows)
- ✅ Double buffering (prevent flickering)
- ✅ Proper cleanup on exit
- ✅ Event-driven architecture
- ✅ Resizable screen support
- ✅ ~500KB dependency (acceptable for binary size)
- ❌ Adds external dependency

**Decision:** The manual approach is fragile. tcell is the standard solution for TUIs in Go (used by many popular tools).

---

## Architecture Changes

### Current Structure
```
internal/ui/
├── tui.go      # Manual escape sequences, screen clearing
├── input.go    # Manual key reading, escape sequence parsing
└── state.go    # State management (unchanged)
```

### New Structure with tcell
```
internal/ui/
├── tui.go      # tcell Screen initialization, rendering
├── input.go    # tcellEvents (Key, Resize, etc.)
└── state.go    # State management (unchanged)
```

---

## Implementation Plan

### Phase 1: Add tcell Dependency

```bash
go get github.com/gdamore/tcell/v2
```

### Phase 2: Refactor UI Layer

**File: `internal/ui/tui.go`**

**Remove:**
- `ClearScreen()`
- `MoveCursor()`
- `HideCursor()` / `ShowCursor()`
- `EnterAlternateScreen()` / `ExitAlternateScreen()`
- Manual escape sequences

**Add:**
```go
import "github.com/gdamore/tcell/v2"

type Terminal struct {
    screen tcell.Screen
}

func SetupTerminal() (*Terminal, error) {
    screen, err := tcell.NewScreen()
    if err != nil {
        return nil, err
    }
    if err := screen.Init(); err != nil {
        return nil, err
    }
    screen.SetStyle(tcell.StyleDefault)
    return &Terminal{screen: screen}, nil
}

func (t *Terminal) Restore() {
    t.screen.Fini()
}

func (t *Terminal) GetSize() (int, int) {
    return t.screen.Size()
}

func Render(state *State, screen tcell.Screen) {
    screen.Clear()
    
    // Header
    drawText(screen, 0, 0, "ifg - [i] [f]or[g]ot that cmd again", tcell.StyleDefault.Bold(true))
    
    // Search prompt
    row := 2
    prompt := "type to search: "
    if state.Mode == ModeNormal {
        prompt = "normal mode: "
    }
    drawText(screen, 0, row, prompt + state.SearchBuf, tcell.StyleDefault)
    
    // Separator
    row += 2
    drawText(screen, 0, row, "---", tcell.StyleDefault)
    
    // Entries
    row += 2
    _, height := screen.Size()
    visibleHeight := height - row - 1
    
    for i := 0; i < visibleHeight && i + state.ScrollOffset < len(state.Filtered); i++ {
        entryIdx := i + state.ScrollOffset
        entry := state.Filtered[entryIdx]
        
        isSelected := entryIdx == state.SelectedIdx
        style := tcell.StyleDefault
        if isSelected {
            style = style.Bold(true)
        }
        
        // Title
        if entry.Title != "" {
            prefix := "  "
            if isSelected {
                prefix = "> "
            }
            drawText(screen, 0, row, prefix + "# " + entry.Title, style)
            row++
        }
        
        // Description
        for _, desc := range entry.Description {
            drawText(screen, 0, row, "  # " + desc, style)
            row++
        }
        
        // Command
        drawText(screen, 0, row, "  " + entry.Command, style)
        row += 2
    }
    
    screen.Show()
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
    for i, ch := range text {
        screen.SetContent(x + i, y, ch, nil, style)
    }
}
```

**File: `internal/ui/input.go`**

**Remove:**
- `ReadKey()` manual implementation
- Escape sequence parsing
- Key type enum (use tcell's Key type)

**Add:**
```go
func ReadKey(screen tcell.Screen) (KeyEvent, error) {
    ev := screen.PollEvent()
    
    switch ev := ev.(type) {
    case *tcell.EventKey:
        switch ev.Key() {
        case tcell.KeyUp:
            return KeyEvent{Type: KeyUp}, nil
        case tcell.KeyDown:
            return KeyEvent{Type: KeyDown}, nil
        case tcell.KeyEnter:
            return KeyEvent{Type: KeyEnter}, nil
        case tcell.KeyEscape:
            return KeyEvent{Type: KeyEscape}, nil
        case tcell.KeyBackspace, tcell.KeyBackspace2:
            return KeyEvent{Type: KeyBackspace}, nil
        case tcell.KeyCtrlC:
            return KeyEvent{Type: KeyCtrlC}, nil
        default:
            if ev.Rune() != 0 {
                return KeyEvent{Type: KeyChar, Char: ev.Rune()}, nil
            }
        }
    case *tcell.EventResize:
        // Handled in main loop
    }
    
    return KeyEvent{Type: KeyUnknown}, nil
}
```

**File: `main.go`**

**Changes:**
- Use tcell screen instead of manual terminal setup
- No need for alternate screen handling (tcell does this)
- Simplified signal handling (tcell handles resize events)

```go
func main() {
    // ... config loading ...
    
    term, err := ui.SetupTerminal()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error setting up terminal: %v\n", err)
        os.Exit(2)
    }
    defer term.Restore()
    
    width, height := term.GetSize()
    state := ui.NewState(entries)
    state.TerminalWidth = width
    state.TerminalHeight = height
    
    // No manual signal handling needed - tcell handles resize
    
    selectedCommand := runInputLoop(state, term)
    
    term.Restore()  // Exit tcell screen FIRST
    
    if selectedCommand != "" {
        fmt.Println(selectedCommand)
        os.Exit(0)
    }
    
    os.Exit(1)
}

func runInputLoop(state *ui.State, term *ui.Terminal) string {
    ui.Render(state, term.screen)
    
    for {
        key, err := ui.ReadKey(term.screen)
        if err != nil {
            continue
        }
        
        // ... same input handling ...
        
        ui.Render(state, term.screen)
    }
}
```

### Phase 3: Testing

**Test strategy:**
1. Verify tcell properly initializes and cleans up
2. Test rendering matches expected format
3. Verify keyboard input works
4. Test on different terminals (gnome-terminal, tmux, screen)
5. Verify command output appears in main terminal after exit

---

## Migration Complexity

**Low complexity:**
- State management (`state.go`) - **NO CHANGES NEEDED**
- Config parsing - **NO CHANGES NEEDED**
- Search logic - **NO CHANGES NEEDED**

**Medium complexity:**
- `tui.go` - Replace escape sequences with tcell Screen API
- `input.go` - Replace manual key reading with tcell events
- `main.go` - Simplify terminal setup

**Estimated time:** 2-3 hours

---

## Binary Size Impact

**Current:** 2.6MB
**With tcell:** ~3.1MB (+500KB)

**Analysis:** 500KB increase is acceptable for the stability and proper terminal handling it provides.

---

## Alternative: Keep Current Approach (Not Recommended)

If we keep the current manual approach, we need to:
1. Add synchronization to prevent staggered output
2. Properly flush buffers after each render
3. Handle terminal-specific quirks
4. Fix alternate screen buffer issues
5. Debug rendering artifacts

This would likely take **longer** than switching to tcell and be **more fragile**.

---

## Recommendation

**Switch to tcell.** It's the industry-standard approach for Go TUIs and will:
- Fix all rendering issues automatically
- Provide better terminal compatibility
- Simplify code (less manual escape handling)
- Be more maintainable going forward

The 500KB binary size increase is a worthwhile trade-off for stability.

---

## Questions for User

1. **Proceed with tcell migration?** (Recommended)
2. **Or investigate manual escape sequence fix first?** (Not recommended - likely more complex than tcell migration)

If proceeding with tcell, the implementation is straightforward and should fix all rendering issues.