package view

import (
	"time"

	"github.com/luinunesmeli/goscriba/scriba"
)

type stepResultMsg struct {
	desc    string
	help    string
	err     error
	ok      string
	elapsed float64
}

func runStep(step scriba.Step) stepResultMsg {
	t := time.Now()
	err, msg := step.Func()
	return stepResultMsg{
		desc:    step.Desc,
		help:    step.Help,
		err:     err,
		ok:      msg,
		elapsed: time.Since(t).Seconds(),
	}
}
