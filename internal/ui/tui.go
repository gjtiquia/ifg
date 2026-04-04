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

	header := "ifg - [i] [f]or[g]ot that cmd again"
	drawText(screen, 0, row, header, tcell.StyleDefault.Bold(true))
	row += 2

	prompt := "type to search: "
	if state.Mode == ModeNormal {
		prompt = "normal mode: "
	}
	drawText(screen, 0, row, prompt+state.SearchBuf, tcell.StyleDefault)
	row += 2

	drawText(screen, 0, row, "---", tcell.StyleDefault)
	row += 2

	if len(state.Filtered) == 0 {
		drawText(screen, 0, row, "No results found", tcell.StyleDefault.Dim(true))
		screen.Show()
		return
	}

	_, height := screen.Size()
	visibleHeight := height - row - 2
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	for i := 0; i < visibleHeight && i+state.ScrollOffset < len(state.Filtered); i++ {
		entryIdx := i + state.ScrollOffset
		entry := state.Filtered[entryIdx]

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
			descPrefix := "  "
			if isSelected {
				descPrefix = "> "
			}
			drawText(screen, 0, row, descPrefix+"# "+desc, style)
			row++
		}

		cmdPrefix := "  "
		if isSelected {
			cmdPrefix = "> "
		}
		drawText(screen, 0, row, cmdPrefix+entry.Command, style)
		row += 2
	}

	if state.ScrollOffset > 0 || state.ScrollOffset+visibleHeight < len(state.Filtered) {
		scrollIndicator := ""
		scrollIndicator += "["
		scrollIndicator += string(rune('0' + (state.ScrollOffset+1)/10))
		scrollIndicator += string(rune('0' + (state.ScrollOffset+1)%10))
		scrollIndicator += "-"
		scrollIndicator += string(rune('0' + min(state.ScrollOffset+visibleHeight, len(state.Filtered))/10))
		scrollIndicator += string(rune('0' + min(state.ScrollOffset+visibleHeight, len(state.Filtered))%10))
		scrollIndicator += " of "
		scrollIndicator += string(rune('0' + len(state.Filtered)/10))
		scrollIndicator += string(rune('0' + len(state.Filtered)%10))
		scrollIndicator += "]"
		scrollText := ""
		scrollText += "["
		n := state.ScrollOffset + 1
		if n < 10 {
			scrollText += string(rune('0' + n))
		} else {
			scrollText += string(rune('0' + n/10))
			scrollText += string(rune('0' + n%10))
		}
		scrollText += "-"
		endIdx := min(state.ScrollOffset+visibleHeight, len(state.Filtered))
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
		drawText(screen, 0, height-1, scrollText, tcell.StyleDefault.Dim(true))
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
