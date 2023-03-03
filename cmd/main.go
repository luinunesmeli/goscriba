package main

import (
	"context"
	"flag"
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

	basePath := cliParams()

	//gitRepo, err := scriba.NewGitRepo("./")
	//gitRepo, err := scriba.NewGitRepo("/Users/luinunes/nodejs/fury_shp-lhts-rostering-frontend", cfg)

	gitRepo, err := scriba.NewGitRepo(basePath, cfg)
	if err != nil {
		handleErr(err)
	}
	owner, repo := gitRepo.GetRepoInfo()

	github := scriba.NewGithubRepo(buildOauthclient(cfg), owner, repo)
	changelog := scriba.NewChangelog(basePath + "docs/guide/pages/changelog.md")

	p := tea.NewProgram(view.NewView(&gitRepo, &github, &changelog))
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

func cliParams() string {
	path := flag.String("path", "./", "project path you want to generate a release")
	if path == nil {
		return ""
	}
	return *path
}
