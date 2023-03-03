package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/scriba"
)

type (
	View struct {
		stepResultList tea.Model
		form           *form
		gitrepo        *scriba.GitRepo
		github         *scriba.GithubRepo
		changelog      *scriba.Changelog
		session        Session
		steps          []scriba.Step
	}
)

type Session struct {
	actual scriba.Step
	state  state
}

func NewView(gitrepo *scriba.GitRepo, github *scriba.GithubRepo, changelog *scriba.Changelog) View {
	f := newForm()
	ctx := context.Background()
	return View{
		gitrepo:        gitrepo,
		github:         github,
		stepResultList: newStepResultList(),
		form:           f,
		changelog:      changelog,
		steps: []scriba.Step{
			//changelog.LoadChangelog(),
			gitrepo.CheckRepoState(),
			gitrepo.CheckoutToDevelop(),
			gitrepo.PullDevelop(),
			github.LoadLatestTag(ctx),
			github.GetPullRequests(ctx),
			f.Show(),
		},
	}
}

func (m View) Init() tea.Cmd {
	return newStateMsg(startStep)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case state:
		m.session.state = msg
		switch msg {
		case startStep:
			if len(m.steps) == 0 {
				return m, tea.Quit
			}
			m.session.actual, m.steps = m.steps[0], m.steps[1:]
			m.stepResultList, _ = m.stepResultList.Update(startStepMsg{step: m.session.actual})

			return m, newStateMsg(executeStep)
		case executeStep:
			result := scriba.RunStep(m.session.actual)
			m.stepResultList, cmd = m.stepResultList.Update(executeStepMsg{result: result})
			if result.Err != nil {
				return m, tea.Quit
			}

			if m.form.show && m.github.LatestTag != "" {
				if err := m.form.SetLatest(m.github.LatestTag, m.github.ActualPRs); err != nil {
					return m, tea.Quit
				}
				m.form, cmd = m.form.Update(msg)
				return m, cmd
			}

			return m, newStateMsg(nextStep)
		case nextStep:
			if len(m.steps) > 0 {
				return m, newStateMsg(startStep)
			}
			return m, tea.Quit
		case confirm:
			m.changelog.PRs = m.github.ActualPRs
			m.steps = []scriba.Step{
				m.gitrepo.CreateRelease(m.form.chosenTag),
				m.gitrepo.CheckoutToRelease(m.form.chosenTag),
				m.changelog.Update(m.form.chosenTag),
				m.gitrepo.Commit(m.form.chosenTag),
				m.gitrepo.PushRelease(m.form.chosenTag),
			}
			return m, newStateMsg(startStep)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	if m.form.show && m.github.LatestTag != "" {
		m.form, cmd = m.form.Update(msg)
	}

	return m, cmd
}

func (m View) View() string {
	output := m.stepResultList.View()

	if m.form.show {
		output += m.form.View()
	}

	return output
}
