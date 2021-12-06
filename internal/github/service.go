package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type GithubServicer interface {
	GetContributionsByUsername(ctx context.Context, options GetContributionsByUsernameOptions) (*Contributions, error)
	GetCurrentContributionStreakByUsername(ctx context.Context, username string) (*CurrentContributionStreak, error)
	GetLongestContributionStreakByUsername(ctx context.Context, username string) (*LongestContributionStreak, error)
	GetFirstContributionYearByUsername(ctx context.Context, username string) (*time.Time, error)
}

type GithubService struct {
	githubClient *githubv4.Client
}

var _ GithubServicer = (*GithubService)(nil)

func NewGithubService(githubClient *githubv4.Client) *GithubService {
	return &GithubService{
		githubClient: githubClient,
	}
}
