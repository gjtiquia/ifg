package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gjtiquia/ifg/internal/config"
	"github.com/gjtiquia/ifg/internal/ui"
)

func main() {
	configPath := config.GetConfigPath()

	var entries []config.Entry
	var err error

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.CreateDefaultConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			os.Exit(2)
		}
	}

	entries, err = config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(2)
	}

	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "No entries found in config file\n")
		os.Exit(2)
	}

	term, err := ui.SetupTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up terminal: %v\n", err)
		os.Exit(2)
	}
	defer term.Restore()

	ui.HideCursor()
	defer ui.ShowCursor()

	width, height, err := term.GetSize()
	if err != nil {
		width = 80
		height = 24
	}

	state := ui.NewState(entries)
	state.TerminalWidth = width
	state.TerminalHeight = height

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)
	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGWINCH {
				width, height, err := term.GetSize()
				if err == nil {
					state.TerminalWidth = width
					state.TerminalHeight = height
					ui.Render(state)
				}
			} else {
				term.Restore()
				ui.ShowCursor()
				os.Exit(1)
			}
		}
	}()

	selectedCommand := runInputLoop(state, term)

	if selectedCommand != "" {
		fmt.Println(selectedCommand)
		os.Exit(0)
	}

	os.Exit(1)
}

func runInputLoop(state *ui.State, term *ui.Terminal) string {
	ui.Render(state)

	for {
		key, err := ui.ReadKey()
		if err != nil {
			continue
		}

		switch state.Mode {
		case ui.ModeInsert:
			switch key.Type {
			case ui.KeyChar:
				state.AppendChar(key.Char)
			case ui.KeyBackspace:
				state.DeleteChar()
			case ui.KeyUp:
				state.NavigateUp()
			case ui.KeyDown:
				state.NavigateDown()
			case ui.KeyEnter:
				return state.GetSelectedCommand()
			case ui.KeyEscape:
				state.SwitchToNormal()
			case ui.KeyCtrlC:
				return ""
			}

		case ui.ModeNormal:
			switch key.Type {
			case ui.KeyChar:
				switch key.Char {
				case 'j':
					state.NavigateDown()
				case 'k':
					state.NavigateUp()
				case 'i', 'I':
					state.SwitchToInsert("start")
				case 'a', 'A':
					state.SwitchToInsert("end")
				}
			case ui.KeyUp:
				state.NavigateUp()
			case ui.KeyDown:
				state.NavigateDown()
			case ui.KeyEnter:
				return state.GetSelectedCommand()
			case ui.KeyEscape:
				return ""
			case ui.KeyCtrlC:
				return ""
			}
		}

		ui.Render(state)
	}
}
