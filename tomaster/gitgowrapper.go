package tomaster

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

// go-git lib has performace problems on big repos with ignored files, specially those with a node_modules/ path
// The lib tries to verify ALL files on some cases, not mattering if they are listed on .gitignore or not
// This is way too slow to be acceptable, peaking up more than 50 seconds on some cases
// So we have this wrapper that executes the git command on shell for some ops like tree state and branch switch
func gitStatusWrapper(wt *git.Worktree) (git.Status, error) {
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

func gitSwitchWrapper(branch string, tree *git.Worktree) error {
	c := exec.Command("git", "switch", branch)
	c.Dir = tree.Filesystem.Root()
	if _, err := c.Output(); err != nil {
		return err
	}

	return nil
}

func gitPulDevelopWrapper() error {
	c := exec.Command("git", "pull", "origin", "develop")
	if _, err := c.Output(); err != nil {
		return err
	}
	return nil
}
