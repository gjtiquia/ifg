package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Key int

const (
	KeyUnknown Key = iota
	KeyUp
	KeyDown
	KeyEnter
	KeyEscape
	KeyBackspace
	KeyCtrlC
	KeyChar
)

type KeyEvent struct {
	Type Key
	Char rune
}

func ReadKey(screen tcell.Screen) KeyEvent {
	ev := screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyUp:
			return KeyEvent{Type: KeyUp}
		case tcell.KeyDown:
			return KeyEvent{Type: KeyDown}
		case tcell.KeyEnter:
			return KeyEvent{Type: KeyEnter}
		case tcell.KeyEscape:
			return KeyEvent{Type: KeyEscape}
		case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyDelete:
			return KeyEvent{Type: KeyBackspace}
		case tcell.KeyCtrlC:
			return KeyEvent{Type: KeyCtrlC}
		default:
			if ev.Rune() != 0 {
				return KeyEvent{Type: KeyChar, Char: ev.Rune()}
			}
		}
	case *tcell.EventResize:
		return KeyEvent{Type: KeyUnknown}
	}

	return KeyEvent{Type: KeyUnknown}
}
