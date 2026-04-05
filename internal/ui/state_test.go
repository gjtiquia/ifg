package ui

import (
	"testing"

	"github.com/gjtiquia/ifg/internal/config"
)

func TestNavigateUp(t *testing.T) {
	entries := makeEntries(20)

	t.Run("decrements SelectedIdx when not at first entry", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 5

		state.NavigateUp()

		if state.SelectedIdx != 4 {
			t.Errorf("expected SelectedIdx 4, got %d", state.SelectedIdx)
		}
	})

	t.Run("does nothing when already at first entry", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 0

		state.NavigateUp()

		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
	})

	t.Run("updates ScrollOffset when selection moves above visible area", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.ScrollOffset = 3
		state.SelectedIdx = 3

		state.NavigateUp()

		if state.ScrollOffset != 2 {
			t.Errorf("expected ScrollOffset 2, got %d", state.ScrollOffset)
		}
		if state.SelectedIdx != 2 {
			t.Errorf("expected SelectedIdx 2, got %d", state.SelectedIdx)
		}
	})

	t.Run("multiple NavigateUp calls work correctly", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 10
		state.ScrollOffset = 5

		for i := 0; i < 5; i++ {
			state.NavigateUp()
		}

		if state.SelectedIdx != 5 {
			t.Errorf("expected SelectedIdx 5, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset != 5 {
			t.Errorf("expected ScrollOffset 5 (selected item still visible), got %d", state.ScrollOffset)
		}
	})
}

func TestNavigateDown(t *testing.T) {
	entries := makeEntries(20)

	t.Run("increments SelectedIdx when not at last entry", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 5

		state.NavigateDown()

		if state.SelectedIdx != 6 {
			t.Errorf("expected SelectedIdx 6, got %d", state.SelectedIdx)
		}
	})

	t.Run("does nothing when already at last entry", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 19

		state.NavigateDown()

		if state.SelectedIdx != 19 {
			t.Errorf("expected SelectedIdx 19, got %d", state.SelectedIdx)
		}
	})

	t.Run("updates ScrollOffset when selection moves below visible area", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 5
		state.ScrollOffset = 0

		for i := 0; i < 10; i++ {
			state.NavigateDown()
		}

		if state.ScrollOffset < 1 {
			t.Errorf("expected ScrollOffset >= 1, got %d", state.ScrollOffset)
		}
	})

	t.Run("small terminal height handles scrolling correctly", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 12

		state.SelectedIdx = 0
		state.ScrollOffset = 0

		for i := 0; i < 15; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 15 {
			t.Errorf("expected SelectedIdx 15, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset < 1 {
			t.Errorf("expected ScrollOffset >= 1, got %d", state.ScrollOffset)
		}
	})
}

func TestScrollBoundaryConditions(t *testing.T) {
	t.Run("empty filtered list - navigation does nothing", func(t *testing.T) {
		state := NewState([]config.Entry{})
		state.TerminalHeight = 24

		state.NavigateUp()
		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}

		state.NavigateDown()
		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
	})

	t.Run("single entry - navigation stays within bounds", func(t *testing.T) {
		entries := makeEntries(1)
		state := NewState(entries)
		state.TerminalHeight = 24

		state.NavigateUp()
		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}

		state.NavigateDown()
		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
	})

	t.Run("very small terminal height", func(t *testing.T) {
		entries := makeEntries(50)
		state := NewState(entries)
		state.TerminalHeight = 10

		for i := 0; i < 30; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 30 {
			t.Errorf("expected SelectedIdx 30, got %d", state.SelectedIdx)
		}

		for i := 0; i < 30; i++ {
			state.NavigateUp()
		}

		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset != 0 {
			t.Errorf("expected ScrollOffset 0, got %d", state.ScrollOffset)
		}
	})

	t.Run("large number of entries - scroll advances properly", func(t *testing.T) {
		entries := makeEntries(100)
		state := NewState(entries)
		state.TerminalHeight = 24

		for i := 0; i < 80; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 80 {
			t.Errorf("expected SelectedIdx 80, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset < 70 {
			t.Errorf("expected ScrollOffset >= 70, got %d", state.ScrollOffset)
		}

		for i := 0; i < 80; i++ {
			state.NavigateUp()
		}

		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset != 0 {
			t.Errorf("expected ScrollOffset 0, got %d", state.ScrollOffset)
		}
	})
}

func TestNavigationAfterSearch(t *testing.T) {
	entries := []config.Entry{
		{Title: "git status", Command: "git status"},
		{Title: "git commit", Command: "git commit"},
		{Title: "docker ps", Command: "docker ps"},
		{Title: "docker build", Command: "docker build"},
	}

	t.Run("search filters entries and resets scroll", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SelectedIdx = 3
		state.ScrollOffset = 2

		state.SearchBuf = "docker"
		state.UpdateSearch()

		if len(state.Filtered) != 2 {
			t.Errorf("expected 2 filtered entries, got %d", len(state.Filtered))
		}
		if state.SelectedIdx >= len(state.Filtered) {
			t.Errorf("SelectedIdx %d out of bounds for %d filtered entries", state.SelectedIdx, len(state.Filtered))
		}
		if state.ScrollOffset != 0 {
			t.Errorf("expected ScrollOffset 0, got %d", state.ScrollOffset)
		}
	})

	t.Run("navigation works correctly after search", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24
		state.SearchBuf = "docker"
		state.UpdateSearch()

		state.NavigateDown()
		if state.SelectedIdx != 1 {
			t.Errorf("expected SelectedIdx 1, got %d", state.SelectedIdx)
		}

		state.NavigateDown()
		if state.SelectedIdx != 1 {
			t.Errorf("expected SelectedIdx 1 (at last), got %d", state.SelectedIdx)
		}

		state.NavigateUp()
		if state.SelectedIdx != 0 {
			t.Errorf("expected SelectedIdx 0, got %d", state.SelectedIdx)
		}
	})
}

func makeEntries(count int) []config.Entry {
	entries := make([]config.Entry, count)
	for i := 0; i < count; i++ {
		entries[i] = config.Entry{
			Title:   "command",
			Command: "cmd",
		}
	}
	return entries
}
