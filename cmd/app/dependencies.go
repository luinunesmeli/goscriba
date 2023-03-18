package app

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/pkg/auth"
	"github.com/luinunesmeli/goscriba/tomaster"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

func buildGitRepo(cfg config.Config) (tomaster.GitRepo, error) {
	repo, err := git.PlainOpen(cfg.Path)
	if err != nil {
		return tomaster.GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}

	gitRepo, err := tomaster.NewGitRepo(repo, cfg, auth.AuthMethod(cfg))
	if err != nil {
		return tomaster.GitRepo{}, err
	}

	return gitRepo, nil
}

func buildGithubClient(ctx context.Context, cfg config.Config, owner, repo string) tomaster.GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GetPersonalAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return tomaster.NewGithubClient(github.NewClient(tc), cfg, owner, repo)
}

func buildChangelog(cfg config.Config) tomaster.Changelog {
	return tomaster.NewChangelog(cfg.Changelog)
}
