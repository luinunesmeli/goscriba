package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/pkg/auth"
	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/scriba"
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

func buildGitRepo(cfg config.Config) (scriba.GitRepo, error) {
	repo, err := git.PlainOpen(cfg.Path)
	if err != nil {
		return scriba.GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}

	gitRepo, err := scriba.NewGitRepo(repo, cfg, auth.AuthMethod(cfg))
	if err != nil {
		return scriba.GitRepo{}, err
	}

	return gitRepo, nil
}

func buildGithubClient(cfg config.Config, owner, repo string) scriba.GithubRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GetPersonalAccessToken()},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	tc.Timeout = time.Second * 5

	return scriba.NewGithubRepo(github.NewClient(tc), cfg, owner, repo)
}

func buildChangelog(cfg config.Config) scriba.Changelog {
	return scriba.NewChangelog(cfg.Changelog)
}
