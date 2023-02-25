package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/scriba"
)

type View struct {
	gitrepo    scriba.GitRepo
	actualStep stepResults
}

func NewView(gitrepo scriba.GitRepo) View {
	return View{
		gitrepo: gitrepo,
	}
}

func (m View) Init() tea.Cmd {
	//m.gitrepo.GetRepoInfo()

	return runSteps(
		//m.gitrepo.CheckRepoState(),
		m.gitrepo.CheckoutToDevelop(),
	)
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case stepResults:
		m.actualStep = msg
		if m.actualStep.checkError() {
			return m, tea.Quit
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m View) View() string {
	if len(m.actualStep) == 0 {
		return ""
	}

	s := ""
	for _, step := range m.actualStep {
		s += fmt.Sprintf("🏃%s... ", step.desc)
		if step.err != nil {
			s += "👎😬\n"
			s += fmt.Sprintf("👹 %s\n", step.err.Error())
			s += fmt.Sprintf("💡 %s\n", step.help)
		} else {
			s += "👍😉\n"
		}
	}
	return s
}
