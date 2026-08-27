package main

import (
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	CELL              = "󱓻 "
	PERIOD            = 70
	LIFE_PERCENTAGE   = 30
	HISTORY_SIZE      = 2
	CYCLE_RESET_DELAY = 5 * time.Second
)

type model struct {
	grid      [][][]bool
	active    int
	history   []uint64
	resetting bool
}

type TickMsg time.Time

type resetGridMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(PERIOD*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func scheduleReset() tea.Cmd {
	return tea.Tick(CYCLE_RESET_DELAY, func(t time.Time) tea.Msg {
		return resetGridMsg{}
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
	m.history = m.history[:0]
	m.resetting = false
}

func (m model) countNeighbours(x, y int) int {
	n := 0
	h := len(m.grid[0])
	w := len(m.grid[0][0])
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			r := ((y+i)%h + h) % h
			c := ((x+j)%w + w) % w
			if m.grid[m.active][r][c] {
				n++
			}
		}
	}
	return n
}

func (m *model) evolve() {
	world := m.grid[m.active]
	next := (m.active + 1) % 2
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[i]); j++ {
			n := m.countNeighbours(j, i)
			if world[i][j] {
				m.grid[next][i][j] = (n == 2 || n == 3)
			} else {
				m.grid[next][i][j] = (n == 3)
			}
		}
	}
	m.active = next
}

func (m model) hashGrid(idx int) uint64 {
	h := fnv.New64a()
	buf := make([]byte, 0, len(m.grid[idx][0]))
	for _, row := range m.grid[idx] {
		buf = buf[:0]
		for _, alive := range row {
			if alive {
				buf = append(buf, 1)
			} else {
				buf = append(buf, 0)
			}
		}
		h.Write(buf)
	}
	return h.Sum64()
}

func (m *model) checkCycle() bool {
	current := m.hashGrid(m.active)

	for _, past := range m.history {
		if past == current {
			return true
		}
	}

	m.history = append(m.history, current)
	if len(m.history) > HISTORY_SIZE {
		m.history = m.history[1:]
	}
	return false
}

func initialModel() model {
	m := model{
		grid:    make([][][]bool, 0),
		active:  0,
		history: make([]uint64, 0, HISTORY_SIZE),
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
		m.evolve()

		if !m.resetting && m.checkCycle() {
			m.resetting = true
			return m, tea.Batch(doTick(), scheduleReset())
		}

		return m, doTick()
	case resetGridMsg:
		m.randomGrid()
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	if len(m.grid) < 2 {
		return "Loading..."
	}
	builder := strings.Builder{}
	for i := 0; i < len(m.grid[m.active]); i++ {
		for j := 0; j < len(m.grid[m.active][i]); j++ {
			if m.grid[m.active][i][j] {
				builder.WriteString(CELL)
			} else {
				builder.WriteString("  ")
			}
		}
		if i < len(m.grid[m.active])-1 {
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