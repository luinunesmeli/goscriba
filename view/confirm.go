package view

import (
	"fmt"

	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

type ConfirmMessage struct {
	Yes    bool
	Tittle string
}

func newConfirm() *selector.Model {
	return &selector.Model{
		Data: []interface{}{
			ConfirmMessage{Tittle: "Yes", Yes: true},
			ConfirmMessage{Tittle: "No", Yes: false},
		},
		HeaderFunc: selector.DefaultHeaderFuncWithAppend("Create release and generate changelog?"),
		SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(ConfirmMessage)
			return common.FontColor(fmt.Sprintf("%s", t.Tittle), selector.ColorSelected)
		},
		UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(ConfirmMessage)
			return common.FontColor(fmt.Sprintf("%s", t.Tittle), selector.ColorUnSelected)
		},
		FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			return ""
		},
		FinishedFunc: func(s interface{}) string {
			c := s.(ConfirmMessage)
			if !c.Yes {
				return "No changes made!"
			}
			return ""
		},
	}
}
