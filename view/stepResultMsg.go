package view

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/scriba"
)

type stepResultMsg struct {
	desc string
	help string
	err  error
}

type stepResults []stepResultMsg

//type resultMsg func() tea.Msg

func runStep(step scriba.Step) tea.Cmd {
	return func() tea.Msg {
		return stepResultMsg{
			desc: step.Desc,
			help: step.Help,
			err:  step.Func(),
		}
	}
}

func runSteps(steps ...scriba.Step) tea.Cmd {
	return func() tea.Msg {
		out := stepResults{}
		for _, step := range steps {
			err := step.Func()
			out = append(out, stepResultMsg{
				desc: step.Desc,
				help: step.Help,
				err:  err,
			})

			if err != nil {
				return out
			}
		}
		return out
	}
}

func (s stepResults) checkError() bool {
	for _, step := range s {
		if step.err != nil {
			return true
		}
	}
	return false
}
