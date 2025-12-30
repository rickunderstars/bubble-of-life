package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CELL   = "󱓻 "
	PERIOD = 1000
)

type model struct {
	alive bool
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(PERIOD*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	return model{
		alive: true,
	}
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
	case TickMsg:
		m.alive = !m.alive
		return m, doTick()
	}

	return m, nil
}

func (m model) View() string {
	s := ""
	c := "  "
	if m.alive {
		c = CELL
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			s += c
		}
		s += "\n"
	}

	s += "\n\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oh no, an error occurred: %v", err)
		os.Exit(1)
	}
}
