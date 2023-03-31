package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/pkg/auth"
	"github.com/luinunesmeli/goscriba/tomaster"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

func buildGitRepo(cfg config.Config, changelog *tomaster.Changelog) (tomaster.GitRepo, error) {
	url, err := getRemoteURL(cfg)
	log.Printf("Remote is `%s`", url)
	if err != nil {
		return tomaster.GitRepo{}, err
	}

	storer := memory.NewStorage()
	fs := memfs.New()

	fmt.Println("Loading repository...")
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           url,
		Auth:          auth.AuthMethod(cfg),
		Progress:      log.Writer(),
		ReferenceName: plumbing.NewBranchReferenceName("develop"),
		SingleBranch:  true,
		Tags:          git.NoTags,
		NoCheckout:    true,
	})
	if err != nil {
		return tomaster.GitRepo{}, err
	}

	gitRepo, err := tomaster.NewGitRepo(repo, changelog, cfg, auth.AuthMethod(cfg))
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

func getRemoteURL(cfg config.Config) (string, error) {
	plainRepo, err := git.PlainOpen(cfg.Path)
	if err != nil {
		return "", err
	}
	repoCfg, err := plainRepo.Config()
	if err != nil {
		return "", err
	}
	if remotes, ok := repoCfg.Remotes["origin"]; ok && len(remotes.URLs) > 0 {
		return repoCfg.Remotes["origin"].URLs[0], nil
	}
	return "", errors.New("could not load remote URL")
}
