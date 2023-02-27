package view

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/luinunesmeli/goscriba/scriba"
)

type (
	View struct {
		versionList *selector.Model
		confirm     *selector.Model
		gitrepo     scriba.GitRepo
		github      scriba.GithubRepo
		session     Session
		steps       []scriba.Step
	}
)

type Session struct {
	actual        scriba.Step
	state         state
	chosenVersion string
	confirm       bool
	LatestTag     string
	PullRequests  scriba.PRs
	output        string
}

const (
	developBranchName = "refs/heads/develop"
	releaseBranchName = "refs/heads/release/%s"
)

func NewView(gitrepo scriba.GitRepo, github scriba.GithubRepo) View {
	ctx := context.Background()
	return View{
		gitrepo:     gitrepo,
		github:      github,
		versionList: newVersionList(),
		confirm:     newConfirm(),
		steps: []scriba.Step{
			gitrepo.CheckRepoState(),
			gitrepo.CheckoutToBranch(developBranchName),
			gitrepo.PullDevelop(),
			github.LoadLatestTag(ctx),
			github.GetPullRequests(ctx),
		},
	}
}

func (m View) Init() tea.Cmd {
	return newStateMsg(startStep)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		switch msg {
		case common.DONE:
			if m.session.state == chooseTag {
				if i, ok := m.versionList.Selected().(TypeMessage); ok {
					m.session.chosenVersion = i.Version
				}
				return m, newStateMsg(confirm)
			}
			if m.session.state == confirm {
				if i, ok := m.confirm.Selected().(ConfirmMessage); ok {
					m.session.confirm = i.Yes
				}
				if !m.session.confirm {
					return m, tea.Quit
				}

				m.steps = []scriba.Step{
					m.gitrepo.CreateRelease(m.session.chosenVersion),
					m.gitrepo.CheckoutToBranch(fmt.Sprintf(releaseBranchName, m.session.chosenVersion)),
				}
				return m, newStateMsg(startStep)
			}
		}
	case state:
		m.session.state = msg
		switch msg {
		case startStep:
			if len(m.steps) == 0 {
				return m, tea.Quit
			}
			m.session.actual, m.steps = m.steps[0], m.steps[1:]
			m.session.output += fmt.Sprintf("🏃%s... ", m.session.actual.Desc)

			return m, newStateMsg(executeStep)
		case executeStep:
			result := runStep(m.session.actual)
			if result.err != nil {
				m.session.output += fmt.Sprintf("👎💩 (took %f)\n", result.elapsed)
				m.session.output += fmt.Sprintf("👹%s\n", result.err.Error())
				m.session.output += fmt.Sprintf("💡%s\n", result.help)
				return m, tea.Quit
			}
			m.session.output += fmt.Sprintf("🤙🤓 (took %f)\n", result.elapsed)

			if result.ok != "" {
				m.session.output += fmt.Sprintf("💡%s\n", result.ok)
			}

			return m, newStateMsg(nextStep)
		case nextStep:
			if m.github.LatestTag != "" {
				m.session.LatestTag = m.github.LatestTag
				m.session.PullRequests = m.github.ActualPRs

				fmt.Println(m.github.LatestTag)

				return m, newStateMsg(chooseTag)
			}

			if len(m.steps) > 0 {
				return m, newStateMsg(startStep)
			}

			return m, tea.Quit
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	if m.session.state == chooseTag {
		m.versionList, cmd = m.versionList.Update(msg)
	}
	if m.session.state == confirm {
		m.confirm, cmd = m.confirm.Update(msg)
	}
	return m, cmd
}

func (m View) View() string {
	output := m.session.output

	if m.session.state == chooseTag {
		output += m.versionList.View()
	}

	if len(m.session.chosenVersion) > 0 {
		output += fmt.Sprintf("\nRelease Version: %s ", m.session.chosenVersion)
		output += fmt.Sprintf("\nWill contain the following Pull Requests:\n")

		output += fmt.Sprintf("\n%s\n", scriba.PRFeature)
		for _, pr := range m.session.PullRequests.Filter(scriba.PRFeature) {
			output += fmt.Sprintf(" * [#%d %s] %s by %s\n", pr.Number, pr.Ref, pr.Title, pr.Author)
		}
	}

	if m.session.state == confirm {
		output += "\n" + m.confirm.View()
	}

	return output
}
