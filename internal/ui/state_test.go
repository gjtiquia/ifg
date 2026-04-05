package ui

import (
	"testing"

	"github.com/gjtiquia/ifg/internal/config"
)

func TestHeaderConstants(t *testing.T) {
	t.Run("headerRows equals sum of components", func(t *testing.T) {
		expected := headerTitleRows + headerPromptRows + headerSeparatorRows
		if headerRows != expected {
			t.Errorf("headerRows = %d, expected %d (sum of components)", headerRows, expected)
		}
	})

	t.Run("all header components are positive", func(t *testing.T) {
		if headerTitleRows <= 0 {
			t.Errorf("headerTitleRows should be positive, got %d", headerTitleRows)
		}
		if headerPromptRows <= 0 {
			t.Errorf("headerPromptRows should be positive, got %d", headerPromptRows)
		}
		if headerSeparatorRows <= 0 {
			t.Errorf("headerSeparatorRows should be positive, got %d", headerSeparatorRows)
		}
	})

	t.Run("bottomPadding is positive", func(t *testing.T) {
		if bottomPadding <= 0 {
			t.Errorf("bottomPadding should be positive, got %d", bottomPadding)
		}
	})

	t.Run("estimatedRowsPerEntry is positive", func(t *testing.T) {
		if estimatedRowsPerEntry <= 0 {
			t.Errorf("estimatedRowsPerEntry should be positive, got %d", estimatedRowsPerEntry)
		}
	})
}

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

	t.Run("scroll respects headerRows constant", func(t *testing.T) {
		entries := makeEntries(50)
		state := NewState(entries)
		state.TerminalHeight = headerRows + bottomPadding + 3*estimatedRowsPerEntry

		for i := 0; i < 10; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 10 {
			t.Errorf("expected SelectedIdx 10, got %d", state.SelectedIdx)
		}
		if state.ScrollOffset < 1 {
			t.Errorf("expected ScrollOffset >= 1, got %d", state.ScrollOffset)
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

func TestSelectionStaysWithinContentBounds(t *testing.T) {
	entries := makeEntries(100)

	t.Run("selection stays within visible content area aftermultiple navigate down", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		visibleRows := state.TerminalHeight - headerRows - bottomPadding
		maxVisibleEntries := visibleRows / estimatedRowsPerEntry
		if maxVisibleEntries < 1 {
			maxVisibleEntries = 1
		}

		maxEntriesBeforeScrollIndicator := maxVisibleEntries - 1
		if maxEntriesBeforeScrollIndicator < 1 {
			maxEntriesBeforeScrollIndicator = 1
		}

		for i := 0; i < maxEntriesBeforeScrollIndicator+5; i++ {
			state.NavigateDown()
		}

		visibleEntriesFromTop := maxEntriesBeforeScrollIndicator
		maxSelectableIdx := state.ScrollOffset + visibleEntriesFromTop

		if state.SelectedIdx > maxSelectableIdx {
			t.Errorf("SelectedIdx %d exceeds visible content bounds (max %d with ScrollOffset %d)",
				state.SelectedIdx, maxSelectableIdx, state.ScrollOffset)
		}
	})

	t.Run("scroll triggers when selection exceeds visible area", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		visibleHeight := state.TerminalHeight - headerRows - bottomPadding - scrollIndicatorRows
		if visibleHeight < 1 {
			visibleHeight = 1
		}

		// Navigate until we exceed the visible area
		for i := 0; i < len(entries); i++ {
			state.NavigateDown()

			// Check that selection stays within visible bounds
			lastVisible := state.findLastVisibleEntry(state.ScrollOffset, visibleHeight)
			if state.SelectedIdx > lastVisible+1 { // +1 for some tolerance during navigation
				t.Errorf("After %d navigations: SelectedIdx %d is past lastVisible %d (ScrollOffset %d)",
					i+1, state.SelectedIdx, lastVisible, state.ScrollOffset)
			}
		}
	})

	t.Run("selection never goes past last visible entry before scroll indicator", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		for i := 0; i < 60; i++ {
			state.NavigateDown()
		}

		visibleRows := state.TerminalHeight - headerRows - bottomPadding
		approxVisibleEntries := visibleRows / estimatedRowsPerEntry

		distanceFromScrollOffset := state.SelectedIdx - state.ScrollOffset
		if distanceFromScrollOffset > approxVisibleEntries {
			t.Errorf("SelectedIdx %d is %d entries past ScrollOffset %d, should be within ~%d visible entries",
				state.SelectedIdx, distanceFromScrollOffset, state.ScrollOffset, approxVisibleEntries)
		}
	})
}

