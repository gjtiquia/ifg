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
	fmt.Println("ifg - interactive command finder")
	fmt.Println()
	fmt.Println("Usage: ifg [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --sh      Print shell integration code")
	fmt.Println("  --bash    Alias for --sh (backward compatibility)")
	fmt.Println("  --zsh     Alias for --sh (backward compatibility)")
	fmt.Println("  --help    Show this help")
	fmt.Println()
	fmt.Println("Shell Integration:")
	fmt.Println("  Add to ~/.bashrc or ~/.zshrc:")
	fmt.Println("    source \"$(ifg --sh)\"")
	fmt.Println()
	fmt.Println("  Then run 'ifg' to select commands interactively.")
	fmt.Println("  Selected commands are added to history.")
	fmt.Println("  Press UP to access and edit before execution.")
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
	ui.Render(state, term.Screen())

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

		ui.Render(state, term.Screen())
	}
}
