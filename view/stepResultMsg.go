package view

import (
	"github.com/luinunesmeli/goscriba/scriba"
)

type stepResultMsg struct {
	desc string
	help string
	err  error
	ok   string
}

type stepResults []stepResultMsg

func runSteps(steps ...scriba.Step) stepResults {
	out := stepResults{}
	for _, step := range steps {
		err, msg := step.Func()
		out = append(out, stepResultMsg{
			desc: step.Desc,
			help: step.Help,
			err:  err,
			ok:   msg,
		})

		if err != nil {
			return out
		}
	}
	return out
}

func (s stepResults) checkError() bool {
	for _, step := range s {
		if step.err != nil {
			return true
		}
	}
	return false
}

func (s stepResults) merge(res stepResults) stepResults {
	merged := s
	for _, msg := range res {
		merged = append(merged, msg)
	}
	return merged
}
