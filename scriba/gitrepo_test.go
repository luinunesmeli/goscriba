package scriba_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/scriba"
)

func TestGitRepo_GetRepoInfo(t *testing.T) {
	repo, _ := scriba.NewGitRepo(config.Config{
		Path: "../",
	})

	owner, repoName := repo.GetRepoInfo()

	assert.NotEmpty(t, owner)
	assert.NotEmpty(t, repoName)
}
