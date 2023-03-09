package scriba

import (
	"time"
)

type (
	Func func() (error, string)

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
	}

	TaskManager struct {
		tasks  []Task
		actual Task
	}
)

func NewTaskManager(task ...Task) TaskManager {
	return TaskManager{tasks: task}
}

func (t *TaskManager) Actual() Task {
	return t.tasks[0]
}

func (t *TaskManager) RunActual() Result {
	t.actual, t.tasks = t.tasks[0], t.tasks[1:]
	return t.actual.Run()
}

func (t *TaskManager) Empty() bool {
	return len(t.tasks) == 0
}

func (t Task) Run() Result {
	now := time.Now()
	err, msg := t.Func()
	return Result{
		Desc:    t.Desc,
		Help:    t.Help,
		Err:     err,
		Ok:      msg,
		Elapsed: time.Since(now).Seconds(),
	}
}
