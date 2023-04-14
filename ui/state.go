package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	state int
)

const (
	confirm = iota + 1
	startStep
	executeStep
	nextStep
)

func newStateMsg(value state) tea.Cmd {
	return func() tea.Msg {
		return value
	}
}
