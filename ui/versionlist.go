package ui

import (
	"fmt"

	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/luinunesmeli/goscriba/tomaster"
)

type TypeMessage struct {
	Type    string
	Version string
}

func newVersionList(latestTag string) (*selector.Model, error) {
	major, minor, patch, err := tomaster.NextReleases(latestTag)
	if err != nil {
		return nil, err
	}

	return &selector.Model{
		Data: []interface{}{
			TypeMessage{Type: "Patch", Version: patch},
			TypeMessage{Type: "Minor", Version: minor},
			TypeMessage{Type: "Major", Version: major},
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
	}, nil
}
