package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/view"
)

func Run(cfg config.Config) error {
	gitRepo, err := buildGitRepo(cfg)
	if err != nil {
		return err
	}
	owner, repoName := gitRepo.GetRepoInfo()

	githubClient := buildGithubClient(cfg, owner, repoName)
	changelog := buildChangelog(cfg)

	p := tea.NewProgram(view.NewView(&gitRepo, &githubClient, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}
