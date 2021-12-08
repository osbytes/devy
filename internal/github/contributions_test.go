package github

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetContributionsByUsername(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()
	from := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	options := GetContributionsByUsernameOptions{
		Username: "devy",
		From:     from,
		To:       to,
	}

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			assert.Equal(githubv4.String(options.Username), params["username"])
			assert.Equal(githubv4.DateTime{Time: from}, params["from"])
			assert.Equal(githubv4.DateTime{Time: to}, params["to"])
			return true
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{
			User: user{
				ContributionsCollection: contributionsCollection{
					ContributionCalendar: contributionCalendar{
						TotalContributions: 100,
						Weeks: []week{
							{
								[]contributionDays{
									{
										ContributionCount: 5,
										Weekday:           1,
										Date:              "2019-01-01",
									},
								},
							},
							{
								[]contributionDays{
									{
										ContributionCount: 10,
										Weekday:           0,
										Date:              "2019-01-02",
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetContributionsByUsername(ctx, options)

	assert.NoError(err)

	assert.Equal(10, resp.Days[0].ContributionCount)
	assert.Equal(0, resp.Days[0].Weekday)
	assert.Equal(time.Date(2019, 01, 02, 0, 0, 0, 0, time.UTC), resp.Days[0].Date)

	assert.Equal(5, resp.Days[1].ContributionCount)
	assert.Equal(1, resp.Days[1].Weekday)
	assert.Equal(time.Date(2019, 01, 01, 0, 0, 0, 0, time.UTC), resp.Days[1].Date)

	githubClient.AssertExpectations(t)
}

func TestGetContributionsByUsername_MultiYear(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()
	from := time.Date(2020, 12, 6, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, 12, 7, 0, 0, 0, 0, time.UTC)

	options := GetContributionsByUsernameOptions{
		Username: "devy",
		From:     from,
		To:       to,
	}

	fmt.Println("check here", from)

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(options.Username) == params["username"] &&
				githubv4.DateTime{Time: time.Date(2020, 12, 7, 0, 0, 0, 0, time.UTC)} == params["from"] &&
				githubv4.DateTime{Time: time.Date(2021, 12, 7, 0, 0, 0, 0, time.UTC)} == params["to"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{
			User: user{
				ContributionsCollection: contributionsCollection{
					ContributionCalendar: contributionCalendar{
						TotalContributions: 100,
						Weeks: []week{
							{
								[]contributionDays{
									{
										ContributionCount: 5,
										Weekday:           1,
										Date:              "2019-01-01",
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(options.Username) == params["username"] &&
				githubv4.DateTime{Time: time.Date(2020, 12, 6, 0, 0, 0, 0, time.UTC)} == params["from"] &&
				githubv4.DateTime{Time: time.Date(2020, 12, 6, 0, 0, 0, 0, time.UTC)} == params["to"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{
			User: user{
				ContributionsCollection: contributionsCollection{
					ContributionCalendar: contributionCalendar{
						TotalContributions: 100,
						Weeks: []week{
							{
								[]contributionDays{
									{
										ContributionCount: 5,
										Weekday:           1,
										Date:              "2019-01-01",
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetContributionsByUsername(ctx, options)

	assert.NoError(err)

	assert.Equal(5, resp.Days[0].ContributionCount)
	assert.Equal(1, resp.Days[0].Weekday)
	assert.Equal(time.Date(2019, 01, 01, 0, 0, 0, 0, time.UTC), resp.Days[0].Date)

	githubClient.AssertExpectations(t)
}

func TestGetContributionsByUsername_DatesZeroValue(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	options := GetContributionsByUsernameOptions{
		Username: "devy",
	}

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(options.Username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{
			User: user{
				ContributionsCollection: contributionsCollection{
					ContributionCalendar: contributionCalendar{
						TotalContributions: 100,
						Weeks: []week{
							{
								[]contributionDays{
									{
										ContributionCount: 5,
										Weekday:           1,
										Date:              "2019-01-01",
									},
								},
							},
							{
								[]contributionDays{
									{
										ContributionCount: 10,
										Weekday:           0,
										Date:              "2019-01-02",
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetContributionsByUsername(ctx, options)

	assert.NoError(err)

	assert.Equal(10, resp.Days[0].ContributionCount)
	assert.Equal(0, resp.Days[0].Weekday)
	assert.Equal(time.Date(2019, 01, 02, 0, 0, 0, 0, time.UTC), resp.Days[0].Date)

	assert.Equal(5, resp.Days[1].ContributionCount)
	assert.Equal(1, resp.Days[1].Weekday)
	assert.Equal(time.Date(2019, 01, 01, 0, 0, 0, 0, time.UTC), resp.Days[1].Date)

	githubClient.AssertExpectations(t)
}

// TODO Tests: error table test on GetContributionsByUsername
// labels: tests
// Need to run a table test on GetContributionsByUsername to hit
// ErrMissingUsername and ErrToDateBeforeFromDate
func TestGetContributionsByUsername_Errors(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(3, 3)
}

// TODO Tests: GetFirstContributionYearByUsername
// labels: tests
func TestGetFirstContributionYearByUsername(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: CurrentContributionStreak.String()
// labels: tests
func TestCurrentContributionStreak_String(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: GetCurrentContributionStreakByUsername
// labels: tests
func TestGetCurrentContributionStreakByUsername(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: LongestContributionStreak.String()
// labels: tests
func TestLongestContributionStreak_String(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: GetLongestContributionStreakByUsername
// labels: tests
func TestGetLongestContributionStreakByUsername(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: TotalContribution.String()
// labels: tests
func TestTotalContribution_String(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}

// TODO Tests: GetTotalContributionsByUsername
// labels: tests
func TestGetTotalContributionsByUsername(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}
