package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	state int
)

const (
	boot state = iota + 1
	latestTag
	chooseTag
	updateDevelop
)

func newStateMsg(value state) tea.Cmd {
	return func() tea.Msg {
		return value
	}
}
