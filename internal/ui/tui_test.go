package ui

import (
	"strings"
	"testing"

	"github.com/gjtiquia/ifg/internal/config"
)

func TestRenderEntriesWithinBounds(t *testing.T) {
	entries := makeEntries(20)
	state := NewState(entries)
	state.TerminalHeight = 20

	screen := NewMockScreen(80, 20)
	Render(state, screen)

	maxRow := screen.MaxRow()
	expectedMaxRow := state.TerminalHeight - bottomPadding - 1

	if maxRow > expectedMaxRow {
		t.Errorf("content rendered past max row: got row %d, expected max %d", maxRow, expectedMaxRow)
	}
}

func TestRenderSelectionIndicator(t *testing.T) {
	entries := []config.Entry{
		{Title: "first", Command: "cmd1"},
		{Title: "second", Command: "cmd2"},
		{Title: "third", Command: "cmd3"},
	}
	state := NewState(entries)
	state.TerminalHeight = 24
	state.SelectedIdx = 1

	screen := NewMockScreen(80, 24)
	Render(state, screen)

	row := findRowContaining(screen, "> # second")
	if row == -1 {
		t.Errorf("selected entry should have '> ' prefix")
	}

	row = findRowContaining(screen, "# first")
	if row != -1 && strings.HasPrefix(screen.RowAt(row), "> ") {
		t.Errorf("non-selected entry should not have '> ' prefix")
	}
}

func TestRenderScrollIndicatorPosition(t *testing.T) {
	entries := makeEntries(50)
	state := NewState(entries)
	state.TerminalHeight = 20
	state.SelectedIdx = 30

	screen := NewMockScreen(80, 20)
	Render(state, screen)

	maxRow := state.TerminalHeight - bottomPadding

	scrollIndicatorRow := findRowContaining(screen, "[")
	if scrollIndicatorRow == -1 {
		t.Errorf("scroll indicator should be visible")
	}

	if scrollIndicatorRow >= maxRow {
		t.Errorf("scroll indicator at row %d should be below maxRow %d", scrollIndicatorRow, maxRow)
	}
}

func TestRenderBottomPaddingEmpty(t *testing.T) {
	entries := makeEntries(50)
	state := NewState(entries)
	state.TerminalHeight = 20

	screen := NewMockScreen(80, 20)
	Render(state, screen)

	maxRow := state.TerminalHeight - bottomPadding
	expectedMaxRow := maxRow - 1

	actualMaxRow := screen.MaxRow()

	if actualMaxRow > expectedMaxRow {
		t.Errorf("content rendered into bottom padding: max row %d, expected <= %d", actualMaxRow, expectedMaxRow)
	}

	for row := actualMaxRow + 1; row < state.TerminalHeight; row++ {
		rowContent := strings.TrimSpace(screen.RowAt(row))
		if rowContent != "" {
			t.Errorf("row %d should be empty for bottom padding, got: %q", row, rowContent)
		}
	}
}

func TestRenderClipsMidEntry(t *testing.T) {
	entries := []config.Entry{
		{
			Title:       "entry with description",
			Description: []string{"line1", "line2", "line3", "line4", "line5"},
			Command:     "cmd",
		},
	}
	state := NewState(entries)
	state.TerminalHeight = headerRows + bottomPadding + 3

	screen := NewMockScreen(80, state.TerminalHeight)
	Render(state, screen)

	maxRow := state.TerminalHeight - bottomPadding
	actualMaxRow := screen.MaxRow()

	if actualMaxRow >= maxRow {
		t.Errorf("rendered past maxRow %d, got max row %d", maxRow, actualMaxRow)
	}
}

func TestRenderEmptyList(t *testing.T) {
	state := NewState([]config.Entry{})
	state.TerminalHeight = 24

	screen := NewMockScreen(80, 24)
	Render(state, screen)

	row := findRowContaining(screen, "No results found")
	if row == -1 {
		t.Errorf("should display 'No results found' for empty list")
	}
}

func TestRenderHeaderPresent(t *testing.T) {
	entries := makeEntries(5)
	state := NewState(entries)
	state.TerminalHeight = 24

	screen := NewMockScreen(80, 24)
	Render(state, screen)

	row := findRowContaining(screen, "ifg - [i] [f]or[g]ot")
	if row != 0 {
		t.Errorf("header should be at row 0, got row %d", row)
	}
}

func TestRenderSearchPrompt(t *testing.T) {
	entries := makeEntries(5)
	state := NewState(entries)
	state.TerminalHeight = 24

	t.Run("insert mode shows type to search", func(t *testing.T) {
		state.Mode = ModeInsert
		screen := NewMockScreen(80, 24)
		Render(state, screen)

		row := findRowContaining(screen, "type to search:")
		if row == -1 {
			t.Errorf("should display 'type to search:' prompt")
		}
	})

	t.Run("normal mode shows search results for", func(t *testing.T) {
		state.Mode = ModeNormal
		screen := NewMockScreen(80, 24)
		Render(state, screen)

		row := findRowContaining(screen, "search results for:")
		if row == -1 {
			t.Errorf("should display 'search results for:' prompt")
		}
	})
}