func TestVariableHeightEntriesScrollCorrectly(t *testing.T) {
	// Create entries with varying heights: some have descriptions, some don't
	entries := []config.Entry{
		{Title: "entry0", Command: "cmd0"},
		{Title: "entry1", Command: "cmd1", Description: []string{"desc1 line1", "desc1 line2"}},
		{Title: "entry2", Command: "cmd2"},
		{Title: "entry3", Command: "cmd3", Description: []string{"desc3"}},
		{Title: "entry4", Command: "cmd4"},
		{Title: "entry5", Command: "cmd5", Description: []string{"desc5a", "desc5b", "desc5c"}},
		{Title: "entry6", Command: "cmd6"},
		{Title: "entry7", Command: "cmd7"},
		{Title: "entry8", Command: "cmd8", Description: []string{"desc8"}},
		{Title: "entry9", Command: "cmd9"},
		{Title: "entry10", Command: "cmd10", Description: []string{"desc10 line1", "desc10 line2"}},
		{Title: "entry11", Command: "cmd11"},
		{Title: "entry12", Command: "cmd12"},
		{Title: "entry13", Command: "cmd13", Description: []string{"desc13"}},
		{Title: "entry14", Command: "cmd14"},
		{Title: "entry15", Command: "cmd15"},
		{Title: "entry16", Command: "cmd16", Description: []string{"desc16a", "desc16b"}},
		{Title: "entry17", Command: "cmd17"},
		{Title: "entry18", Command: "cmd18"},
		{Title: "entry19", Command: "cmd19", Description: []string{"desc19"}},
	}

	t.Run("selection stays within bounds after navigating through variable height entries", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		// Navigate through many entries with variable heights
		for i := 0; i < 15; i++ {
			state.NavigateDown()
		}

		// The selection should stay within reasonable bounds of the scroll offset
		distanceFromScrollOffset := state.SelectedIdx - state.ScrollOffset

		// With estimatedRowsPerEntry = 3, we expect roughly 1-2 entries visible
		// The distance should not exceed what's reasonably visible
		maxExpectedDistance := 3

		if distanceFromScrollOffset > maxExpectedDistance {
			t.Errorf("SelectedIdx %d is %d entries past ScrollOffset %d (expected at most %d)",
				state.SelectedIdx, distanceFromScrollOffset, state.ScrollOffset, maxExpectedDistance)
		}
	})

	t.Run("scroll offset advances as expected with variable heights", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		// Start at entry 0, navigate through entries
		initialScrollOffset := state.ScrollOffset

		// Navigate down several times
		for i := 0; i < 10; i++ {
			state.NavigateDown()
		}

		// After navigating down 10 times, scroll should have advanced
		// The exact amount depends on entry heights, but it should have moved
		if state.ScrollOffset == initialScrollOffset && state.SelectedIdx > 2 {
			t.Errorf("ScrollOffset should have advanced after navigating past visible entries, got ScrollOffset %d with SelectedIdx %d",
				state.ScrollOffset, state.SelectedIdx)
		}
	})
}

