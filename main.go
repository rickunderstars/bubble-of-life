package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CELL            = "󱓻 "
	PERIOD          = 90
	LIFE_PERCENTAGE = 30
)

type model struct {
	grid [][][]bool
	turn bool
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(PERIOD*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) resizeGrid(w, h int) {
	gridW := w / utf8.RuneCountInString(CELL)
	gridH := h

	g := make([][][]bool, 2)
	for i := range g {
		g[i] = make([][]bool, gridH)
		for k := range g[i] {
			g[i][k] = make([]bool, gridW)
		}
	}

	if len(m.grid) == 0 {
		m.grid = g
		return
	}

	for k := 0; k < 2; k++ {
		for i := 0; i < len(m.grid[0]) && i < gridH; i++ {
			for j := 0; j < len(m.grid[0][i]) && j < gridW; j++ {
				g[k][i][j] = m.grid[k][i][j]
			}
		}
	}
	m.grid = g
}

func (m *model) randomGrid() {
	for i := 0; i < len(m.grid[0]); i++ {
		for j := 0; j < len(m.grid[0][i]); j++ {

			if rand.IntN(100) < LIFE_PERCENTAGE {
				m.grid[0][i][j] = true
				m.grid[1][i][j] = true
			} else {
				m.grid[0][i][j] = false
				m.grid[1][i][j] = false
			}

		}
	}
}

func (m model) countNeighbours(x, y int) int {
	n := 0
	h := len(m.grid[0])
	w := len(m.grid[0][0])

	index := 0
	if !m.turn {
		index = 1
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			r := ((y+i)%h + h) % h
			c := ((x+j)%w + w) % w

			if m.grid[index][r][c] {
				n++
			}
		}
	}

	return n
}

func (m *model) evolve() {
}

func initialModel() model {
	m := model{
		grid: make([][][]bool, 0),
		turn: true,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			return m, tea.Quit
		case "r", "R":
			m.randomGrid()
			return m, nil
		}

	case tea.WindowSizeMsg:
		if len(m.grid) == 0 {
			m.resizeGrid(msg.Width, msg.Height)
			m.randomGrid()
		} else {

			m.resizeGrid(msg.Width, msg.Height)
		}
		return m, nil

	case TickMsg:
		return m, doTick()
	}

	return m, nil
}

func (m model) View() string {

	if len(m.grid) < 2 {
		return "Loading..."
	}

	index := 0
	if !m.turn {
		index = 1
	}

	builder := strings.Builder{}

	for i := 0; i < len(m.grid[index]); i++ {
		for j := 0; j < len(m.grid[index][i]); j++ {

			if m.grid[index][i][j] {
				builder.WriteString(CELL)
			} else {
				builder.WriteString("  ")
			}
		}
		if i < len(m.grid[index])-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oh no, an error occurred -> %v", err)
		os.Exit(1)
	}
}
