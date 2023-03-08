package scriba

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitRepo struct {
	repo *git.Repository
	cfg  Config
}

const (
	developBranchName = "refs/heads/develop"
	releaseBranchName = "refs/heads/release/%s"
)

func NewGitRepo(path string, cfg Config) (GitRepo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}
	return GitRepo{repo: repo, cfg: cfg}, nil
}

func (g GitRepo) CheckoutToDevelop() Step {
	return g.CheckoutToBranch(developBranchName)
}

func (g GitRepo) CheckoutToRelease(tag string) Step {
	return g.CheckoutToBranch(fmt.Sprintf(releaseBranchName, tag))
}

func (g GitRepo) CheckoutToBranch(branch string) Step {
	return Step{
		Desc: fmt.Sprintf("Checkout to `%s` branch", branch),
		Help: fmt.Sprintf("Looks like `%s` branch don't exist or some code wasn't commited.", branch),
		Func: func() (error, string) {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			checkoutOpts := &git.CheckoutOptions{
				Branch: plumbing.ReferenceName(branch),
				Keep:   false,
			}
			if err = tree.Checkout(checkoutOpts); err != nil {
				return err, ""
			}
			return nil, ""
		},
	}
}

func (g GitRepo) CheckRepoState() Step {
	return Step{
		Desc: "Checking if current branch is clear",
		Help: "Commit or stash first your changes before creating a release",
		Func: func() (error, string) {
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

func (g GitRepo) PullDevelop() Step {
	return Step{
		Desc: "Pull changes from remote",
		Help: "Cannot pull changes or there are uncommited changes!",
		Func: func() (error, string) {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			opts := git.PullOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "token_user", // yes, this can be anything except an empty string
					Password: g.cfg.GithubTokenAPI,
				},
			}

			if err = tree.Pull(&opts); err != nil && err.Error() == "already up-to-date" {
				return nil, "Already up-to-date! No changes made!"
			}
			return nil, ""
		},
	}
}

func (g GitRepo) CreateRelease(tag string) Step {
	return Step{
		Desc: fmt.Sprintf("Create release/%s", tag),
		Help: "Couldn't create release!",
		Func: func() (error, string) {
			headRef, err := g.repo.Head()
			if err != nil {
				return nil, ""
			}

			branch := fmt.Sprintf("refs/heads/release/%s", tag)
			ref := plumbing.NewHashReference(plumbing.ReferenceName(branch), headRef.Hash())
			err = g.repo.Storer.SetReference(ref)
			if err != nil {
				return nil, ""
			}

			return nil, "Branch release created"
		},
	}
}

func (g GitRepo) ReleaseExists(tag string) Step {
	return Step{
		Desc: fmt.Sprintf("Create release/%s", tag),
		Help: "Couldn't create release!",
		Func: func() (error, string) {
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

func (g GitRepo) Commit(tag string) Step {
	return Step{
		Desc: "Commit changelog changes...",
		Help: "Some errors found when commiting changes",
		Func: func() (error, string) {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err, ""
			}

			opts := &git.CommitOptions{All: true}
			hash, err := tree.Commit(fmt.Sprintf("Automatic release commit %s", tag), opts)
			if err != nil {
				return err, ""
			}

			return nil, fmt.Sprintf("Commited with hash `%s`", hash.String())
		},
	}
}

func (g GitRepo) PushRelease(tag string) Step {
	return Step{
		Desc: fmt.Sprintf("Push release/%s to remote", tag),
		Help: "Couldn't push release to remote!",
		Func: func() (error, string) {
			opts := &git.PushOptions{
				RemoteName: fmt.Sprintf(releaseBranchName, tag),
			}

			if err := g.repo.Push(opts); err != nil {
				return err, ""
			}

			return nil, ""
		},
	}
}

func (g GitRepo) GetRepoInfo() (string, string) {
	c, _ := g.repo.Config()
	r := c.Remotes

	parts := strings.Split(r["origin"].URLs[0], "/")
	return parts[3], strings.TrimSuffix(parts[4], ".git")
}
