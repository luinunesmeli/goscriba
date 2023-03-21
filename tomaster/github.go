package tomaster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

type GithubClient struct {
	client    *github.Client
	config    config.Config
	owner     string
	repo      string
	LatestTag string
	ActualPRs PRs
}

const (
	head           = "develop"
	initialRelease = "0.0.0"
)

func NewGithubClient(client *github.Client, cfg config.Config, owner, repo string) GithubClient {
	return GithubClient{
		client: client,
		owner:  owner,
		repo:   repo,
		config: cfg,
	}
}

func (r *GithubClient) LoadLatestTag(ctx context.Context) Task {
	return Task{
		Desc: "Loading latest tag",
		Help: "Couldn't get version. Do you have permission to read this repo?",
		Func: func(session Session) (error, string) {
			rel, resp, err := r.client.Repositories.GetLatestRelease(ctx, r.owner, r.repo)
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					r.LatestTag = initialRelease
					return nil, "I haven't found any releases, so looks like this is the first release 🥇!"
				}
				return err, ""
			}

			r.LatestTag = rel.GetTagName()
			return nil, fmt.Sprintf("Latest tag is %s!", r.LatestTag)
		},
	}
}

func (r *GithubClient) DiffBaseHead(ctx context.Context) Task {
	return Task{
		Desc: "Comparing `master` and `develop`",
		Help: "Couldn't get diff!",
		Func: func(session Session) (error, string) {
			commits, _, err := r.client.Repositories.CompareCommits(
				ctx, r.owner, r.repo, r.config.Base, head, &github.ListOptions{},
			)
			if err != nil {
				return err, ""
			}
			prOptions := &github.PullRequestListOptions{State: "open"}
			for _, commit := range commits.Commits {
				pr, _, _ := r.client.PullRequests.ListPullRequestsWithCommit(ctx, r.owner, r.repo, commit.GetSHA(), prOptions)
				for _, p := range pr {
					prType := getPRType(p.GetHead())
					if prType == "" {
						continue
					}

					r.ActualPRs = r.ActualPRs.Append(PR{
						PRType: prType,
						Title:  p.GetTitle(),
						PRLink: p.GetLinks().GetHTML().GetHRef(),
						Author: commit.GetAuthor().GetName(),
						Number: p.GetNumber(),
						Ref:    p.GetHead().GetRef(),
					})
				}
			}
			return nil, ""
		},
	}
}

func (r *GithubClient) CreatePullRequest(ctx context.Context) Task {
	return Task{
		Desc: "Generating the Pull Request for you.",
		Help: "Couldn't generate the Pull Request!",
		Func: func(session Session) (error, string) {
			title := fmt.Sprintf("Release version %s", session.ChosenVersion)
			base := r.config.Base
			head := fmt.Sprintf("release/%s", session.ChosenVersion)
			changelog := session.Changelog

			newPR := &github.NewPullRequest{
				Title: &title,
				Head:  &head,
				Base:  &base,
				Body:  &changelog,
			}
			pr, _, err := r.client.PullRequests.Create(ctx, r.owner, r.repo, newPR)
			if err != nil {
				return err, ""
			}

			return nil, fmt.Sprintf("Access at: %s", pr.GetHTMLURL())
		},
	}
}
