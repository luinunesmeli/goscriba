package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/luinunesmeli/goscriba/tomaster"
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
		prs                 tomaster.PRs
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
				output += fmt.Sprintf(" * [#%d %s] %s by %s\n", pr.Number, pr.Ref, pr.Title, pr.Author.Name)
			}
			output += "\n"
		}

		return output
	}
	return ""
}

func (f *form) Show() tomaster.Task {
	return tomaster.Task{
		Desc: "Select your version",
		Help: "Show form",
		Func: func(ctx context.Context, session tomaster.Session) (error, string, tomaster.Session) {
			f.show = true
			f.latestTag = session.LastestVersion
			f.prs = session.PRs

			list, err := newVersionList(f.latestTag)
			if err != nil {
				return err, "", session
			}
			f.showTagSelection = true
			f.tagSelect = list

			return nil, "", session
		},
	}
}

func (f *form) GetSelectedVersion() tomaster.Task {
	return tomaster.Task{
		Desc: "A new version number was selected!",
		Help: "Version is empty!",
		Func: func(ctx context.Context, session tomaster.Session) (error, string, tomaster.Session) {
			session.ChosenVersion = f.chosenTag
			return nil, "", session
		},
	}
}

func newConfirmTag() tea.Cmd {
	return func() tea.Msg {
		return confirmTag{}
	}
}
