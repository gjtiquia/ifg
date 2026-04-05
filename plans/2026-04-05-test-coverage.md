# Test Coverage Improvements Plan

## Overview

Improve test coverage from current levels:
- `main` package: 0% (skipped - requires terminal/integration tests)
- `internal/config`: 85%
- `internal/search`: 100%
- `internal/ui`: 52.4% → **81.2%** ✅

## Target Coverage

- All packages: 80%+
- Critical paths: 100%

---

## Final Coverage Report

| Package | Before | After | Status |
|---------|--------|-------|--------|
| `internal/ui` | 52.4% | 81.2% | ✅ Complete |
| `internal/search` | 100% | 100% | ✅ Already complete |
| `internal/config` | 85% | 85% | ⏸️ Not priority |
| `main` | 0% | 0% | ⏭️ Skipped (requires terminal) |

---

## Completed Phases

### Phase 1: Cursor/Input Operations ✅

| Function | Coverage |
|----------|----------|
| `AppendChar` | 100% |
| `DeleteChar` | 100% |
| `MoveCursorLeft` | 100% |
| `MoveCursorRight` | 100% |

**Tests added in `internal/ui/state_test.go`:**
- `TestAppendChar` - 5 subtests
- `TestDeleteChar` - 5 subtests
- `TestMoveCursorLeft` - 4 subtests
- `TestMoveCursorRight` - 5subtests

---

### Phase 2: Vim Motions ✅

| Function | Coverage |
|----------|----------|
| `MoveWordForward` | 100% |
| `MoveWORDForward` | 100% |
| `MoveWordBackward` | 100% |
| `MoveWORDBackward` | 100% |
| `MoveWordEnd` | 100% |
| `MoveWORDEnd` | 100% |
| `isWordChar` | 100% |
| `isSpace` | 100% |

**Tests added in `internal/ui/state_test.go`:**
- `TestMoveWordForward` - 7 subtests
- `TestMoveWORDForward` - 5 subtests
- `TestMoveWordBackward` - 8 subtests
- `TestMoveWORDBackward` -6 subtests
- `TestMoveWordEnd` - 8 subtests
- `TestMoveWORDEnd` - 6 subtests
- `TestIsWordChar` - 13 subtests
- `TestIsSpace` - 9 subtests

---

### Phase 3: Mode Switching ✅

| Function | Coverage |
|----------|----------|
| `SwitchToNormal` | 100% |
| `SwitchToInsert` | 100% |

**Tests added in `internal/ui/state_test.go`:**
- `TestSwitchToNormal` - 3 subtests
- `TestSwitchToInsert` - 10 subtests

---

### Phase 4: State Getters ✅

| Function | Coverage |
|----------|----------|
| `GetSelectedCommand` | 100% |

**Tests added in `internal/ui/state_test.go`:**
- `TestGetSelectedCommand` - 7 subtests

---

### Phase 5: Style Conversion ✅

| Function | Coverage |
|----------|----------|
| `NewStyle` | 100% |
| `ToTcellStyle` | 100% |

**Tests added in `internal/ui/screen_test.go`:**
- `TestNewStyle` - 1 test
- `TestToTcellStyle` - 8 subtests

---

## Skipped Phases

### Phase 6: Mock Helpers ⏭️

**Reason:** Low priority - these are test helpers, not production code.

| Function | Status |
|----------|--------|
| `Show` | Skipped |
| `ContentAt` | Skipped |
| `HasContentAt` | Skipped |

---

### Phase 7: Input Handling ⏭️

**Reason:** Requires terminal mocking or integration tests.

| Function | Status |
|----------|--------|
| `ReadKey` | Skipped |

---

### Phase 8: Screen/TUI ⏭️

**Reason:** Requires actual terminal for setup/cleanup operations.

| Function | Status |
|----------|--------|
| `SetupTerminal` | Skipped |
| `Restore` | Skipped |
| `GetSize` | Skipped |
| `Screen` | Skipped |
| `WrappedScreen` | Skipped |

---

### Phase 9: Main Package ⏭️

**Reason:** Main function controls process lifecycle with `os.Exit()`. Requires integration tests or process-level mocking.

---

## Files Created/Modified

| File | Action | Status |
|------|--------|--------|
| `internal/ui/state_test.go` | Added tests | ✅ |
| `internal/ui/screen_test.go` | Created new | ✅ |

---

## Summary

**Coverage improvement:** 52.4% → 81.2% (+28.8%)

All critical path functions (state management, cursor operations, vim motions, mode switching) now have 100% test coverage. Terminal-dependent code is appropriately tested via integration or skipped.

---

## Remaining Uncovered Functions

These require terminal/integration tests and are intentionally skipped:

**`internal/ui/input.go`:**
- `ReadKey` - requires mock tcell.Screen with event simulation

**`internal/ui/tui.go`:**
- `SetupTerminal` - requires actual TTY
- `Restore` - requires actual TTY
- `GetSize` - requires actual TTY
- `Screen` - requires actual TTY
- `WrappedScreen` - requires actual TTY

**`internal/ui/tcell.go`:**
- All `TcellScreen` methods - require actual screen

**`internal/ui/mock.go`:**
- `Show`, `ContentAt`, `HasContentAt` - test helpers, low priority

**`main` package:**
- `main`, `run`, `runInputLoop` - require process-level mocking