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

func EnterAlternateScreen() {
	fmt.Print("\x1b[?1049h")
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")
	HideCursor()
}

func ExitAlternateScreen() {
	ShowCursor()
	fmt.Print("\x1b[?1049l")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func Render(state *State) {
	fmt.Print("\x1b[2J\x1b[H")

	fmt.Println("\x1b[1mifg - [i] [f]or[g]ot that cmd again\x1b[0m")
	fmt.Println()

	prompt := "type to search: "
	if state.Mode == ModeNormal {
		prompt = "normal mode: "
	}
	fmt.Printf("%s%s\n", prompt, state.SearchBuf)
	fmt.Println()

	fmt.Println("---")
	fmt.Println()

	if len(state.Filtered) == 0 {
		fmt.Println("\x1b[90mNo results found\x1b[0m")
		return
	}

	visibleHeight := state.TerminalHeight - 8
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	for i := 0; i < visibleHeight && i+state.ScrollOffset < len(state.Filtered); i++ {
		entryIdx := i + state.ScrollOffset
		entry := state.Filtered[entryIdx]

		isSelected := entryIdx == state.SelectedIdx
		prefix := "  "
		if isSelected {
			prefix = "> "
		}

		if isSelected {
			fmt.Print("\x1b[1m")
		}

		if entry.Title != "" {
			fmt.Printf("%s# %s\n", prefix, entry.Title)
		}

		for _, desc := range entry.Description {
			fmt.Printf("%s# %s\n", prefix, desc)
		}

		fmt.Printf("%s%s\n", prefix, entry.Command)

		if isSelected {
			fmt.Print("\x1b[0m")
		}

		fmt.Println()
	}

	if state.ScrollOffset > 0 || state.ScrollOffset+visibleHeight < len(state.Filtered) {
		fmt.Printf("\x1b[90m[%d-%d of %d]\x1b[0m\n", state.ScrollOffset+1, min(state.ScrollOffset+visibleHeight, len(state.Filtered)), len(state.Filtered))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
