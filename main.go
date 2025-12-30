package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CELL   = "󱓻 "
	PERIOD = 90
)

type model struct {
	a [][]bool
	b [][]bool
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

	newA := make([][]bool, gridH)
	newB := make([][]bool, gridH)
	for i := range newA {
		newA[i] = make([]bool, gridW)
	}
	for i := range newB {
		newB[i] = make([]bool, gridW)
	}
	for i := 0; i < len(m.a) && i < gridH; i++ {
		for j := 0; j < len(m.a[i]) && j < gridW; j++ {
			newA[i][j] = m.a[i][j]
			newB[i][j] = m.b[i][j]
		}
	}
	m.a, m.b = newA, newB
}

func (m *model) randomGrid() {
	for i := 0; i < len(m.a) && i < 15; i++ {
		for j := 0; j < len(m.a[i]) && j < 15; j++ {
			m.a[i][j] = true
			m.b[i][j] = true
		}
	}
}

func initialModel() model {
	m := model{
		a: make([][]bool, 0),
		b: make([][]bool, 0),
	}

	m.resizeGrid(20, 20)

	m.randomGrid()

	return m
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.resizeGrid(msg.Width, msg.Height)
		return m, nil

	case TickMsg:
		return m, doTick()
	}

	return m, nil
}

func (m model) View() string {
	builder := strings.Builder{}

	for i := 0; i < len(m.a); i++ {
		for j := 0; j < len(m.a[i]); j++ {

			if m.a[i][j] {
				builder.WriteString(CELL)
			} else {
				builder.WriteString("  ")
			}
		}
		if i < len(m.a)-1 {
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
