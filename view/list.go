package view

import (
	"fmt"

	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

type TypeMessage struct {
	Type    string
	Version string
}

func newVersionList() *selector.Model {
	return &selector.Model{
		Data: []interface{}{
			TypeMessage{Type: "Patch", Version: "0.0.2"},
			TypeMessage{Type: "Minor", Version: "0.1.0"},
		},
		HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select the type of release:"),
		SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Type, t.Version), selector.ColorSelected)
		},
		UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Type, t.Version), selector.ColorUnSelected)
		},
		FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			return ""
		},
		FinishedFunc: func(s interface{}) string {
			return ""
		},
	}
}
