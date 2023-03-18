package scriba

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

type GitRepo struct {
	repo          *git.Repository
	tree          *git.Worktree
	cfg           config.Config
	releaseBranch string
	authMethod    transport.AuthMethod
}

const (
	developBranchName = "develop"
	releaseHead       = "refs/heads/release/%s"
	releaseBranch     = "release/%s"
)

func NewGitRepo(repo *git.Repository, cfg config.Config, authMethod transport.AuthMethod) (GitRepo, error) {
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
	}, nil
}

func (g *GitRepo) CheckoutToDevelop() Task {
	return Task{
		Desc: "Checkout to develop",
		Help: "Looks like some code wasnt commited at develop.",
		Func: func(session Session) (error, string) {
			return gitSwitchWrapper(developBranchName, g.tree), ""
		},
	}
}

func (g *GitRepo) CheckoutToRelease() Task {
	return Task{
		Desc: "Checkout to release branch",
		Help: "Looks like the release branch isn't creates.",
		Func: func(session Session) (error, string) {
			return gitSwitchWrapper(g.releaseBranch, g.tree), ""
		},
	}
}

func (g *GitRepo) CheckRepoState() Task {
	return Task{
		Desc: "Checking if current branch is clear",
		Help: "Commit or stash first your changes before creating a release",
		Func: func(session Session) (error, string) {
			status, err := gitStatusWrapper(g.tree)
			if !status.IsClean() {
				return errors.New("current branch has uncommited changes"), ""
			}
			return err, ""
		},
	}
}

func (g *GitRepo) PullDevelop() Task {
	return Task{
		Desc: "Pull changes from remote",
		Help: "Cannot pull changes or there are uncommited changes!",
		Func: func(session Session) (error, string) {
			opts := git.FetchOptions{
				RemoteName: "origin",
				RefSpecs: []gitconfig.RefSpec{
					gitconfig.RefSpec("+refs/heads/develop:refs/remotes/origin/develop"),
				},
				Auth: g.authMethod,
			}

			if err := g.repo.Fetch(&opts); err != nil {
				if err.Error() == "already up-to-date" {
					return nil, "Already up-to-date! No changes made!"
				}
			}
			return nil, ""
		},
	}
}

func (g *GitRepo) CreateRelease() Task {
	return Task{
		Desc: "Create release...",
		Help: "Couldn't create release!",
		Func: func(session Session) (error, string) {
			headRef, err := g.repo.Head()
			if err != nil {
				return nil, ""
			}

			releaseHeadBranch := fmt.Sprintf(releaseHead, session.ChosenVersion)
			ref := plumbing.NewHashReference(plumbing.ReferenceName(releaseHeadBranch), headRef.Hash())
			if err = g.repo.Storer.SetReference(ref); err != nil {
				return nil, ""
			}

			g.releaseBranch = fmt.Sprintf(releaseBranch, session.ChosenVersion)
			return nil, fmt.Sprintf("Created branch release %s", g.releaseBranch)
		},
	}
}

func (g *GitRepo) ReleaseExists(tag string) Task {
	return Task{
		Desc: fmt.Sprintf("Create release/%s", tag),
		Help: "Couldn't create release!",
		Func: func(session Session) (error, string) {
			branchTag := fmt.Sprintf("refs/heads/release/%s", tag)
			branch, err := g.repo.Branch(branchTag)
			if err != nil {
				return nil, ""
			}

			if branch != nil {
				return fmt.Errorf("`%s` already exists", branchTag), ""
			}
			return nil, ""
		},
	}
}

func (g *GitRepo) Commit() Task {
	return Task{
		Desc: "Commit changelog changes...",
		Help: "Some errors found when commiting changes",
		Func: func(session Session) (error, string) {
			opts := &git.CommitOptions{All: false}
			hash, err := g.tree.Commit(fmt.Sprintf("Automatic release commit %s", session.ChosenVersion), opts)
			if err != nil {
				return err, ""
			}

			return nil, fmt.Sprintf("Commited with hash `%s`", hash.String())
		},
	}
}

func (g *GitRepo) PushReleaseBranch() Task {
	return Task{
		Desc: "Push release to remote",
		Help: "Couldn't push release to remote!",
		Func: func(session Session) (error, string) {
			refSpec := fmt.Sprintf(
				"refs/heads/release/%s:refs/heads/release/%s",
				session.ChosenVersion,
				session.ChosenVersion,
			)

			opts := &git.PushOptions{
				RemoteName: "origin",
				RefSpecs:   []gitconfig.RefSpec{gitconfig.RefSpec(refSpec)},
				Auth:       g.authMethod,
			}
			if err := g.repo.Push(opts); err != nil {
				return err, ""
			}

			return nil, fmt.Sprintf("Pushed release/%s", refSpec)
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
