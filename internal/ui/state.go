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
	// 'i/I': Insert BEFORE cursor - keep cursor at current position
	// 'a/A': Insert AFTER cursor - move cursor right by one
	if cursorPos == "end" {
		// Move cursor right to insert after current position
		if s.CursorIdx < len(s.SearchBuf) {
			s.CursorIdx++
		}
	}
	// If "start", keep CursorIdx as-is (insert before current position)
}

func (s *State) GetSelectedCommand() string {
	if len(s.Filtered) > 0 && s.SelectedIdx >= 0 && s.SelectedIdx < len(s.Filtered) {
		return s.Filtered[s.SelectedIdx].Command
	}
	return ""
}

func (s *State) MoveCursorLeft() {
	if s.CursorIdx > 0 {
		s.CursorIdx--
	}
}

func (s *State) MoveCursorRight() {
	if s.CursorIdx < len(s.SearchBuf) {
		s.CursorIdx++
	}
}

func (s *State) MoveWordForward() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	for i < len(buf) && isWordChar(buf[i]) {
		i++
	}

	for i < len(buf) && !isWordChar(buf[i]) {
		i++
	}

	s.CursorIdx = i
}

func (s *State) MoveWORDForward() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	for i < len(buf) && !isSpace(buf[i]) {
		i++
	}

	for i < len(buf) && isSpace(buf[i]) {
		i++
	}

	s.CursorIdx = i
}

func (s *State) MoveWordBackward() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	for i > 0 && !isWordChar(buf[i-1]) {
		i--
	}

	for i > 0 && isWordChar(buf[i-1]) {
		i--
	}

	s.CursorIdx = i
}

func (s *State) MoveWORDBackward() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	for i > 0 && isSpace(buf[i-1]) {
		i--
	}

	for i > 0 && !isSpace(buf[i-1]) {
		i--
	}

	s.CursorIdx = i
}

func (s *State) MoveWordEnd() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	if i < len(buf) {
		i++
	}

	for i < len(buf) && !isWordChar(buf[i]) {
		i++
	}

	for i < len(buf) && isWordChar(buf[i]) {
		i++
	}

	if i > 0 {
		s.CursorIdx = i - 1
	} else {
		s.CursorIdx = 0
	}
}

func (s *State) MoveWORDEnd() {
	buf := []rune(s.SearchBuf)
	i := s.CursorIdx

	if i < len(buf) {
		i++
	}

	for i < len(buf) && isSpace(buf[i]) {
		i++
	}

	for i < len(buf) && !isSpace(buf[i]) {
		i++
	}

	if i > 0 {
		s.CursorIdx = i - 1
	} else {
		s.CursorIdx = 0
	}
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
