package scriba_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/luinunesmeli/goscriba/scriba"
)

func TestGithubRepo_CreatePullRequest(t *testing.T) {
	github := scriba.NewGithubRepo(buildOauthclient(), scriba.Config{}, "luinunesmeli", "goscriba")
	github.CreatePullRequest(context.Background()).Func()
}

func buildOauthclient() *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	tc.Timeout = time.Second * 5

	return tc
}
