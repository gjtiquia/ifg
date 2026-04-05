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

func findRowContaining(screen *MockScreen, text string) int {
	for y := 0; y < screen.height; y++ {
		if strings.Contains(screen.RowAt(y), text) {
			return y
		}
	}
	return -1
}