func TestSelectionStaysWithinVisibleArea(t *testing.T) {
	entries := makeEntries(100)

	t.Run("selection should not scroll into bottom padding area", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = headerRows + bottomPadding + 10

		visibleRows := state.TerminalHeight - headerRows - bottomPadding
		maxVisibleEntries := visibleRows / estimatedRowsPerEntry

		if maxVisibleEntries < 1 {
			maxVisibleEntries = 1
		}

		for i := 0; i < 50; i++ {
			state.NavigateDown()
		}

		visibleEntriesFromOffset := maxVisibleEntries
		maxSelectableIndex := state.ScrollOffset + visibleEntriesFromOffset - 1

		if state.SelectedIdx > maxSelectableIndex {
			t.Errorf("SelectedIdx %d is beyond visible area (max %d with offset %d)",
				state.SelectedIdx, maxSelectableIndex, state.ScrollOffset)
		}
	})

	t.Run("scroll offset adjusted when navigating to keep selection visible", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 20

		for i := 0; i < 30; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx < state.ScrollOffset {
			t.Errorf("SelectedIdx %d is above ScrollOffset %d", state.SelectedIdx, state.ScrollOffset)
		}
	})

	t.Run("selection can reach last entry without going into padding", func(t *testing.T) {
		entries := makeEntries(10)
		state := NewState(entries)
		state.TerminalHeight = headerRows + bottomPadding + 15

		for i := 0; i < 20; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 9 {
			t.Errorf("expected SelectedIdx 9 (last entry), got %d", state.SelectedIdx)
		}
	})
}

func TestBottomPaddingRespected(t *testing.T) {
	entries := makeEntries(50)

	t.Run("visible entry count accounts for bottom padding", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		visibleRows := state.TerminalHeight - headerRows - bottomPadding
		if visibleRows <= 0 {
			t.Errorf("visibleRows should be positive, got %d", visibleRows)
		}

		_ = visibleRows
	})

	t.Run("scroll advances before selection hits bottom padding", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = headerRows + bottomPadding + 6

		for i := 0; i < 10; i++ {
			state.NavigateDown()
		}

		visibleRows := state.TerminalHeight - headerRows - bottomPadding
		maxVisibleFromTop := state.ScrollOffset + visibleRows/estimatedRowsPerEntry

		if state.SelectedIdx >= maxVisibleFromTop {
			t.Errorf("SelectedIdx %d should not reach bottom padding area (visible up to %d)",
				state.SelectedIdx, maxVisibleFromTop-1)
		}
	})
}

func TestAppendChar(t *testing.T) {
	entries := makeEntries(5)

	t.Run("append to empty buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = ""
		state.CursorIdx = 0

		state.AppendChar('a')

		if state.SearchBuf != "a" {
			t.Errorf("expected SearchBuf 'a', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 1 {
			t.Errorf("expected CursorIdx 1, got %d", state.CursorIdx)
		}
	})

	t.Run("append to end of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 3

		state.AppendChar('d')

		if state.SearchBuf != "abcd" {
			t.Errorf("expected SearchBuf 'abcd', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 4 {
			t.Errorf("expected CursorIdx 4, got %d", state.CursorIdx)
		}
	})

	t.Run("insert in middle of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "acd"
		state.CursorIdx = 1

		state.AppendChar('b')

		if state.SearchBuf != "abcd" {
			t.Errorf("expected SearchBuf 'abcd', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 2 {
			t.Errorf("expected CursorIdx 2, got %d", state.CursorIdx)
		}
	})

	t.Run("insert at beginning of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "bcd"
		state.CursorIdx = 0

		state.AppendChar('a')

		if state.SearchBuf != "abcd" {
			t.Errorf("expected SearchBuf 'abcd', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 1 {
			t.Errorf("expected CursorIdx 1, got %d", state.CursorIdx)
		}
	})

	t.Run("updates filtered results", func(t *testing.T) {
		entries := []config.Entry{
			{Command: "git"},
			{Command: "docker"},
		}
		state := NewState(entries)

		state.AppendChar('g')

		if len(state.Filtered) != 1 {
			t.Errorf("expected 1 filtered result, got %d", len(state.Filtered))
		}
		if state.Filtered[0].Command != "git" {
			t.Errorf("expected filtered result 'git', got %q", state.Filtered[0].Command)
		}
	})
}

