package scriba

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

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
