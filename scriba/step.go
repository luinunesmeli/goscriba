package scriba

import "time"

type StepFunc func() (error, string)

type Step struct {
	Desc string
	Help string
	Func StepFunc
}

type StepResult struct {
	Desc    string
	Help    string
	Err     error
	Ok      string
	Elapsed float64
}

func RunStep(step Step) StepResult {
	t := time.Now()
	err, msg := step.Func()
	return StepResult{
		Desc:    step.Desc,
		Help:    step.Help,
		Err:     err,
		Ok:      msg,
		Elapsed: time.Since(t).Seconds(),
	}
}
