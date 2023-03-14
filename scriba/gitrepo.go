package scriba

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

type GitRepo struct {
	repo          *git.Repository
	cfg           config.Config
	releaseBranch string
}

const (
	developBranchName = "refs/heads/develop"
	releaseBranchName = "refs/heads/release/%s"
)

func NewGitRepo(cfg config.Config) (GitRepo, error) {
	repo, err := git.PlainOpen(cfg.Path)
	if err != nil {
		return GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}
	return GitRepo{repo: repo, cfg: cfg}, nil
}

func (g *GitRepo) CheckoutToDevelop() Task {
	return Task{
		Desc: "Checkout to develop",
		Help: "Looks like some code wasnt commited at develop.",
		Func: g.CheckoutToBranch(false),
	}
}

func (g *GitRepo) CheckoutToRelease() Task {
	return Task{
		Desc: "Checkout to release branch",
		Help: "Looks like the release branch isn't creates.",
		Func: g.CheckoutToBranch(true),
	}
}

func (g *GitRepo) CheckoutToBranch(asRelease bool) Func {
	return func(session Session) (error, string) {
		tree, err := g.repo.Worktree()
		if err != nil {
			return err, ""
		}

		branch := developBranchName
		if asRelease {
			branch = g.releaseBranch
		}

		checkoutOpts := &git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branch),
			Keep:   true,
		}
		if err = tree.Checkout(checkoutOpts); err != nil {
			return err, ""
		}
		return nil, ""
	}
}

func (g *GitRepo) CheckRepoState() Task {
	return Task{
		Desc: "Checking if current branch is clear",
		Help: "Commit or stash first your changes before creating a release",
		Func: func(session Session) (error, string) {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			status, err := gitStatus(tree)
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
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			opts := git.PullOptions{
				ReferenceName: plumbing.NewBranchReferenceName("develop"),
				SingleBranch:  true,
				Auth:          defaultAuth(g.cfg.GithubTokenAPI),
			}
			if err = tree.Pull(&opts); err != nil {
				if err.Error() == "already up-to-date" {
					return nil, "Already up-to-date! No changes made!"
				}
				return err, ""
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

			g.releaseBranch = fmt.Sprintf(releaseBranchName, session.ChosenVersion)
			ref := plumbing.NewHashReference(plumbing.ReferenceName(g.releaseBranch), headRef.Hash())
			if err = g.repo.Storer.SetReference(ref); err != nil {
				return nil, ""
			}
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
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			opts := &git.CommitOptions{All: true}
			hash, err := tree.Commit(fmt.Sprintf("Automatic release commit %s", session.ChosenVersion), opts)
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
				Auth:       defaultAuth(g.cfg.GithubTokenAPI),
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
	return parts[3], strings.TrimSuffix(parts[4], ".git")
}

func defaultAuth(token string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: "token_user", // yes, this can be anything except an empty string
		Password: token,
	}
}
