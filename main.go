package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gjtiquia/ifg/internal/config"
	"github.com/gjtiquia/ifg/internal/ui"
)

//go:embed shell/ifg.sh
var shellWrapper string

//go:embed shell/config.sh
var defaultConfigContent string

func printHelp() {
	fmt.Println("# ifg - [i] [f]or[g]ot")
	fmt.Println()
	fmt.Println("for when you are trying to rmb that command")
	fmt.Println()
	fmt.Println("## usage:")
	fmt.Println()
	fmt.Println("  ifg [flags]")
	fmt.Println()
	fmt.Println("## config dir:")
	fmt.Println()
	fmt.Println("  " + config.GetConfigDir())
	fmt.Println()
	fmt.Println("## flags:")
	fmt.Println()
	fmt.Println("  --help    show this help")
	fmt.Println("  --sh      print shell integration code")
	fmt.Println("  --bash    alias for --sh")
	fmt.Println("  --zsh     alias for --sh")
	fmt.Println()
	fmt.Println("## shell integration:")
	fmt.Println()
	fmt.Println("  add to ~/.bashrc or ~/.zshrc:")
	fmt.Println("  eval \"$(ifg --sh)\"")
	fmt.Println()
	fmt.Println("## repo:")
	fmt.Println()
	fmt.Println("  github.com/gjtiquia/ifg")
	fmt.Println()
	fmt.Println("## license:")
	fmt.Println()
	fmt.Println("  MIT")
	fmt.Println()
}

func main() {
	// Check for shell integration flags
	if len(os.Args) > 1 {
		switch os.Args[1] {

		case "--sh", "--bash", "--zsh":
			fmt.Print(shellWrapper)
			os.Exit(0)

		case "--help", "-h":
			printHelp()
			os.Exit(0)

		default:
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	}

	config.SetDefaultConfig(defaultConfigContent)
	configDir := config.GetConfigDir()

	var entries []config.Entry
	var err error

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := config.CreateDefaultConfig(configDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			os.Exit(2)
		}
	}

	entries, err = config.LoadConfig(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(2)
	}

	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "No entries found in config\n")
		os.Exit(2)
	}

	term, err := ui.SetupTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up terminal: %v\n", err)
		os.Exit(2)
	}
	defer term.Restore()

	width, height := term.GetSize()
	state := ui.NewState(entries)
	state.TerminalWidth = width
	state.TerminalHeight = height

	selectedCommand := runInputLoop(state, term)

	term.Restore()

	if selectedCommand != "" {
		fmt.Println(selectedCommand)
		os.Exit(0)
	}

	os.Exit(1)
}

func runInputLoop(state *ui.State, term *ui.Terminal) string {
	ui.Render(state, term.WrappedScreen())

	for {
		key := ui.ReadKey(term.Screen())

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
				case 'h':
					state.MoveCursorLeft()
				case 'l':
					state.MoveCursorRight()
				case 'w':
					state.MoveWordForward()
				case 'W':
					state.MoveWORDForward()
				case 'b':
					state.MoveWordBackward()
				case 'B':
					state.MoveWORDBackward()
				case 'e':
					state.MoveWordEnd()
				case 'E':
					state.MoveWORDEnd()
				case 'i':
					state.SwitchToInsert("before")
				case 'I':
					state.SwitchToInsert("start")
				case 'a':
					state.SwitchToInsert("after")
				case 'A':
					state.SwitchToInsert("end")
				case 'q', 'Q':
					return ""
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

		ui.Render(state, term.WrappedScreen())
	}
}
