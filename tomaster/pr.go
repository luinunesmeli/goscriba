package tomaster

import (
	"strings"

	"github.com/google/go-github/v50/github"
)

type (
	PRs []PR

	PRType string

	PR struct {
		PRType PRType
		Title  string
		PRLink string
		Author string
		Number int
		Ref    string
	}
)

const (
	Feature     PRType = "feature"
	Enhancement PRType = "enhancement"
	Fix         PRType = "fix"
	Bugfix      PRType = "bugfix"
)

func getPRType(branch *github.PullRequestBranch) PRType {
	switch {
	case strings.HasPrefix(branch.GetRef(), string(Feature)):
		return Feature
	case strings.HasPrefix(branch.GetRef(), string(Enhancement)):
		return Enhancement
	case strings.HasPrefix(branch.GetRef(), string(Fix)):
		return Fix
	case strings.HasPrefix(branch.GetRef(), string(Bugfix)):
		return Bugfix
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
