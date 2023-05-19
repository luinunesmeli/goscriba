package app

import (
	"context"
	"fmt"
	"log"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/pkg/auth"
	"github.com/luinunesmeli/goscriba/tomaster"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

func buildGitRepo(repo *git.Repository, tree *git.Worktree, cfg config.Config) (tomaster.GitRepo, error) {
	gitRepo, err := tomaster.NewGitRepo(repo, tree, cfg, auth.AuthMethod(cfg))
	if err != nil {
		return tomaster.GitRepo{}, err
	}

	return gitRepo, nil
}

func buildGithub(client *github.Client, cfg config.Config, owner, repo string) tomaster.GithubClient {
	return tomaster.NewGithubClient(client, cfg, owner, repo)
}

func githubClient(ctx context.Context, cfg config.Config) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GetPersonalAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func buildChangelog(cfg config.Config, tree *git.Worktree) tomaster.Changelog {
	return tomaster.NewChangelog(cfg, tree)
}

func getWorktree(cfg config.Config, repo *git.Repository) (*git.Worktree, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return tree, err
	}

	ignoreList, err := cfg.ReadGitignore()
	for _, ignore := range ignoreList {
		tree.Excludes = append(tree.Excludes, gitignore.ParsePattern(ignore, []string{}))
	}

	return tree, err
}

func cloneRepository(ctx context.Context, cfg config.Config) (*git.Repository, error) {
	storer := memory.NewStorage()
	fs := memfs.New()

	fmt.Println("Loading repository...")

	return git.CloneContext(ctx, storer, fs, &git.CloneOptions{
		URL:           cfg.Repo.URL,
		Auth:          auth.AuthMethod(cfg),
		ReferenceName: plumbing.NewBranchReferenceName("develop"),
		SingleBranch:  true,
		NoCheckout:    true,
		Progress:      log.Writer(),
		Tags:          git.NoTags,
	})
}
