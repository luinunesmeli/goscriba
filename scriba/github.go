package scriba

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GithubRepo struct {
	client    *github.Client
	LatestTag string
}

func NewGithubRepo(cfg Config, ctx context.Context) GithubRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubTokenAPI},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Timeout = time.Second * 5

	return GithubRepo{
		client: github.NewClient(tc),
	}
}

func (r *GithubRepo) LoadLatestTag(ctx context.Context) Step {
	return Step{
		Desc: "Loading latest tag",
		Help: "Couldn't get version. Do you have permission to read this repo?",
		Func: func() (error, string) {
			rel, _, err := r.client.Repositories.GetLatestRelease(ctx, "luinunesmeli", "goscriba")
			if err != nil {
				return err, ""
			}
			r.LatestTag = rel.GetTagName()
			return nil, fmt.Sprintf("Latest tag is %s!", r.LatestTag)
		},
	}
}

func (r *GithubRepo) GetCommits(ctx context.Context, days int) Step {
	return Step{
		Desc: "Get commits",
		Help: "Couldn't get commits",
		Func: func() (error, string) {
			//sub := time.Now().AddDate(0, 0, days*-1)
			opts := &github.ListOptions{}

			commits, _, err := r.client.Repositories.CompareCommits(
				ctx, "luinunesmeli", "goscriba", "develop", "master", opts,
			)
			if err != nil {
				return err, ""
			}

			fmt.Sprintf(commits.String())

			//opts := &github.PullRequestListOptions{
			//	State: "closed",
			//	Base:  "develop",
			//}
			//
			//commits, _, err := r.client.PullRequests.List()
			//if err != nil {
			//	return err, ""
			//}
			//
			//fmt.Sprintf(commits[0].String())

			return nil, ""
		},
	}
}
