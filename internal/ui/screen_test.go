package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewStyle(t *testing.T) {
	style := NewStyle()

	if style.Bold {
		t.Error("expected Bold to be false")
	}
	if style.Dim {
		t.Error("expected Dim to be false")
	}
	if style.Reverse {
		t.Error("expected Reverse to be false")
	}
}

func TestToTcellStyle(t *testing.T) {
	getAttrs := func(s tcell.Style) tcell.AttrMask {
		_, _, attr := s.Decompose()
		return attr
	}

	t.Run("empty style returns StyleDefault", func(t *testing.T) {
		style := Style{}
		tcellStyle := ToTcellStyle(style)

		if tcellStyle != tcell.StyleDefault {
			t.Error("expected StyleDefault for empty style")
		}
	})

	t.Run("bold only", func(t *testing.T) {
		style := Style{Bold: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrBold == 0 {
			t.Error("expected Bold to be set")
		}
		if attr&tcell.AttrDim != 0 {
			t.Error("expected Dim to be unset")
		}
		if attr&tcell.AttrReverse != 0 {
			t.Error("expected Reverse to be unset")
		}
	})

	t.Run("dim only", func(t *testing.T) {
		style := Style{Dim: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrDim == 0 {
			t.Error("expected Dim to be set")
		}
		if attr&tcell.AttrBold != 0 {
			t.Error("expected Bold to be unset")
		}
		if attr&tcell.AttrReverse != 0 {
			t.Error("expected Reverse to be unset")
		}
	})

	t.Run("reverse only", func(t *testing.T) {
		style := Style{Reverse: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrReverse == 0 {
			t.Error("expected Reverse to be set")
		}
		if attr&tcell.AttrBold != 0 {
			t.Error("expected Bold to be unset")
		}
		if attr&tcell.AttrDim != 0 {
			t.Error("expected Dim to be unset")
		}
	})

	t.Run("bold and dim", func(t *testing.T) {
		style := Style{Bold: true, Dim: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrBold == 0 {
			t.Error("expected Bold to be set")
		}
		if attr&tcell.AttrDim == 0 {
			t.Error("expected Dim to be set")
		}
		if attr&tcell.AttrReverse != 0 {
			t.Error("expected Reverse to be unset")
		}
	})

	t.Run("bold and reverse", func(t *testing.T) {
		style := Style{Bold: true, Reverse: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrBold == 0 {
			t.Error("expected Bold to be set")
		}
		if attr&tcell.AttrReverse == 0 {
			t.Error("expected Reverse to be set")
		}
		if attr&tcell.AttrDim != 0 {
			t.Error("expected Dim to be unset")
		}
	})

	t.Run("dim and reverse", func(t *testing.T) {
		style := Style{Dim: true, Reverse: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrDim == 0 {
			t.Error("expected Dim to be set")
		}
		if attr&tcell.AttrReverse == 0 {
			t.Error("expected Reverse to be set")
		}
		if attr&tcell.AttrBold != 0 {
			t.Error("expected Bold to be unset")
		}
	})

	t.Run("all flags set", func(t *testing.T) {
		style := Style{Bold: true, Dim: true, Reverse: true}
		tcellStyle := ToTcellStyle(style)
		attr := getAttrs(tcellStyle)

		if attr&tcell.AttrBold == 0 {
			t.Error("expected Bold to be set")
		}
		if attr&tcell.AttrDim == 0 {
			t.Error("expected Dim to be set")
		}
		if attr&tcell.AttrReverse == 0 {
			t.Error("expected Reverse to be set")
		}
	})
}
