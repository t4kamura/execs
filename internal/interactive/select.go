package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	items    []string
	title    string
	cursor   int
	selected string
}

func SelectItem(title string, items []string) (string, error) {
	p := tea.NewProgram(model{items: items, title: title})
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := m.(model); ok && m.selected != "" {
		return m.selected, nil
	}

	return "", fmt.Errorf("invalid choice")
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.selected = m.items[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.items) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.items) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.title + "\n\n")

	for i := 0; i < len(m.items); i++ {
		if m.cursor == i {
			s.WriteString("> ")
		} else {
			s.WriteString("  ")
		}
		s.WriteString(m.items[i])
		s.WriteString("\n")
	}

	return s.String()
}