func TestRenderScrollIndicatorDoesNotOverlapEntries(t *testing.T) {
	entries := makeEntries(50)
	state := NewState(entries)
	state.TerminalHeight = 24
	state.SelectedIdx = 0

	screen := NewMockScreen(80, state.TerminalHeight)
	Render(state, screen)

	maxRow := state.TerminalHeight - bottomPadding

	scrollIndicatorRow := -1
	for y := 0; y < state.TerminalHeight; y++ {
		rowContent := screen.RowAt(y)
		if strings.Contains(rowContent, "[") && strings.Contains(rowContent, " of ") {
			scrollIndicatorRow = y
			break
		}
	}

	if scrollIndicatorRow == -1 {
		t.Fatalf("scroll indicator should be visible")
	}

	if scrollIndicatorRow >= maxRow {
		t.Errorf("scroll indicator at row %d should be below maxRow %d", scrollIndicatorRow, maxRow)
	}

	scrollRowContent := strings.TrimSpace(screen.RowAt(scrollIndicatorRow))
	hasEntryMarker := strings.Contains(scrollRowContent, "> ") || (strings.Contains(scrollRowContent, "#") && !strings.HasPrefix(scrollRowContent, "["))
	if hasEntryMarker {
		t.Errorf("scroll indicator overlaps with entry content at row %d: %q", scrollIndicatorRow, scrollRowContent)
	}

	entryMaxRow := -1
	for y := 0; y < state.TerminalHeight; y++ {
		rowContent := screen.RowAt(y)
		if strings.Contains(rowContent, "> ") || strings.Contains(rowContent, "#") {
			if y > entryMaxRow {
				entryMaxRow = y
			}
		}
	}

	if entryMaxRow >= scrollIndicatorRow && scrollIndicatorRow != -1 {
		t.Errorf("entries end at row %d, should be before scroll indicator at row %d", entryMaxRow, scrollIndicatorRow)
	}
}

func TestVariableHeightEntriesStayWithinBounds(t *testing.T) {
	// Create entries with varying heights
	entries := []config.Entry{
		{Title: "short", Command: "cmd"},
		{Title: "tall", Command: "cmd", Description: []string{"line1", "line2", "line3", "line4"}},
		{Title: "short", Command: "cmd"},
		{Title: "medium", Command: "cmd", Description: []string{"line1", "line2"}},
		{Title: "short", Command: "cmd"},
		{Title: "tall", Command: "cmd", Description: []string{"line1", "line2", "line3"}},
		{Title: "short", Command: "cmd"},
		{Title: "short", Command: "cmd"},
		{Title: "medium", Command: "cmd", Description: []string{"line1"}},
		{Title: "short", Command: "cmd"},
		{Title: "tall", Command: "cmd", Description: []string{"line1", "line2", "line3", "line5"}},
		{Title: "short", Command: "cmd"},
		{Title: "short", Command: "cmd"},
		{Title: "medium", Command: "cmd", Description: []string{"line1", "line2"}},
		{Title: "short", Command: "cmd"},
	}

	t.Run("selected entry visible after navigating through variable height entries", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		t.Logf("Entry heights: %v", state.EntryHeights)

		// Navigate through multiple variable-height entries
		for i := 0; i < 8; i++ {
			state.NavigateDown()
			t.Logf("After NavigateDown %d: SelectedIdx=%d, ScrollOffset=%d", i+1, state.SelectedIdx, state.ScrollOffset)
		}

		// Render and check if selected entry is visible
		screen := NewMockScreen(80, state.TerminalHeight)
		Render(state, screen)

		// Debug: print all rows
		for y := 0; y < state.TerminalHeight; y++ {
			rowContent := strings.TrimSpace(screen.RowAt(y))
			if rowContent != "" {
				t.Logf("Row %d: %q", y, rowContent)
			}
		}

		// The selected entry should be visible in the rendered output
		selectedEntry := entries[state.SelectedIdx]
		selectedMarker := "> # " + selectedEntry.Title

		found := false
		for y := 0; y < state.TerminalHeight; y++ {
			if strings.Contains(screen.RowAt(y), selectedMarker) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Selected entry %d (title: %q) not visible in render after navigating",
				state.SelectedIdx, selectedEntry.Title)
		}
	})

	t.Run("entries do not overflow past content area with variable heights", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 20 // Smaller terminal to force scrolling

		// Navigate to middle of list
		for i := 0; i < 7; i++ {
			state.NavigateDown()
		}

		screen := NewMockScreen(80, state.TerminalHeight)
		Render(state, screen)

		maxRow := state.TerminalHeight - bottomPadding

		// Check no entries rendered past content end
		for y := maxRow; y < state.TerminalHeight; y++ {
			rowContent := screen.RowAt(y)
			if strings.Contains(rowContent, "#") || strings.Contains(rowContent, "> ") {
				t.Errorf("Entry content found in bottom padding area at row %d: %q", y, rowContent)
			}
		}
	})

	t.Run("scroll indicator not overlapping with entries after extensive navigation", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 22

		// Navigate through many variable-height entries
		for i := 0; i < 12; i++ {
			state.NavigateDown()
		}

		screen := NewMockScreen(80, state.TerminalHeight)
		Render(state, screen)

		scrollIndicatorRow := -1
		for y := 0; y < state.TerminalHeight; y++ {
			rowContent := screen.RowAt(y)
			if strings.Contains(rowContent, "[") && strings.Contains(rowContent, " of ") {
				scrollIndicatorRow = y
				break
			}
		}

		if scrollIndicatorRow == -1 {
			t.Fatalf("scroll indicator should be visible after navigation")
		}

		scrollRowContent := strings.TrimSpace(screen.RowAt(scrollIndicatorRow))
		hasEntryMarker := strings.Contains(scrollRowContent, "> ") ||
			(strings.Contains(scrollRowContent, "#") && !strings.HasPrefix(scrollRowContent, "["))

		if hasEntryMarker {
			t.Errorf("After navigation, scroll indicator overlaps with entry at row %d: %q",
				scrollIndicatorRow, scrollRowContent)
		}
	})
}

func findRowContaining(screen *MockScreen, text string) int {
	for y := 0; y < screen.height; y++ {
		if strings.Contains(screen.RowAt(y), text) {
			return y
		}
	}
	return -1
}
