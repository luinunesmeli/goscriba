package view

import (
	"time"

	"github.com/luinunesmeli/goscriba/scriba"
)

type stepResult struct {
	desc    string
	help    string
	err     error
	ok      string
	elapsed float64
}

func runStep(step scriba.Step) stepResult {
	t := time.Now()
	err, msg := step.Func()
	return stepResult{
		desc:    step.Desc,
		help:    step.Help,
		err:     err,
		ok:      msg,
		elapsed: time.Since(t).Seconds(),
	}
}
