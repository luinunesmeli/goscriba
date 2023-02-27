package view

import (
	"context"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/prompt"
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
	daysInput     *prompt.Model
	chosenVersion string
	chosenDays    int
}

const (
	developBranchName = "refs/heads/develop"
	releaseBranchName = "refs/heads/release/%s"
)

func NewView(gitrepo scriba.GitRepo, github scriba.GithubRepo) View {
	return View{
		gitrepo:     gitrepo,
		github:      github,
		versionList: newVersionList(),
		daysInput:   newDaysInput(),
		chosenDays:  -1,
	}
}

func (m View) Init() tea.Cmd {
	return newStateMsg(checkoutRepository)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == common.DONE {
		if m.state == chooseTag {
			i, ok := m.versionList.Selected().(TypeMessage)
			if ok {
				m.chosenVersion = i.Version
			}
			return m, newStateMsg(setDays)
		}
		if m.state == setDays {
			days := m.daysInput.Value()
			m.chosenDays, _ = strconv.Atoi(days)
			return m, newStateMsg(createRelease)
		}
	}

	switch msg := msg.(type) {
	case state:
		m.state = msg
		switch msg {
		case checkoutRepository:
			m.actualStep = m.actualStep.merge(runSteps(
			//m.gitrepo.CheckRepoState(),
			//m.gitrepo.CheckoutToBranch(developBranchName),
			//m.gitrepo.PullDevelop(),
			))

			if m.actualStep.checkError() {
				return m, tea.Quit
			}
			return m, newStateMsg(fetchLatestTag)
		case fetchLatestTag:
			m.actualStep = m.actualStep.merge(runSteps(
				m.github.LoadLatestTag(context.Background()),
			))
			if m.actualStep.checkError() {
				return m, tea.Quit
			}
			m.latestTag = m.github.LatestTag
			return m, newStateMsg(chooseTag)
		case createRelease:
			m.actualStep = m.actualStep.merge(runSteps(
				m.gitrepo.CreateRelease(m.chosenVersion),
				//m.gitrepo.CheckoutToBranch(fmt.Sprintf(releaseBranchName, m.chosenVersion)),
			))
			if m.actualStep.checkError() {
				return m, tea.Quit
			}
			m.latestTag = m.github.LatestTag
			return m, newStateMsg(listCommits)
		case listCommits:
			m.actualStep = m.actualStep.merge(runSteps(
				m.github.GetCommits(context.Background(), m.chosenDays),
			))
			if m.actualStep.checkError() {
				return m, tea.Quit
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
	if m.state == chooseTag {
		m.versionList, cmd = m.versionList.Update(msg)
	}
	if m.state == setDays {
		m.daysInput, cmd = m.daysInput.Update(msg)
	}

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
	if m.state == setDays {
		output += m.daysInput.View()
	}

	if m.chosenVersion != "" {
		output += fmt.Sprintf("\nCreate version: %s | ", m.chosenVersion)
	}
	if m.chosenDays >= 0 {
		output += fmt.Sprintf("Release days: %d\n", m.chosenDays)
	}

	return output
}
