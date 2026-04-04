package ui

import (
	"bufio"
	"os"
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
	Ctrl bool
	Alt  bool
}

func ReadKey() (KeyEvent, error) {
	reader := bufio.NewReader(os.Stdin)

	b, err := reader.ReadByte()
	if err != nil {
		return KeyEvent{}, err
	}

	if b == 3 {
		return KeyEvent{Type: KeyCtrlC}, nil
	}

	if b == 13 || b == 10 {
		return KeyEvent{Type: KeyEnter}, nil
	}

	if b == 27 {
		_, err := reader.Peek(1)
		if err != nil {
			return KeyEvent{Type: KeyEscape}, nil
		}

		next, err := reader.ReadByte()
		if err != nil {
			return KeyEvent{Type: KeyEscape}, nil
		}

		if next == '[' {
			code, err := reader.ReadByte()
			if err != nil {
				return KeyEvent{Type: KeyUnknown}, nil
			}

			switch code {
			case 'A':
				return KeyEvent{Type: KeyUp}, nil
			case 'B':
				return KeyEvent{Type: KeyDown}, nil
			}

			return KeyEvent{Type: KeyUnknown}, nil
		}

		return KeyEvent{Type: KeyEscape}, nil
	}

	if b == 127 || b == 8 {
		return KeyEvent{Type: KeyBackspace}, nil
	}

	if b >= 32 && b <= 126 {
		return KeyEvent{Type: KeyChar, Char: rune(b)}, nil
	}

	return KeyEvent{Type: KeyUnknown}, nil
}
