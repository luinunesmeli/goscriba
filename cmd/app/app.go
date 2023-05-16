package app

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/ui"
)

func Run(cfg config.Config) error {
	log.Printf("Remote is `%s`", cfg.Repo.URL)
	log.Printf("Repository name `%s` with owner `%s`", cfg.Repo.Name, cfg.Repo.Owner)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client := githubClient(ctx, cfg)
	github := buildGithub(client, cfg, cfg.Repo.Owner, cfg.Repo.Name)

	author, err := github.GetGithubUsername(ctx)
	if err != nil {
		return err
	}
	cfg.Repo.Author = config.Author{
		Name:  author.Name,
		Email: author.Email,
		Login: author.Login,
	}

	repo, err := cloneRepository(ctx, cfg)
	if err != nil {
		return err
	}

	tree, err := getWorktree(cfg, repo)
	if err != nil {
		return err
	}

	gitRepo, err := buildGitRepo(repo, tree, cfg)
	if err != nil {
		return err
	}

	changelog := buildChangelog(cfg, tree)

	p := tea.NewProgram(ui.NewView(&gitRepo, &github, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}
