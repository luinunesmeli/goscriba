package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	state int
)

const (
	chooseTag = iota + 1
	createRelease
	confirm
	startStep
	executeStep
	nextStep
)

func newStateMsg(value state) tea.Cmd {
	return func() tea.Msg {
		return value
	}
}
