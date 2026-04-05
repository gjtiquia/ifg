package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

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

func (t *Terminal) Screen() tcell.Screen {
	return t.screen
}

func (t *Terminal) WrappedScreen() Screen {
	return NewTcellScreen(t.screen)
}

func Render(state *State, screen Screen) {
	screen.Clear()

	row := 0

	header := "ifg - [i] [f]or[g]ot"
	drawText(screen, 0, row, header, Style{Bold: true})
	row += 2

	prompt := "type to search: "
	if state.Mode == ModeNormal {
		prompt = "search results for: "
	}

	if len(state.SearchBuf) > 0 {
		drawText(screen, 0, row, prompt, NewStyle())
		x := len(prompt)

		for i, ch := range state.SearchBuf {
			if i == state.CursorIdx {
				screen.SetContent(x+i, row, ch, Style{Reverse: true})
			} else {
				screen.SetContent(x+i, row, ch, NewStyle())
			}
		}

		if state.CursorIdx == len(state.SearchBuf) {
			screen.SetContent(x+state.CursorIdx, row, ' ', Style{Reverse: true})
		}
	} else {
		drawText(screen, 0, row, prompt, NewStyle())
		if state.Mode == ModeInsert {
			screen.SetContent(len(prompt), row, ' ', Style{Reverse: true})
		}
	}
	row += 2

	drawText(screen, 0, row, "---", NewStyle())
	row += 2

	if len(state.Filtered) == 0 {
		drawText(screen, 0, row, "No results found", Style{Dim: true})
		screen.Show()
		return
	}

	_, height := screen.Size()
	maxRow := height - bottomPadding
	if maxRow < row+1 {
		maxRow = row + 1
	}

	// Reserve one row for scroll indicator
	scrollIndicatorRow := maxRow - 1
	contentEndRow := scrollIndicatorRow - 1
	if contentEndRow < row+1 {
		contentEndRow = row + 1
	}

	var lastVisibleIdx int
	for i := 0; i+state.ScrollOffset < len(state.Filtered); i++ {
		if row >= contentEndRow {
			break
		}
		entryIdx := i + state.ScrollOffset
		entry := state.Filtered[entryIdx]
		lastVisibleIdx = entryIdx

		isSelected := entryIdx == state.SelectedIdx
		style := NewStyle()
		if isSelected {
			style.Bold = true
		}

		prefix := "  "
		if isSelected {
			prefix = "> "
		}

		if entry.Title != "" {
			drawText(screen, 0, row, prefix+"# "+entry.Title, style)
			row++
		}

		for _, desc := range entry.Description {
			if row >= contentEndRow {
				break
			}
			descPrefix := "  "
			if isSelected {
				descPrefix = "> "
			}
			drawText(screen, 0, row, descPrefix+"# "+desc, style)
			row++
		}

		if row >= contentEndRow {
			break
		}
		cmdPrefix := "  "
		if isSelected {
			cmdPrefix = "> "
		}
		drawText(screen, 0, row, cmdPrefix+entry.Command, style)
		row += 2
	}

	if state.ScrollOffset > 0 || lastVisibleIdx < len(state.Filtered)-1 {
		scrollText := "[" + strconv.Itoa(state.ScrollOffset+1) + "-" + strconv.Itoa(lastVisibleIdx+1) + " of " + strconv.Itoa(len(state.Filtered)) + "]"
		drawText(screen, 0, scrollIndicatorRow, scrollText, Style{Dim: true})
	}

	screen.Show()
}

func drawText(screen Screen, x, y int, text string, style Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, style)
	}
}
