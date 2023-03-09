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
	config    Config
	LatestTag string
	ActualPRs PRs
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
	Feature     PRType = "feature"
	Enhancement PRType = "enhancement"
	Fix         PRType = "fix"
	Bugfix      PRType = "bugfix"

	head = "develop"
)

func NewGithubRepo(client *http.Client, cfg Config, owner, repo string) GithubRepo {
	return GithubRepo{
		client: github.NewClient(client),
		owner:  owner,
		repo:   repo,
		config: cfg,
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
		Desc: "Comparing `master` and `develop`",
		Help: "Couldn't get diff!",
		Func: func() (error, string) {
			opts := &github.ListOptions{}
			commits, _, err := r.client.Repositories.CompareCommits(
				ctx, r.owner, r.repo, r.config.Base, head, opts,
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

func (r *GithubRepo) CreatePullRequest(ctx context.Context, tag, body string) Step {
	return Step{
		Desc: "Generating the Pull Request for you.",
		Help: "Couldn't generate the Pull Request!",
		Func: func() (error, string) {
			title := fmt.Sprintf("Release version %s", tag)
			base := r.config.Base
			head := fmt.Sprintf("release/%s", tag)

			newPR := &github.NewPullRequest{
				Title: &title,
				Head:  &head,
				Base:  &base,
				Body:  &body,
			}

			pr, _, err := r.client.PullRequests.Create(ctx, r.owner, r.repo, newPR)
			if err != nil {
				return err, ""
			}

			return nil, fmt.Sprintf("Access at: %s", pr.GetHTMLURL())
		},
	}
}

func getPRType(branch *github.PullRequestBranch) PRType {
	switch {
	case strings.HasPrefix(branch.GetRef(), string(Feature)):
		return Feature
	case strings.HasPrefix(branch.GetRef(), string(Enhancement)):
		return Enhancement
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

func (p PRs) AsMap() map[PRType]PRs {
	out := map[PRType]PRs{}
	for _, pr := range p {
		if _, ok := out[pr.PRType]; !ok {
			out[pr.PRType] = PRs{}
		}
		out[pr.PRType] = append(out[pr.PRType], pr)
	}
	return out
}

func (p PRs) Append(value PR) PRs {
	if len(p) == 0 {
		return PRs{value}
	}

	out := PRs{}
	for _, pr := range p {
		if pr.Number == value.Number {
			return p
		}
		out = append(out, pr)
	}
	return out
}
