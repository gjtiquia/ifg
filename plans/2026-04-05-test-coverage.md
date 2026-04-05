# Test Coverage Improvements Plan

## Overview

Improve test coverage from current levels:
- `main` package: 0%
- `internal/config`: 85%
- `internal/search`: 100%
- `internal/ui`: 52.4%

## Target Coverage

- All packages: 80%+
- Critical paths: 100%

---

## Package: `main`

### Current Coverage: 0%

### Files to Create

| File | Purpose |
|------|---------|
| `main_test.go` | Test main function behavior |

### Test Cases

1. **TestPrintHelp**
   - Verify help output contains expected sections
   - Check executable name and usage

2. **TestShellIntegrationFlag**
   - `--sh`, `--bash`, `--zsh` all output shell wrapper
   - Verify output contains function definition

3. **TestUnknownFlag**
   - Unknown flag exits with error
   - Error message shown

4. **TestHelpFlag**
   - `--help` and `-h` show help and exit 0

### Challenge

`main()` function has side effects (terminal setup). Consider:
- Extract logic into testable functions
- Use build tags or flags to skip in CI
- Mock terminal for integration tests

---

## Package: `internal/config`

### Current Coverage: 85%

### Uncovered Functions

| Function | File:Line | Priority |
|----------|-----------|----------|
| `SetDefaultConfig` | config.go:15 | Low |

### Test Cases

1. **TestSetDefaultConfig**
   - Set default config content
   - Verify it's used by `CreateDefaultConfig`

---

## Package: `internal/ui`

### Current Coverage: 52.4%

### Uncovered Functions byCategory

#### Cursor/Input Operations (Priority: High)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `AppendChar` | state.go:72 | Insert character at cursor |
| `DeleteChar` | state.go:78 | Delete character before cursor |
| `MoveCursorLeft` | state.go:173 | Move cursor left |
| `MoveCursorRight` | state.go:179 | Move cursor right |

**Test Cases:**

1. **TestAppendChar**
   - Append to empty buffer
   - Append to middle of buffer (cursor not at end)
   - Append to end of buffer
   - Append unicode characters
   - Verify cursor advances
   - Verify search updates

2. **TestDeleteChar**
   - Delete from empty buffer (no-op)
   - Delete from beginning (no-op)
   - Delete from middle
   - Delete from end
   - Verify cursor decrements
   - Verify search updates

3. **TestMoveCursorLeft**
   - Move left from position 0 (no-op)
   - Move left from middle
   - Move left from end

4. **TestMoveCursorRight**
   - Move right from end (no-op)
   - Move right from middle
   - Move right from beginning

#### Vim Motions (Priority: High)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `MoveWordForward` | state.go:185 | `w` motion |
| `MoveWORDForward` | state.go:200 | `W` motion |
| `MoveWordBackward` | state.go:215 | `b` motion |
| `MoveWORDBackward` | state.go:230 | `B` motion |
| `MoveWordEnd` | state.go:245 | `e` motion |
| `MoveWORDEnd` | state.go:268 | `E` motion |
| `isWordChar` | state.go:291 | Helper |
| `isSpace` | state.go:295 | Helper |

**Test Cases:**

1. **TestMoveWordForward**
   - Empty buffer
   - Single word
   - Multiple words with spaces
   - Word followed by punctuation
   - Punctuation followed by word
   - End of buffer

2. **TestMoveWORDForward**
   - Word with internal punctuation (e.g., "foo_bar")
   - Multiple WORDs separated by spaces
   - Leading/trailing whitespace

3. **TestMoveWordBackward**
   - From end to previous word
   - From middle of word
   - From after punctuation
   - From start (no-op)

4. **TestMoveWORDBackward**
   - Similar to `MoveWordBackward` but WORD semantics

5. **TestMoveWordEnd**
   - From start of word to end
   - From middle of word
   - Skip non-word chars
   - On last word (to end)

6. **TestMoveWORDEnd**
   - Similar to `MoveWordEnd` but WORD semantics

7. **TestIsWordChar**
   - Letters: a-z, A-Z
   - Digits: 0-9
   - Underscore
   - Non-word: punctuation, space, unicode

8. **TestIsSpace**
   - Space, tab, newline, carriage return
   - Non-space characters

#### Mode Switching (Priority: Medium)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `SwitchToNormal` | state.go:143 | Switch to normal mode |
| `SwitchToInsert` | state.go:147 | Switch to insert mode |

**Test Cases:**

1. **TestSwitchToNormal**
   - Mode changes from Insert to Normal
   - Preserves buffer and cursor position

2. **TestSwitchToInsert**
   - `"before"` - cursor stays at position
   - `"after"` - cursor moves right by one
   - `"start"` - cursor moves to beginning
   - `"end"` - cursor moves to end
   - Invalid position defaults to current

#### State Getters (Priority: Medium)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `GetSelectedCommand` | state.go:166 | Get selected command |

**Test Cases:**

1. **TestGetSelectedCommand**
   - Empty filtered list returns ""
   - Valid selection returns command
   - Out of bounds returns ""

#### Input Handling (Priority: Low - requires mocking)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `ReadKey` | input.go:25 | Read keyboard input |

**Test Cases:**

1. **TestReadKey**
   - Character key
   - Up/Down arrows
   - Enter
   - Escape
   - Backspace
   - Ctrl-C
   - Resize event

#### Screen/TUI (Priority: Low - requires terminal)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `SetupTerminal` | tui.go:11 | Initialize terminal |
| `Restore` | tui.go:23 | Cleanup terminal |
| `GetSize` | tui.go:27 | Get terminal dimensions |
| `Screen` | tui.go:31 | Get screen interface |
| `WrappedScreen` | tui.go:35 | Get wrapped screen |

These require actual terminal and are hard to unit test. Consider:
- Integration tests
- Manual testing
- Accept lower coverage for these

#### Style Conversion (Priority: Low)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `ToTcellStyle` | screen.go:24 | Convert Style to tcell.Style |

**Test Cases:**

1. **TestToTcellStyle**
   - Empty style (no flags)
   - Bold only
   - Dim only
   - Reverse only
   - Combinations (Bold+Dim, etc.)

#### Mock Helpers (Priority: Low - test helpers)

| Function | File:Line | Description |
|----------|-----------|-------------|
| `Show` | mock.go:41 | No-op for mock |
| `ContentAt` | mock.go:59 | Get cell content |
| `HasContentAt` | mock.go:80 | Check if cell has content |

**Test Cases:**

1. **TestMockScreenContentAt**
   - Cell within bounds returns content
   - Cell out of bounds returns empty

2. **TestMockScreenHasContentAt**
   - Set content and verify
   - No content returns false

---

## Implementation Order

1. **Phase 1**: Cursor/Input operations (state_test.go)
2. **Phase 2**: Vim motions (state_test.go)
3. **Phase 3**: Mode switching (state_test.go)
4. **Phase 4**: State getters (state_test.go)
5. **Phase 5**: Style conversion (screen_test.go)
6. **Phase 6**: Mock helpers (tui_test.go additions)
7. **Phase 7**: Input handling (input_test.go - may need mock tcell.Screen)
8. **Phase 8**: Main package (main_test.go)

---

## Files to Create/Modify

| File | Action |
|------|--------|
| `main_test.go` | Create new |
| `internal/ui/state_test.go` | Add tests |
| `internal/ui/screen_test.go` | Create new |
| `internal/ui/input_test.go` | Create new (may be limited) |

---

## Notes

- `tcell.go` functions cannot be unit tested without actual terminal
- Consider integration tests for terminal-related code
- `ReadKey` depends on tcell.Screen.PollEvent() which requires terminal