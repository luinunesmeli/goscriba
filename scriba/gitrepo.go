package scriba

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
)

const (
	developBranchName = "refs/heads/develop"
)

type GitRepo struct {
	repo *git.Repository
}

type Step struct {
	Desc string
	Help string
	Func func() error
}

func NewGitRepo() (GitRepo, error) {
	repo, err := git.PlainOpen("./")
	if err != nil {
		return GitRepo{}, fmt.Errorf("actual directory doesn't contains a git repository: %w", err)
	}

	return GitRepo{
		repo: repo,
	}, nil
}

func (g GitRepo) CheckoutToDevelop() Step {
	return Step{
		Desc: "Checkout to `develop` branch",
		Help: "Looks like `develop` branch don't exist or is unreachable.",
		Func: func() error {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err
			}

			checkoutOpts := &git.CheckoutOptions{
				Branch: developBranchName,
			}
			if err = tree.Checkout(checkoutOpts); err != nil {
				return err
			}

			return nil
		},
	}
}

func (g GitRepo) CheckRepoState() Step {
	return Step{
		Desc: "Checking if current branch is dirty",
		Help: "Commit or stash first your changes before creating a release",
		Func: func() error {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err
			}

			treeStatus, err := tree.Status()
			if err != nil {
				return err
			}

			if !treeStatus.IsClean() {
				return errors.New("current branch has uncommited changes")
			}

			return nil
		},
	}
}

func (g GitRepo) CheckChangelog() Step {
	return Step{
		Desc: "Looking for changelog file",
		Help: "Don't worry, if it doesn't exist I will create for you",
		Func: func() error {
			tree, err := g.repo.Worktree()
			if err != nil {
				return err
			}

			treeStatus, err := tree.Status()
			if err != nil {
				return err
			}

			if treeStatus.IsClean() {
				return errors.New("current branch is dirty")
			}

			return nil
		},
	}
}

func (g GitRepo) GetRepoInfo() {
	c, _ := g.repo.Config()
	fmt.Sprintf(c.User.Name)
}