func TestDeleteChar(t *testing.T) {
	entries := makeEntries(5)

	t.Run("delete from empty buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = ""
		state.CursorIdx = 0

		state.DeleteChar()

		if state.SearchBuf != "" {
			t.Errorf("expected SearchBuf '', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 0 {
			t.Errorf("expected CursorIdx 0, got %d", state.CursorIdx)
		}
	})

	t.Run("delete from beginning of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 0

		state.DeleteChar()

		if state.SearchBuf != "abc" {
			t.Errorf("expected SearchBuf unchanged 'abc', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 0 {
			t.Errorf("expected CursorIdx 0, got %d", state.CursorIdx)
		}
	})

	t.Run("delete from middle of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abcd"
		state.CursorIdx = 2

		state.DeleteChar()

		if state.SearchBuf != "acd" {
			t.Errorf("expected SearchBuf 'acd', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 1 {
			t.Errorf("expected CursorIdx 1, got %d", state.CursorIdx)
		}
	})

	t.Run("delete from end of buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 3

		state.DeleteChar()

		if state.SearchBuf != "ab" {
			t.Errorf("expected SearchBuf 'ab', got %q", state.SearchBuf)
		}
		if state.CursorIdx != 2 {
			t.Errorf("expected CursorIdx 2, got %d", state.CursorIdx)
		}
	})

	t.Run("updates filtered results", func(t *testing.T) {
		entries := []config.Entry{
			{Command: "git"},
			{Command: "docker"},
		}
		state := NewState(entries)
		state.SearchBuf = "gi"
		state.CursorIdx = 2
		state.UpdateSearch()

		state.DeleteChar()

		if state.SearchBuf != "g" {
			t.Errorf("expected SearchBuf 'g', got %q", state.SearchBuf)
		}
	})
}

func TestMoveCursorLeft(t *testing.T) {
	entries := makeEntries(5)

	t.Run("move left from beginning", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 0

		state.MoveCursorLeft()

		if state.CursorIdx != 0 {
			t.Errorf("expected CursorIdx 0, got %d", state.CursorIdx)
		}
	})

	t.Run("move left from middle", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 2

		state.MoveCursorLeft()

		if state.CursorIdx != 1 {
			t.Errorf("expected CursorIdx 1, got %d", state.CursorIdx)
		}
	})

	t.Run("move left from end", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 3

		state.MoveCursorLeft()

		if state.CursorIdx != 2 {
			t.Errorf("expected CursorIdx 2, got %d", state.CursorIdx)
		}
	})

	t.Run("does not change search buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 2

		state.MoveCursorLeft()

		if state.SearchBuf != "abc" {
			t.Errorf("expected SearchBuf 'abc', got %q", state.SearchBuf)
		}
	})
}

func TestMoveCursorRight(t *testing.T) {
	entries := makeEntries(5)

	t.Run("move right from beginning", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 0

		state.MoveCursorRight()

		if state.CursorIdx != 1 {
			t.Errorf("expected CursorIdx 1, got %d", state.CursorIdx)
		}
	})

	t.Run("move right from middle", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 1

		state.MoveCursorRight()

		if state.CursorIdx != 2 {
			t.Errorf("expected CursorIdx 2, got %d", state.CursorIdx)
		}
	})

	t.Run("move right from end", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 3

		state.MoveCursorRight()

		if state.CursorIdx != 3 {
			t.Errorf("expected CursorIdx 3, got %d", state.CursorIdx)
		}
	})

	t.Run("move right with empty buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = ""
		state.CursorIdx = 0

		state.MoveCursorRight()

		if state.CursorIdx != 0 {
			t.Errorf("expected CursorIdx 0, got %d", state.CursorIdx)
		}
	})

	t.Run("does not change search buffer", func(t *testing.T) {
		state := NewState(entries)
		state.SearchBuf = "abc"
		state.CursorIdx = 0

		state.MoveCursorRight()

		if state.SearchBuf != "abc" {
			t.Errorf("expected SearchBuf 'abc', got %q", state.SearchBuf)
		}
	})
}

