package app

import (
	"context"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	githubClient := buildGithubClient(ctx, cfg, owner, repoName)
	changelog := buildChangelog(cfg)

	p := tea.NewProgram(view.NewView(ctx, &gitRepo, &githubClient, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}
