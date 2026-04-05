package ui

import (
	"github.com/gdamore/tcell/v2"
)

type TcellScreen struct {
	screen tcell.Screen
}

func NewTcellScreen(ts tcell.Screen) *TcellScreen {
	return &TcellScreen{screen: ts}
}

func (t *TcellScreen) Clear() {
	t.screen.Clear()
}

func (t *TcellScreen) Size() (int, int) {
	return t.screen.Size()
}

func (t *TcellScreen) SetContent(x, y int, ch rune, style Style) {
	t.screen.SetContent(x, y, ch, nil, ToTcellStyle(style))
}

func (t *TcellScreen) Show() {
	t.screen.Show()
}
