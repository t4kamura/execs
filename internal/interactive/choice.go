package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices     []string
	questionMsg string
	cursor      int
	selected    string
}

func AskChoices(questionMsg string, choices []string) (string, error) {
	p := tea.NewProgram(model{choices: choices, questionMsg: questionMsg})
	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	// Assert the final tea.Model to our local model and print the choice.
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
			// Send the choice on the channel and exit.
			m.selected = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.questionMsg + "\n\n")

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
