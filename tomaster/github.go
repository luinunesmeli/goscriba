package tomaster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v50/github"

	"github.com/luinunesmeli/goscriba/pkg/config"
	"github.com/luinunesmeli/goscriba/pkg/datapool"
)

type GithubClient struct {
	client   *github.Client
	config   config.Config
	owner    string
	repo     string
	authors  datapool.Pool[string, Author]
	prAuthor Author
}

const (
	head           = "develop"
	initialRelease = "0.0.0"
	closedState    = "closed"
)

func NewGithubClient(client *github.Client, cfg config.Config, owner, repo string) GithubClient {
	return GithubClient{
		client:  client,
		owner:   owner,
		repo:    repo,
		config:  cfg,
		authors: datapool.NewPool[string, Author](),
	}
}

func (r *GithubClient) LoadLatestTag(ctx context.Context) Task {
	return Task{
		Desc: "Loading latest tag",
		Help: "Couldn't get version. Do you have permissions to read this repo?",
		Func: func(session Session) (error, string, Session) {
			rel, resp, err := r.client.Repositories.GetLatestRelease(ctx, r.owner, r.repo)
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					session.LastestVersion = initialRelease
					return nil, "I haven't found any releases, so looks like this is the first release 🥇!", session
				}
				return err, "", session
			}

			session.LastestVersion = rel.GetTagName()
			return nil, fmt.Sprintf("Latest tag is %s!", session.LastestVersion), session
		},
	}
}

func (r *GithubClient) DiffBaseHead(ctx context.Context) Task {
	return Task{
		Desc: "Comparing `master` and `develop`",
		Help: "Couldn't resolve diff!",
		Func: func(session Session) (error, string, Session) {
			commits, _, err := r.client.Repositories.CompareCommits(
				ctx, r.owner, r.repo, r.config.Base, head, &github.ListOptions{},
			)
			if err != nil {
				return err, "", session
			}

			cachedCommits := datapool.NewPool[string, string]()
			cachedPR := datapool.NewPool[int, int]()

			session.PRs = PRs{}
			prOptions := &github.PullRequestListOptions{}
			for _, commit := range commits.Commits {
				if cachedCommits.Has(commit.GetSHA()) {
					continue
				}

				pr, _, _ := r.client.PullRequests.ListPullRequestsWithCommit(ctx, r.owner, r.repo, commit.GetSHA(), prOptions)
				for _, p := range pr {
					if cachedPR.Has(p.GetNumber()) {
						continue
					}
					cachedPR.Add(p.GetNumber(), p.GetNumber())

					if p.GetState() != closedState {
						continue
					}

					commitsPR, _, _ := r.client.PullRequests.ListCommits(
						ctx, r.owner, r.repo, p.GetNumber(), &github.ListOptions{},
					)

					for _, repositoryCommit := range commitsPR {
						cachedCommits.Add(repositoryCommit.GetSHA(), repositoryCommit.GetSHA())
					}

					prType := getPRType(p.GetHead())
					log.Printf("Number: %d Head: %s Type: %s Merged %t", p.GetNumber(), p.GetHead().GetRef(), prType, p.GetMerged())
					if !p.GetMerged() && prType == "" {
						continue
					}

					author, err := r.getAuthor(ctx, p.User.GetLogin())
					if err != nil {
						return err, "", session
					}

					session.PRs = append(session.PRs, PR{
						PRType: prType,
						Title:  p.GetTitle(),
						Link:   p.GetLinks().GetHTML().GetHRef(),
						Author: author,
						Number: p.GetNumber(),
						Ref:    p.GetHead().GetRef(),
					})
				}
			}

			if len(session.PRs) == 0 {
				return errors.New("no closed pull requests on `develop` can be merged on `master`"), "", session
			}

			return nil, "", session
		},
	}
}

func (r *GithubClient) CreatePullRequest(ctx context.Context) Task {
	return Task{
		Desc: "Generating the Pull Request for you.",
		Help: "Couldn't generate the Pull Request!",
		Func: func(session Session) (error, string, Session) {
			title := fmt.Sprintf("Release version %s", session.ChosenVersion)
			head := fmt.Sprintf("release/%s", session.ChosenVersion)

			newPR := &github.NewPullRequest{
				Title: &title,
				Head:  &head,
				Base:  &r.config.Base,
				Body:  &session.Changelog,
			}

			pr, _, err := r.client.PullRequests.Create(ctx, r.owner, r.repo, newPR)
			if err != nil {
				return err, "", Session{}
			}

			session.PRUrl = pr.GetHTMLURL()
			session.PRNumber = pr.GetNumber()

			if _, _, err = r.client.Issues.AddAssignees(ctx, r.owner, r.repo, pr.GetNumber(), []string{r.prAuthor.Login}); err != nil {
				return nil, "could not add assignee to Pull Request", session
			}
			if _, _, err = r.client.Issues.AddLabelsToIssue(ctx, r.owner, r.repo, pr.GetNumber(), []string{r.config.Changelog.ReleaseLabel}); err != nil {
				return nil, "could not add label to Pull Request", session
			}

			return nil, fmt.Sprintf("Access at: %s", pr.GetHTMLURL()), session
		},
	}
}

func (r *GithubClient) GetGithubUsername(ctx context.Context) (Author, error) {
	author, err := r.getAuthor(ctx, "")
	r.prAuthor = author
	return author, err
}

func (r *GithubClient) getAuthor(ctx context.Context, login string) (Author, error) {
	if val, ok := r.authors.Get(login); ok {
		return val, nil
	}

	prUser, _, err := r.client.Users.Get(ctx, login)
	if err != nil {
		return Author{}, err
	}
	author := Author{
		Login: prUser.GetLogin(),
		Name:  prUser.GetName(),
		Email: prUser.GetEmail(),
	}
	r.authors.Add(author.Login, author)

	return author, nil
}
