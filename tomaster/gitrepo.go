package tomaster

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

type GitRepo struct {
	repo       *git.Repository
	tree       *git.Worktree
	cfg        config.Config
	authMethod transport.AuthMethod
	changelog  *Changelog
}

const (
	developRef  = "refs/remotes/origin/develop"
	releaseRef  = "refs/heads/release/%s"
	pushRefSpec = "refs/heads/release/%s:refs/heads/release/%s"
)

func NewGitRepo(repo *git.Repository, changelog *Changelog, cfg config.Config, authMethod transport.AuthMethod) (GitRepo, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}
	for _, ignore := range cfg.Gitignore {
		tree.Excludes = append(tree.Excludes, gitignore.ParsePattern(ignore, []string{}))
	}

	return GitRepo{
		repo:       repo,
		tree:       tree,
		cfg:        cfg,
		authMethod: authMethod,
		changelog:  changelog,
	}, nil
}

func (g *GitRepo) CreateRelease() Task {
	return Task{
		Desc: "Create release...",
		Help: "Couldn't create release!",
		Func: func(session Session) (error, string, Session) {
			headRef, err := storer.ResolveReference(g.repo.Storer, developRef)
			if err != nil {
				return nil, "", session
			}

			refName := plumbing.ReferenceName(fmt.Sprintf(releaseRef, session.ChosenVersion))
			ref := plumbing.NewHashReference(refName, headRef.Hash())

			return g.repo.Storer.SetReference(ref), "", session
		},
	}
}

func (g *GitRepo) Commit() Task {
	return Task{
		Desc: "Commit changelog changes...",
		Help: "Some errors found when commiting changes",
		Func: func(session Session) (error, string, Session) {
			repoConfig, err := g.repo.ConfigScoped(gitconfig.SystemScope)
			if err != nil {
				return err, "", session
			}

			localRef := plumbing.NewBranchReferenceName(fmt.Sprintf("release/%s", session.ChosenVersion))
			if err := g.tree.Checkout(&git.CheckoutOptions{Branch: localRef}); err != nil {
				return nil, "", session
			}

			if err, session.Changelog = g.changelog.UpdateChangelog(session, repoConfig.User.Name, g.tree); err != nil {
				return err, "", session
			}

			if _, err := g.tree.Add(g.cfg.Changelog); err != nil {
				return err, "", session
			}

			opts := &git.CommitOptions{}
			hash, err := g.tree.Commit(fmt.Sprintf("Automatic release commit %s", session.ChosenVersion), opts)

			return err, fmt.Sprintf("Commited with hash `%s`", hash), session
		},
	}
}

func (g *GitRepo) PushReleaseBranch() Task {
	return Task{
		Desc: "Push release to remote",
		Help: "Couldn't push release to remote!",
		Func: func(session Session) (error, string, Session) {
			refSpec := fmt.Sprintf(pushRefSpec, session.ChosenVersion, session.ChosenVersion)
			err := g.repo.Push(&git.PushOptions{
				RemoteName: "origin",
				RefSpecs:   []gitconfig.RefSpec{gitconfig.RefSpec(refSpec)},
				Auth:       g.authMethod,
				Force:      true,
			})
			return err, fmt.Sprintf("Pushed release/%s", session.ChosenVersion), session
		},
	}
}

func (g *GitRepo) GetRepoInfo() (string, string) {
	c, _ := g.repo.Config()
	r := c.Remotes

	parts := strings.Split(r["origin"].URLs[0], "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	return owner, strings.TrimSuffix(repo, ".git")
}
