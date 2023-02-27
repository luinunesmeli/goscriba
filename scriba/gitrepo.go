package scriba

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepo struct {
	repo *git.Repository
}

type Step struct {
	Desc string
	Help string
	Func func() (error, string)
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

			treeStatus, err := tree.Status()
			if err != nil {
				return err, ""
			}

			if !treeStatus.IsClean() {
				return errors.New("current branch has uncommited changes"), ""
			}

			return nil, ""
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

			opts := git.PullOptions{}
			err = tree.Pull(&opts)

			if err.Error() == "already up-to-date" {
				return nil, "Already up-to-date! No changes made!"
			}
			return err, ""
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

//func (g GitRepo) CheckChangelog() Step {
//	return Step{
//		Desc: "Looking for changelog file",
//		Help: "Don't worry, if it doesn't exist I will create for you",
//		Func: func() error {
//			tree, err := g.repo.Worktree()
//			if err != nil {
//				return err
//			}
//
//			treeStatus, err := tree.Status()
//			if err != nil {
//				return err
//			}
//
//			if treeStatus.IsClean() {
//				return errors.New("current branch is dirty")
//			}
//
//			return nil
//		},
//	}
//}

func (g GitRepo) GetRepoInfo() {
	c, _ := g.repo.Config()
	fmt.Sprintf(c.User.Name)
}
