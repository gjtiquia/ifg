package ui

type Cell struct {
	Char  rune
	Style Style
}

type cellKey struct {
	x int
	y int
}

type MockScreen struct {
	width  int
	height int
	cells  map[cellKey]Cell
}

func NewMockScreen(width, height int) *MockScreen {
	return &MockScreen{
		width:  width,
		height: height,
		cells:  make(map[cellKey]Cell),
	}
}

func (m *MockScreen) Clear() {
	m.cells = make(map[cellKey]Cell)
}

func (m *MockScreen) Size() (int, int) {
	return m.width, m.height
}

func (m *MockScreen) SetContent(x, y int, ch rune, style Style) {
	if x >= 0 && x < m.width && y >= 0 && y < m.height {
		m.cells[cellKey{x: x, y: y}] = Cell{Char: ch, Style: style}
	}
}

func (m *MockScreen) Show() {}

func (m *MockScreen) RowAt(y int) string {
	if y < 0 || y >= m.height {
		return ""
	}
	chars := make([]rune, m.width)
	for i := range chars {
		chars[i] = ' '
	}
	for x := 0; x < m.width; x++ {
		if cell, ok := m.cells[cellKey{x: x, y: y}]; ok {
			chars[x] = cell.Char
		}
	}
	return string(chars)
}

func (m *MockScreen) ContentAt(x, y int) (rune, Style) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return ' ', Style{}
	}
	cell, ok := m.cells[cellKey{x: x, y: y}]
	if !ok {
		return ' ', Style{}
	}
	return cell.Char, cell.Style
}

func (m *MockScreen) MaxRow() int {
	maxRow := -1
	for key := range m.cells {
		if key.y > maxRow {
			maxRow = key.y
		}
	}
	return maxRow
}

func (m *MockScreen) HasContentAt(x, y int) bool {
	_, ok := m.cells[cellKey{x: x, y: y}]
	return ok
}