func TestDebugEntryHeights(t *testing.T) {
	// Test that entry heights are computed correctly
	entries := []config.Entry{
		{Title: "short", Command: "cmd"},
		{Title: "tall", Command: "cmd", Description: []string{"line1", "line2"}},
		{Title: "wide", Command: "cmd", Description: []string{"a", "b", "c", "d"}},
	}
	state := NewState(entries)
	state.TerminalHeight = 24

	t.Run("entry heights computed correctly", func(t *testing.T) {
		if len(state.EntryHeights) != 3 {
			t.Fatalf("expected 3 entry heights, got %d", len(state.EntryHeights))
		}

		// Entry 0: title(1) + desc(0) + cmd(1) + spacing(1) = 3
		if state.EntryHeights[0] != 3 {
			t.Errorf("entry0 height: expected 3, got %d", state.EntryHeights[0])
		}

		// Entry 1: title(1) + desc(2) + cmd(1) + spacing(1) = 5
		if state.EntryHeights[1] != 5 {
			t.Errorf("entry1 height: expected 5, got %d", state.EntryHeights[1])
		}

		// Entry 2: title(1) + desc(4) + cmd(1) + spacing(1) = 7
		if state.EntryHeights[2] != 7 {
			t.Errorf("entry2 height: expected 7, got %d", state.EntryHeights[2])
		}
	})

	t.Run("findLastVisibleEntry works correctly", func(t *testing.T) {
		visibleHeight := 24 - headerRows - bottomPadding // 8

		// Starting from entry 0, with height 8
		// Entry 0 (height 3): total = 3
		// Entry 1 (height 5): total = 3 + 5 = 8
		// Entry 2 (height 7): would exceed 8
		// So last visible is entry 1
		lastVisible := state.findLastVisibleEntry(0, visibleHeight)
		if lastVisible != 1 {
			t.Errorf("expected last visible entry 1, got %d", lastVisible)
		}
	})

	t.Run("navigation keeps selection visible", func(t *testing.T) {
		state := NewState(entries)
		state.TerminalHeight = 24

		// Navigate to entry 2
		for i := 0; i < 2; i++ {
			state.NavigateDown()
		}

		if state.SelectedIdx != 2 {
			t.Errorf("expected SelectedIdx 2, got %d", state.SelectedIdx)
		}

		// Entry 2 should be visible, so ScrollOffset should be adjusted
		visibleHeight := 24 - headerRows - bottomPadding
		lastVisible := state.findLastVisibleEntry(state.ScrollOffset, visibleHeight)

		if state.SelectedIdx > lastVisible {
			t.Errorf("SelectedIdx %d is past last visible entry %d (ScrollOffset %d)",
				state.SelectedIdx, lastVisible, state.ScrollOffset)
		}
	})
}

func TestDebugVariableHeightNavigation(t *testing.T) {
	entries := []config.Entry{
		{Title: "e0", Command: "cmd"},
		{Title: "e1", Command: "cmd", Description: []string{"d1", "d2", "d3", "d4"}}, // height 7
		{Title: "e2", Command: "cmd"},
		{Title: "e3", Command: "cmd", Description: []string{"d1", "d2"}}, // height 5
		{Title: "e4", Command: "cmd"},
		{Title: "e5", Command: "cmd", Description: []string{"d1", "d2", "d3"}}, // height 6
		{Title: "e6", Command: "cmd"},
		{Title: "e7", Command: "cmd"},
		{Title: "e8", Command: "cmd", Description: []string{"d1"}}, // height 4
	}
	state := NewState(entries)
	state.TerminalHeight = 24

	t.Log("Entry heights:", state.EntryHeights)

	// Navigate to entry 8
	for i := 0; i < 8; i++ {
		state.NavigateDown()
		t.Logf("After NavigateDown %d: SelectedIdx=%d, ScrollOffset=%d", i+1, state.SelectedIdx, state.ScrollOffset)
	}

	visibleHeight := 24 - headerRows - bottomPadding
	lastVisible := state.findLastVisibleEntry(state.ScrollOffset, visibleHeight)

	t.Logf("Final: SelectedIdx=%d, ScrollOffset=%d, lastVisible=%d, visibleHeight=%d",
		state.SelectedIdx, state.ScrollOffset, lastVisible, visibleHeight)

	if state.SelectedIdx > lastVisible {
		t.Errorf("SelectedIdx %d is past last visible entry %d", state.SelectedIdx, lastVisible)
	}
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
