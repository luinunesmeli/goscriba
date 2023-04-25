package tomaster

import (
	"time"
)

type (
	Func func(Session) (error, string, Session)

	Task struct {
		Desc string
		Help string
		Func Func
	}

	Result struct {
		Desc    string
		Help    string
		Err     error
		Ok      string
		Elapsed float64
		Session Session
	}

	Manager struct {
		tasks    []Task
		actual   Task
		rollback []Func
	}

	Session struct {
		ChosenVersion  string
		LastestVersion string
		PRs            PRs
		Changelog      string
		PRUrl          string
		PRNumber       int
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
	return t.actual.Run(session)
}

func (t *Manager) Empty() bool {
	return len(t.tasks) == 0
}

func (t Task) Run(session Session) Result {
	now := time.Now()
	err, msg, session := t.Func(session)
	return Result{
		Desc:    t.Desc,
		Help:    t.Help,
		Err:     err,
		Ok:      msg,
		Elapsed: time.Since(now).Seconds(),
		Session: session,
	}
}
