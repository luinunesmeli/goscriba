package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

type TypeMessage struct {
	Type    string
	Version string
}

func newVersionList(latestTag string) (*selector.Model, error) {
	minor, patch, err := nextTags(latestTag)
	if err != nil {
		return nil, err
	}

	return &selector.Model{
		Data: []interface{}{
			TypeMessage{Type: "Patch", Version: patch},
			TypeMessage{Type: "Minor", Version: minor},
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

func nextTags(latestTag string) (string, string, error) {
	parts := strings.Split(latestTag, ".")

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", err
	}
	major := strings.TrimSuffix(parts[0], "v")

	versionFmt := "%s.%d.%d"
	return fmt.Sprintf(versionFmt, major, minor+1, 0), fmt.Sprintf(versionFmt, major, minor, patch+1), nil
}
