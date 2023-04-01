package app

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/view"
)

func Run(cfg config.Config) error {
	log.Printf("Remote is `%s`", cfg.Repo.URL)
	log.Printf("Repository name `%s` with owner `%s`", cfg.Repo.Name, cfg.Repo.Owner)

	changelog := buildChangelog(cfg)

	gitRepo, err := buildGitRepo(cfg, &changelog)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	githubClient := buildGithubClient(ctx, cfg, cfg.Repo.Owner, cfg.Repo.Name)

	p := tea.NewProgram(view.NewView(ctx, &gitRepo, &githubClient, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}
