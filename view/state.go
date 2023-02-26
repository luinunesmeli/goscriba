package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	state int
)

const (
	checkoutRepository state = iota + 1
	fetchLatestTag
	chooseTag
	setDays
	createRelease
)

func newStateMsg(value state) tea.Cmd {
	return func() tea.Msg {
		return value
	}
}
