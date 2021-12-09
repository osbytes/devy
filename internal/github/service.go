package github

import (
	"context"
	"time"
)

type GithubClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

type GithubServicer interface {
	GetContributionsByUsername(ctx context.Context, options GetContributionsByUsernameOptions) (*Contributions, error)
	GetCurrentContributionStreakByUsername(ctx context.Context, username string) (*CurrentContributionStreak, error)
	GetLongestContributionStreakByUsername(ctx context.Context, username string) (*LongestContributionStreak, error)
	GetFirstContributionYearByUsername(ctx context.Context, username string) (*time.Time, error)
	GetTotalContributionsByUsername(ctx context.Context, username string) (*TotalContribution, error)
	GetLastRepoByUsername(ctx context.Context, username string) (*LastRepo, error)
	GetLanguagesByUsername(ctx context.Context, username string) (Languages, error)
}

type GithubService struct {
	githubClient GithubClient
}

var _ GithubServicer = (*GithubService)(nil)

func NewGithubService(githubClient GithubClient) *GithubService {
	return &GithubService{
		githubClient: githubClient,
	}
}
