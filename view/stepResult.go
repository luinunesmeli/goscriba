package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/scriba"
)

type (
	stepResultList struct {
		output string
	}

	resultMsg struct {
		state   state
		step    scriba.Step
		result  stepResult
		content string
	}
)

func newStepResultList() stepResultList {
	return stepResultList{}
}

func (s stepResultList) Init() tea.Cmd {
	return nil
}

func (s stepResultList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	resultMsg := msg.(resultMsg)
	var cmd tea.Cmd
	switch resultMsg.state {
	case startStep:
		s.output += fmt.Sprintf("🏃%s... ", resultMsg.step.Desc)
		return s, cmd
	case executeStep:
		if resultMsg.result.err != nil {
			s.output += fmt.Sprintf("👎💩 (took %f)\n", resultMsg.result.elapsed)
			s.output += fmt.Sprintf("👹%s\n", resultMsg.result.err.Error())
			s.output += fmt.Sprintf("💡%s\n", resultMsg.result.help)
			return s, tea.Quit
		}
		s.output += fmt.Sprintf("🤙🤓 (took %f)\n", resultMsg.result.elapsed)

		if resultMsg.result.ok != "" {
			s.output += fmt.Sprintf("💡%s\n", resultMsg.result.ok)
		}
	}
	return s, cmd
}

func (s stepResultList) View() string {
	return s.output
}
