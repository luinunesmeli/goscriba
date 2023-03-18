package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/tomaster"
)

type (
	View struct {
		stepResultList tea.Model
		form           *form
		gitrepo        *tomaster.GitRepo
		github         *tomaster.GithubRepo
		changelog      *tomaster.Changelog
		manager        tomaster.Manager
		config         config.Config
	}
)

func NewView(gitrepo *tomaster.GitRepo, github *tomaster.GithubRepo, changelog *tomaster.Changelog, config config.Config) View {
	v := View{
		gitrepo:        gitrepo,
		github:         github,
		stepResultList: newStepResultList(),
		form:           newForm(),
		changelog:      changelog,
		config:         config,
	}

	ctx := context.Background()
	steps := []tomaster.Task{
		v.changelog.LoadChangelog(),
		v.gitrepo.CheckRepoState(),
		v.gitrepo.CheckoutToDevelop(),
		v.gitrepo.PullDevelop(),
		v.github.LoadLatestTag(ctx),
		v.github.GetPullRequests(ctx),
		v.form.Show(),
		v.gitrepo.CreateRelease(),
		v.gitrepo.CheckoutToRelease(),
		v.changelog.Update(),
		v.gitrepo.Commit(),
	}
	if config.AutoPR {
		steps = append(steps, []tomaster.Task{
			v.gitrepo.PushReleaseBranch(),
			v.github.CreatePullRequest(ctx),
		}...)
	}
	v.manager = tomaster.NewManager(steps...)
	return v
}

func (m View) Init() tea.Cmd {
	return newStateMsg(startStep)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case state:
		switch msg {
		case startStep:
			if m.manager.Empty() {
				return m, tea.Quit
			}
			m.stepResultList, _ = m.stepResultList.Update(startStepMsg{step: m.manager.Actual()})
			return m, newStateMsg(executeStep)
		case executeStep:
			result := m.manager.RunActual(tomaster.NewSession(m.form.chosenTag, m.github.ActualPRs))
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
			if m.manager.Empty() {
				return m, tea.Quit
			}
			return m, newStateMsg(startStep)
		case confirm:
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
