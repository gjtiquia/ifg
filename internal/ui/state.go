package ui

import (
	"github.com/gjtiquia/ifg/internal/config"
	"github.com/gjtiquia/ifg/internal/search"
)

type Mode int

const (
	ModeInsert Mode = iota
	ModeNormal
)

type State struct {
	Mode           Mode
	SearchBuf      string
	CursorIdx      int
	SelectedIdx    int
	Entries        []config.Entry
	Filtered       []config.Entry
	TerminalHeight int
	TerminalWidth  int
	ScrollOffset   int
}

func NewState(entries []config.Entry) *State {
	return &State{
		Mode:         ModeInsert,
		SearchBuf:    "",
		CursorIdx:    0,
		SelectedIdx:  0,
		Entries:      entries,
		Filtered:     entries,
		ScrollOffset: 0,
	}
}

func (s *State) UpdateSearch() {
	s.Filtered = search.Match(s.Entries, s.SearchBuf)
	if s.SelectedIdx >= len(s.Filtered) {
		s.SelectedIdx = len(s.Filtered) - 1
		if s.SelectedIdx < 0 {
			s.SelectedIdx = 0
		}
	}
	s.ScrollOffset = 0
}

func (s *State) AppendChar(ch rune) {
	s.SearchBuf = s.SearchBuf[:s.CursorIdx] + string(ch) + s.SearchBuf[s.CursorIdx:]
	s.CursorIdx++
	s.UpdateSearch()
}

func (s *State) DeleteChar() {
	if s.CursorIdx > 0 {
		s.SearchBuf = s.SearchBuf[:s.CursorIdx-1] + s.SearchBuf[s.CursorIdx:]
		s.CursorIdx--
		s.UpdateSearch()
	}
}

func (s *State) NavigateUp() {
	if s.SelectedIdx > 0 {
		s.SelectedIdx--
		if s.SelectedIdx < s.ScrollOffset {
			s.ScrollOffset = s.SelectedIdx
		}
	}
}

func (s *State) NavigateDown() {
	if s.SelectedIdx < len(s.Filtered)-1 {
		s.SelectedIdx++
		visibleHeight := s.TerminalHeight - 2
		if s.SelectedIdx >= s.ScrollOffset+visibleHeight {
			s.ScrollOffset = s.SelectedIdx - visibleHeight + 1
		}
	}
}

func (s *State) SwitchToNormal() {
	s.Mode = ModeNormal
}

func (s *State) SwitchToInsert(cursorPos string) {
	s.Mode = ModeInsert
	if cursorPos == "start" {
		s.CursorIdx = 0
	} else if cursorPos == "end" {
		s.CursorIdx = len(s.SearchBuf)
	}
}

func (s *State) GetSelectedCommand() string {
	if len(s.Filtered) > 0 && s.SelectedIdx >= 0 && s.SelectedIdx < len(s.Filtered) {
		return s.Filtered[s.SelectedIdx].Command
	}
	return ""
}
