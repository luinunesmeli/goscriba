package app

import (
	"context"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/scriba"
	"github.com/luinunesmeli/goscriba/view"
)

func Run(cfg config.Config) error {
	gitRepo, err := scriba.NewGitRepo(cfg)
	if err != nil {
		return err
	}
	owner, repo := gitRepo.GetRepoInfo()

	github := scriba.NewGithubRepo(buildOauthclient(cfg), cfg, owner, repo)
	changelog := scriba.NewChangelog(cfg.Changelog)

	p := tea.NewProgram(view.NewView(&gitRepo, &github, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}

func buildOauthclient(cfg config.Config) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GetPersonalAccessToken()},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	tc.Timeout = time.Second * 5

	return tc
}
