package view

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/pkg/task"
	"github.com/luinunesmeli/goscriba/tomaster"
)

type (
	View struct {
		stepResultList tea.Model
		form           *form
		gitrepo        *tomaster.GitRepo
		github         *tomaster.GithubClient
		changelog      *tomaster.Changelog
		manager        task.Manager
		config         config.Config
		session        task.Session
	}
)

func NewView(ctx context.Context, gitrepo *tomaster.GitRepo, github *tomaster.GithubClient, changelog *tomaster.Changelog, config config.Config) View {
	v := View{
		gitrepo:        gitrepo,
		github:         github,
		stepResultList: newStepResultList(),
		form:           newForm(),
		changelog:      changelog,
		config:         config,
	}

	steps := []task.Task{
		v.changelog.LoadChangelog(),
		v.github.LoadLatestTag(ctx),
		v.github.DiffBaseHead(ctx),
		v.form.Show(),
		v.form.GetSelectedVersion(),
		v.gitrepo.CreateRelease(),
		v.gitrepo.Commit(),
		v.gitrepo.PushReleaseBranch(),
		v.github.CreatePullRequest(ctx),
	}

	v.manager = task.NewManager(steps...)
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
			result := m.manager.RunActual(m.session)
			m.session = result.Session

			m.stepResultList, cmd = m.stepResultList.Update(executeStepMsg{result: result})
			if result.Err != nil {
				return m, tea.Quit
			}

			if m.form.show {
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

	if m.form.show {
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
