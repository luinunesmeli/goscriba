package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/tomaster"
)

type (
	stepResultList struct {
		output string
	}

	executeStepMsg struct {
		result tomaster.Result
	}

	startStepMsg struct {
		step tomaster.Task
	}
)

func newStepResultList() stepResultList {
	return stepResultList{}
}

func (s stepResultList) Init() tea.Cmd {
	return nil
}

func (s stepResultList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case startStepMsg:
		s.output += fmt.Sprintf("🏃%s... ", msg.step.Desc)
		return s, cmd
	case executeStepMsg:
		if msg.result.Err != nil {
			s.output += fmt.Sprintf("👎💩 (took %f)\n", msg.result.Elapsed)
			s.output += fmt.Sprintf("👹Error: %s\n", msg.result.Err.Error())
			s.output += fmt.Sprintf("💡%s\n", msg.result.Help)
			return s, tea.Quit
		}
		s.output += fmt.Sprintf("🤙🤓 (took %f)\n", msg.result.Elapsed)

		if msg.result.Ok != "" {
			s.output += fmt.Sprintf("💡%s\n", msg.result.Ok)
		}
	}
	return tea.Model(s), cmd
}

func (s stepResultList) View() string {
	return s.output
}
