package scriba

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepo struct {
	repo *git.Repository
}

func NewGitRepo(path string) (GitRepo, error) {
	repo, err := git.PlainOpen(path)
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

func (g GitRepo) GetRepoInfo() (string, string) {
	c, _ := g.repo.Config()
	r := c.Remotes

	parts := strings.Split(r["origin"].URLs[0], "/")
	return parts[3], strings.TrimSuffix(parts[4], ".git")
}

func gitStatus(wt *git.Worktree) (git.Status, error) {
	c := exec.Command("git", "status", "--porcelain", "-z")
	c.Dir = wt.Filesystem.Root()
	output, err := c.Output()
	if err != nil {
		stat, err := wt.Status()
		return stat, err
	}

	lines := strings.Split(string(output), "\000")
	stat := make(map[string]*git.FileStatus, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(strings.TrimLeft(line, " "), " ", 2)
		if len(parts) == 2 {
			stat[strings.Trim(parts[1], " ")] = &git.FileStatus{
				Staging: git.StatusCode([]byte(parts[0])[0]),
			}
		} else {
			stat[strings.Trim(parts[0], " ")] = &git.FileStatus{
				Staging: git.Unmodified,
			}
		}
	}
	return stat, err
}
