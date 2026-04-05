package ui

import (
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

func Render(state *State, screen tcell.Screen) {
	screen.Clear()

	row := 0

	header := "ifg - [i] [f]or[g]ot"
	drawText(screen, 0, row, header, tcell.StyleDefault.Bold(true))
	row += 2

	prompt := "type to search: "
	if state.Mode == ModeNormal {
		prompt = "search results for: "
	}

	// Draw search buffer with cursor (both Insert and Normal modes)
	if len(state.SearchBuf) > 0 {
		drawText(screen, 0, row, prompt, tcell.StyleDefault)
		x := len(prompt)

		for i, ch := range state.SearchBuf {
			if i == state.CursorIdx {
				screen.SetContent(x+i, row, ch, nil, tcell.StyleDefault.Reverse(true))
			} else {
				screen.SetContent(x+i, row, ch, nil, tcell.StyleDefault)
			}
		}

		if state.CursorIdx == len(state.SearchBuf) {
			screen.SetContent(x+state.CursorIdx, row, ' ', nil, tcell.StyleDefault.Reverse(true))
		}
	} else {
		// Empty buffer
		drawText(screen, 0, row, prompt, tcell.StyleDefault)
		if state.Mode == ModeInsert {
			screen.SetContent(len(prompt), row, ' ', nil, tcell.StyleDefault.Reverse(true))
		}
	}
	row += 2

	drawText(screen, 0, row, "---", tcell.StyleDefault)
	row += 2

	if len(state.Filtered) == 0 {
		drawText(screen, 0, row, "No results found", tcell.StyleDefault.Dim(true))
		screen.Show()
		return
	}

	_, height := screen.Size()
	maxRow := height - bottomPadding
	if maxRow < row+1 {
		maxRow = row + 1
	}

	var lastVisibleIdx int
	for i := 0; i+state.ScrollOffset < len(state.Filtered); i++ {
		if row >= maxRow {
			break
		}
		entryIdx := i + state.ScrollOffset
		entry := state.Filtered[entryIdx]
		lastVisibleIdx = entryIdx

		isSelected := entryIdx == state.SelectedIdx
		style := tcell.StyleDefault
		if isSelected {
			style = style.Bold(true)
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
			if row >= maxRow {
				break
			}
			descPrefix := "  "
			if isSelected {
				descPrefix = "> "
			}
			drawText(screen, 0, row, descPrefix+"# "+desc, style)
			row++
		}

		if row >= maxRow {
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
		scrollText := "["
		n := state.ScrollOffset + 1
		if n < 10 {
			scrollText += string(rune('0' + n))
		} else {
			scrollText += string(rune('0' + n/10))
			scrollText += string(rune('0' + n%10))
		}
		scrollText += "-"
		endIdx := lastVisibleIdx + 1
		if endIdx < 10 {
			scrollText += string(rune('0' + endIdx))
		} else {
			scrollText += string(rune('0' + endIdx/10))
			scrollText += string(rune('0' + endIdx%10))
		}
		scrollText += " of "
		total := len(state.Filtered)
		if total < 10 {
			scrollText += string(rune('0' + total))
		} else {
			scrollText += string(rune('0' + total/10))
			scrollText += string(rune('0' + total%10))
		}
		scrollText += "]"
		drawText(screen, 0, maxRow-1, scrollText, tcell.StyleDefault.Dim(true))
	}

	screen.Show()
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
