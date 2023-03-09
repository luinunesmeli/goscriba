package scriba_test

import (
	"testing"

	"github.com/luinunesmeli/goscriba/scriba"
)

func TestGitRepo_PushRelease(t *testing.T) {
	repo, _ := scriba.NewGitRepo(scriba.Config{
		GithubTokenAPI: "",
		Path:           "../",
	})

	repo.PushReleaseBranch("v0.0.3").Func()
}
