package view

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/luinunesmeli/goscriba/pkg/task"
	"github.com/luinunesmeli/goscriba/scriba"
)

type (
	confirmTag struct {
		latestTag string
	}

	form struct {
		tagSelect           *selector.Model
		confirm             *selector.Model
		show                bool
		latestTag           string
		prs                 scriba.PRs
		chosenTag           string
		yes                 bool
		showTagSelection    bool
		showTagConformation bool
	}
)

func newForm() *form {
	return &form{
		confirm: newConfirm(),
	}
}

func (f *form) Update(msg tea.Msg) (*form, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case string:
		switch msg {
		case common.DONE:
			if f.showTagSelection {
				if i, ok := f.tagSelect.Selected().(TypeMessage); ok {
					f.chosenTag = i.Version
				}

				f.showTagSelection = false
				f.showTagConformation = true

				return f, newConfirmTag()
			}
			if f.showTagConformation {
				if i, ok := f.confirm.Selected().(ConfirmMessage); ok {
					f.yes = i.Yes
				}
				if !f.yes {
					return f, tea.Quit
				}
				f.showTagConformation = false
				f.showTagSelection = false
				f.show = false
				return f, newStateMsg(confirm)
			}
		}
	}

	if f.showTagSelection {
		f.tagSelect, cmd = f.tagSelect.Update(msg)
	}

	if f.showTagConformation {
		f.confirm, cmd = f.confirm.Update(msg)
	}

	return f, cmd
}

func (f *form) View() string {
	if f.showTagSelection {
		return f.tagSelect.View()
	}
	if f.showTagConformation {
		output := f.confirm.View()

		output += fmt.Sprintf("\nRelease Version: %s ", f.chosenTag)
		output += fmt.Sprintf("\nWill contain the following Pull Requests:\n")

		for prType, prs := range f.prs.AsMap() {
			output += fmt.Sprintf("\n%s\n", strings.ToUpper(string(prType)))
			for _, pr := range prs {
				output += fmt.Sprintf(" * [#%d %s] %s by %s\n", pr.Number, pr.Ref, pr.Title, pr.Author)
			}
			output += "\n"
		}

		return output
	}
	return ""
}

func (f *form) Show() task.Task {
	return task.Task{
		Desc: "Select your version",
		Help: "Show form",
		Func: func(_ task.Session) (error, string) {
			f.show = true
			return nil, ""
		},
	}
}

func (f *form) SetLatest(tag string, prs scriba.PRs) error {
	f.latestTag = tag
	f.prs = prs

	list, err := newVersionList(f.latestTag)
	if err != nil {
		return err
	}
	f.showTagSelection = true
	f.tagSelect = list
	return err
}

func newConfirmTag() tea.Cmd {
	return func() tea.Msg {
		return confirmTag{}
	}
}
