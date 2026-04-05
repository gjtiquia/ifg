package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Style struct {
	Bold    bool
	Dim     bool
	Reverse bool
}

type Screen interface {
	Clear()
	Size() (int, int)
	SetContent(x, y int, ch rune, style Style)
	Show()
}

func NewStyle() Style {
	return Style{}
}

func ToTcellStyle(s Style) tcell.Style {
	style := tcell.StyleDefault
	if s.Bold {
		style = style.Bold(true)
	}
	if s.Dim {
		style = style.Dim(true)
	}
	if s.Reverse {
		style = style.Reverse(true)
	}
	return style
}
