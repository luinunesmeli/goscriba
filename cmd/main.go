package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/scriba"
	"github.com/luinunesmeli/goscriba/view"
)

func main() {
	cfg, err := scriba.LoadConfig()
	if err != nil {
		handleErr(err)
	}

	fmt.Println(cfg)

	gitRepo, err := scriba.NewGitRepo(cfg)
	if err != nil {
		handleErr(err)
	}
	owner, repo := gitRepo.GetRepoInfo()

	github := scriba.NewGithubRepo(buildOauthclient(cfg), cfg, owner, repo)
	changelog := scriba.NewChangelog(cfg.Changelog)

	p := tea.NewProgram(view.NewView(&gitRepo, &github, &changelog, cfg))
	if _, err = p.Run(); err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func buildOauthclient(cfg scriba.Config) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubTokenAPI},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	tc.Timeout = time.Second * 5

	return tc
}
