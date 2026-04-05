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
	EntryHeights   []int
	TerminalHeight int
	TerminalWidth  int
	ScrollOffset   int
}

func NewState(entries []config.Entry) *State {
	state := &State{
		Mode:         ModeInsert,
		SearchBuf:    "",
		CursorIdx:    0,
		SelectedIdx:  0,
		Entries:      entries,
		Filtered:     entries,
		ScrollOffset: 0,
	}
	state.computeEntryHeights()
	return state
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
	s.computeEntryHeights()
}

func (s *State) computeEntryHeights() {
	s.EntryHeights = make([]int, len(s.Filtered))
	for i, entry := range s.Filtered {
		s.EntryHeights[i] = computeEntryHeight(entry)
	}
}

func computeEntryHeight(entry config.Entry) int {
	height := 0
	if entry.Title != "" {
		height++
	}
	height += len(entry.Description)
	height++ // command line
	height++ // spacing after entry
	return height
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

const (
	headerTitleRows     = 2
	headerPromptRows    = 2
	headerSeparatorRows = 2
	headerRows          = headerTitleRows + headerPromptRows + headerSeparatorRows

	bottomPadding         = 10
	estimatedRowsPerEntry = 3 // Deprecated: kept for tests, not used in navigation
	scrollIndicatorRows   = 2 // Scroll indicator + gap before it
)

func (s *State) NavigateDown() {
	if s.SelectedIdx < len(s.Filtered)-1 {
		s.SelectedIdx++
		s.ensureSelectionVisible()
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

func (s *State) ensureSelectionVisible() {
	visibleHeight := s.TerminalHeight - headerRows - bottomPadding - scrollIndicatorRows
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	for s.ScrollOffset < s.SelectedIdx {
		lastVisible := s.findLastVisibleEntry(s.ScrollOffset, visibleHeight)
		if lastVisible >= s.SelectedIdx {
			break
		}
		s.ScrollOffset++
	}
}

func (s *State) findLastVisibleEntry(startIdx, visibleHeight int) int {
	row := 0
	for i := startIdx; i < len(s.EntryHeights); i++ {
		entryHeight := s.EntryHeights[i]
		if row+entryHeight > visibleHeight {
			if i == startIdx {
				return startIdx
			}
			return i - 1
		}
		row += entryHeight
	}
	return len(s.EntryHeights) - 1
}

func (s *State) SwitchToNormal() {
	s.Mode = ModeNormal
}

func (s *State) SwitchToInsert(cursorPos string) {
	s.Mode = ModeInsert
	switch cursorPos {
	case "before":
		// 'i': Keep cursor as-is, insert before current position
	case "after":
		// 'a': Move right by one, insert after current position
		if s.CursorIdx < len(s.SearchBuf) {
			s.CursorIdx++
		}
	case "start":
		// 'I': Move to beginning of line
		s.CursorIdx = 0
	case "end":
		// 'A': Move to end of line
		s.CursorIdx = len(s.SearchBuf)
	}
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
