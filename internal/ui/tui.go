package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	originalState *term.State
}

func SetupTerminal() (*Terminal, error) {
	fd := int(os.Stdin.Fd())
	originalState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}

	return &Terminal{originalState: originalState}, nil
}

func (t *Terminal) Restore() error {
	fd := int(os.Stdin.Fd())
	return term.Restore(fd, t.originalState)
}

func (t *Terminal) GetSize() (int, int, error) {
	fd := int(os.Stdin.Fd())
	return term.GetSize(fd)
}

func ClearScreen() {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%dH", row, col)
}

func Render(state *State) {
	ClearScreen()

	searchPrompt := "search: "
	if state.Mode == ModeNormal {
		searchPrompt = "normal: "
	}

	MoveCursor(1, 1)
	fmt.Printf("\x1b[1m%s\x1b[0m%s", searchPrompt, state.SearchBuf)

	if len(state.Filtered) == 0 {
		MoveCursor(3, 1)
		fmt.Print("\x1b[90mNo results found\x1b[0m")
		return
	}

	visibleHeight := state.TerminalHeight - 2
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	for i := 0; i < visibleHeight && i+state.ScrollOffset < len(state.Filtered); i++ {
		entryIdx := i + state.ScrollOffset
		entry := state.Filtered[entryIdx]

		row := i + 3
		MoveCursor(row, 1)

		prefix := "  "
		if entryIdx == state.SelectedIdx {
			prefix = "\x1b[32;1m>\x1b[0m "
			fmt.Print("\x1b[1m")
		}

		if len(entry.Title) > 0 {
			fmt.Printf("%s%s", prefix, entry.Title)
		} else {
			fmt.Printf("%s%s", prefix, entry.Command)
		}

		if entryIdx == state.SelectedIdx {
			fmt.Print("\x1b[0m")
		}
	}

	if state.ScrollOffset > 0 || state.ScrollOffset+visibleHeight < len(state.Filtered) {
		MoveCursor(state.TerminalHeight, 1)
		fmt.Printf("\x1b[90m[%d-%d of %d]\x1b[0m", state.ScrollOffset+1, min(state.ScrollOffset+visibleHeight, len(state.Filtered)), len(state.Filtered))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RenderCommand(command string) {
	ClearScreen()
	MoveCursor(1, 1)
	fmt.Print(command)
}
