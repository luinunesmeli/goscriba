package tomaster

import (
	"time"
)

type (
	Func func(Session) (error, string)

	Task struct {
		Desc     string
		Help     string
		Func     Func
		Rollback Func
	}

	Result struct {
		Desc    string
		Help    string
		Err     error
		Ok      string
		Elapsed float64
	}

	Manager struct {
		tasks    []Task
		actual   Task
		rollback []Func
	}

	Session struct {
		ChosenVersion string
		PRs           PRs
		Changelog     string
	}
)

func NewManager(task ...Task) Manager {
	return Manager{tasks: task}
}

func (t *Manager) Actual() Task {
	return t.tasks[0]
}

func (t *Manager) RunActual(session Session) Result {
	t.actual, t.tasks = t.tasks[0], t.tasks[1:]

	if t.rollback != nil {
		t.rollback = append(t.rollback, t.actual.Rollback)
	}

	return t.actual.Run(session)
}

func (t *Manager) Empty() bool {
	return len(t.tasks) == 0
}

func (t Task) Run(session Session) Result {
	now := time.Now()
	err, msg := t.Func(session)
	return Result{
		Desc:    t.Desc,
		Help:    t.Help,
		Err:     err,
		Ok:      msg,
		Elapsed: time.Since(now).Seconds(),
	}
}

func NewSession(chosenVersion string, prs PRs, changelog string) Session {
	return Session{
		ChosenVersion: chosenVersion,
		PRs:           prs,
		Changelog:     changelog,
	}
}
