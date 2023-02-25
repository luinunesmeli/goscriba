package scriba

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GithubRepo struct {
	client *github.Client
}

func NewGithubRepo(cfg Config, ctx context.Context) GithubRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubTokenAPI},
	)
	tc := oauth2.NewClient(ctx, ts)

	return GithubRepo{
		client: github.NewClient(tc),
	}
}

func (r GithubRepo) GetLatestRelease(ctx context.Context) Step {
	return Step{
		Desc: "Looking for latest release version",
		Help: "Couldn't get version. Do you have permission to read this repo?",
		Func: func() error {
			//r.client.Repositories.GetLatestRelease(ctx, )

			return nil
		},
	}
}
