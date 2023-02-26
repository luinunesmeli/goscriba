package view

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"

	"github.com/luinunesmeli/goscriba/scriba"
)

type View struct {
	gitrepo       scriba.GitRepo
	github        scriba.GithubRepo
	actualStep    stepResults
	latestTag     string
	err           error
	state         state
	versionList   *selector.Model
	chosenVersion string
}

func NewView(gitrepo scriba.GitRepo, github scriba.GithubRepo) View {
	return View{
		gitrepo:     gitrepo,
		github:      github,
		versionList: ss(),
	}
}

func (m View) Init() tea.Cmd {
	return newStateMsg(boot)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == common.DONE {
		if m.state == chooseTag {
			i, ok := m.versionList.Selected().(TypeMessage)
			if ok {
				m.chosenVersion = i.Version
			}
			return m, newStateMsg(updateDevelop)
		}
	}

	switch msg := msg.(type) {
	case state:
		m.state = msg
		switch msg {
		case boot:
			m.actualStep = m.actualStep.merge(runSteps(
				m.gitrepo.CheckRepoState(),
				m.gitrepo.CheckoutToDevelop(),
				m.gitrepo.PullDevelop(),
				m.github.LoadLatestTag(context.Background()),
			))

			if m.actualStep.checkError() {
				return m, tea.Quit
			}

			m.latestTag = m.github.LatestTag

			return m, newStateMsg(chooseTag)
			//case latestTag:
			//	m.actualStep = m.actualStep.merge(runSteps(m.github.LoadLatestTag(context.Background())))
			//	m.latestTag = m.github.LatestTag
			//	return m, newStateMsg(chooseTag)
			//case updateDevelop:
			//	m.actualStep = m.actualStep.merge(runSteps(m.gitrepo.PullDevelop()))
			//	return m, nil
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.versionList, cmd = m.versionList.Update(msg)

	return m, cmd
}

func (m View) View() string {
	output := ""
	for _, step := range m.actualStep {
		output += fmt.Sprintf("🏃%s... ", step.desc)
		if step.err != nil {
			output += "👎😬\n"
			output += fmt.Sprintf("👹%s\n", step.err.Error())
			output += fmt.Sprintf("💡%s\n", step.help)
		} else {
			output += "👍😉\n"
			if step.ok != "" {
				output += fmt.Sprintf("💡%s\n", step.ok)
			}
		}
	}

	if m.state == chooseTag {
		output += m.versionList.View()
	}

	if m.chosenVersion != "" {
		output += fmt.Sprintf("\n%s", m.chosenVersion)
	}

	return output
}
