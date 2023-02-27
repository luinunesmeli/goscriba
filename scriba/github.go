package scriba

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v50/github"
)

type GithubRepo struct {
	client    *github.Client
	LatestTag string
	ActualPRs []PR
	owner     string
	repo      string
}

type PRs []PR

type PRType string

type PR struct {
	PRType PRType
	Title  string
	PRLink string
	Author string
	Number int
	Ref    string
}

const (
	PRFeature PRType = "feature"
	base             = "main"
	head             = "develop"
)

func NewGithubRepo(client *http.Client, owner, repo string) GithubRepo {
	return GithubRepo{
		client: github.NewClient(client),
		owner:  owner,
		repo:   repo,
	}
}

func (r *GithubRepo) LoadLatestTag(ctx context.Context) Step {
	return Step{
		Desc: "Loading latest tag",
		Help: "Couldn't get version. Do you have permission to read this repo?",
		Func: func() (error, string) {
			rel, _, err := r.client.Repositories.GetLatestRelease(ctx, r.owner, r.repo)
			if err != nil {
				return err, ""
			}
			r.LatestTag = rel.GetTagName()
			return nil, fmt.Sprintf("Latest tag is %s!", r.LatestTag)
		},
	}
}

func (r *GithubRepo) GetPullRequests(ctx context.Context) Step {
	return Step{
		Desc: "Comparing `main` and `develop`",
		Help: "Couldn't get diff!",
		Func: func() (error, string) {
			opts := &github.ListOptions{}
			commits, _, err := r.client.Repositories.CompareCommits(
				ctx, r.owner, r.repo, base, head, opts,
			)
			if err != nil {
				return err, ""
			}

			opt2 := &github.PullRequestListOptions{
				State: "closed",
			}
			for _, commit := range commits.Commits {
				if commit.SHA == nil {
					continue
				}

				pr, _, _ := r.client.PullRequests.ListPullRequestsWithCommit(
					ctx,
					r.owner,
					r.repo,
					*commit.SHA,
					opt2,
				)

				for _, p := range pr {
					prType := getPRType(p.GetHead())
					if prType == "" {
						continue
					}

					r.ActualPRs = append(r.ActualPRs, PR{
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

func getPRType(branch *github.PullRequestBranch) PRType {
	switch {
	case strings.HasPrefix(branch.GetRef(), string(PRFeature)):
		return PRFeature
	default:
		return ""
	}
}

func (p PRs) Filter(prType PRType) PRs {
	out := PRs{}
	for _, pr := range p {
		if pr.PRType == prType {
			out = append(out, pr)
		}
	}
	return out
}
