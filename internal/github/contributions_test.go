package github

import (
	"bot/pkg/date"
	"context"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGithubService_GetContributionsByUsername(t *testing.T) {
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

func TestGithubService_GetContributionsByUsername__MultiYear(t *testing.T) {
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
	assert.Equal(200, resp.TotalContributions)

	githubClient.AssertExpectations(t)
}

func TestGithubService_GetContributionsByUsername__DatesZeroValue(t *testing.T) {
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
			return githubv4.String(options.Username) == params["username"] &&
				date.WithinDuration(time.Now(), params["to"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(time.Now().AddDate(-1, 0, 0), params["from"].(githubv4.DateTime).Time, time.Millisecond)
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

func TestGithubService_GetContributionsByUsername__Errors(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	tests := []struct {
		name    string
		options GetContributionsByUsernameOptions
		wantErr error
	}{
		{
			name:    "test for ErrMissingUsername",
			options: GetContributionsByUsernameOptions{},
			wantErr: ErrMissingUsername,
		},
		{
			name: "test for ErrToDateBeforeFromDate",
			options: GetContributionsByUsernameOptions{
				Username: "test",
				From:     time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC),
				To:       time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrToDateBeforeFromDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := githubService.GetContributionsByUsername(ctx, tt.options)

			assert.Nil(resp)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func TestGithubService_GetFirstContributionYearByUsername(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionYears"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionYears)
		(*a) = contributionYears{
			User: userContributionYears{
				ContributionsCollection: contributionsCollectionContributionYears{
					ContributionYears: []int{
						2019,
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetFirstContributionYearByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), *resp)
	githubClient.AssertExpectations(t)
}

// TODO Tests: CurrentContributionStreak.String()
// labels: tests
func TestCurrentContributionStreak_String(t *testing.T) {

}

// TODO Tests: CurrentContributionStreak.String() no streak
// labels: tests
func TestCurrentContributionStreak_String__NoStreak(t *testing.T) {

}

func TestGithubService_GetCurrentContributionStreakByUsername(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"
	to := time.Now()

	from := time.Date(to.Year()-1, to.Month(), to.Day(), to.Hour(), to.Minute(), to.Second(), to.Nanosecond(), to.Location())

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return date.WithinDuration(from, params["from"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(to, params["to"].(githubv4.DateTime).Time, time.Millisecond) &&
				githubv4.String(username) == params["username"]
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
										Weekday:           7,
										Date:              time.Date(time.Now().Year(), 7, 10, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 5,
										Weekday:           6,
										Date:              time.Date(time.Now().Year(), 7, 9, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 0,
										Weekday:           5,
										Date:              time.Date(time.Now().Year(), 7, 8, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetCurrentContributionStreakByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(time.Date(time.Now().Year(), 7, 9, 0, 0, 0, 0, time.UTC), resp.StartedAt)
	assert.Equal(2, resp.Streak)

}

// TODO Tests: LongestContributionStreak.String()
// labels: tests, good first issue
func TestLongestContributionStreak_String(t *testing.T) {

}

func TestGithubService_GetLongestContributionStreakByUsername(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"
	year := time.Now()

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionYears"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionYears)
		(*a) = contributionYears{
			User: userContributionYears{
				ContributionsCollection: contributionsCollectionContributionYears{
					ContributionYears: []int{
						year.Year(),
					},
				},
			},
		}
	}).Once()

	from := time.Date(year.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"] &&
				date.WithinDuration(from, params["from"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(year, params["to"].(githubv4.DateTime).Time, time.Millisecond)
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
										ContributionCount: 0,
										Weekday:           7,
										Date:              time.Date(time.Now().Year(), 7, 10, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 5,
										Weekday:           7,
										Date:              time.Date(time.Now().Year(), 7, 10, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 5,
										Weekday:           6,
										Date:              time.Date(time.Now().Year(), 7, 9, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 0,
										Weekday:           5,
										Date:              time.Date(time.Now().Year(), 7, 8, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetLongestContributionStreakByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(time.Date(time.Now().Year(), 7, 10, 0, 0, 0, 0, time.UTC), resp.EndedAt)
	assert.Equal(time.Date(time.Now().Year(), 7, 9, 0, 0, 0, 0, time.UTC), resp.StartedAt)
	assert.Equal(2, resp.Streak)

	// TODO Tests: githubClient.AssertExpectations failing
	//panic: interface conversion: interface {} is nil, not githubv4.DateTime [recovered]
	//panic: interface conversion: interface {} is nil, not githubv4.DateTime
	// githubClient.AssertExpectations(t)
}

func TestGithubService_GetLongestContributionStreakByUsername__NoEndDateCurrentStreak(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"
	year := time.Now()

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionYears"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionYears)
		(*a) = contributionYears{
			User: userContributionYears{
				ContributionsCollection: contributionsCollectionContributionYears{
					ContributionYears: []int{
						year.Year(),
					},
				},
			},
		}
	}).Once()

	from := time.Date(year.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"] &&
				date.WithinDuration(from, params["from"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(year, params["to"].(githubv4.DateTime).Time, time.Millisecond)
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
										Weekday:           7,
										Date:              time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 5,
										Weekday:           6,
										Date:              time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
									{
										ContributionCount: 0,
										Weekday:           5,
										Date:              time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-2, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
									},
								},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetLongestContributionStreakByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(time.Time{}, resp.EndedAt)
	assert.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.UTC), resp.StartedAt)
	assert.Equal(2, resp.Streak)

}
func TestGithubService_GetLongestContributionStreakByUsername__NoContributionDays(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"
	year := time.Now()

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionYears"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionYears)
		(*a) = contributionYears{
			User: userContributionYears{
				ContributionsCollection: contributionsCollectionContributionYears{
					ContributionYears: []int{
						year.Year(),
					},
				},
			},
		}
	}).Once()

	from := time.Date(year.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"] &&
				date.WithinDuration(from, params["from"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(year, params["to"].(githubv4.DateTime).Time, time.Millisecond)
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
								[]contributionDays{},
							},
						},
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetLongestContributionStreakByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(time.Time{}, resp.EndedAt)
	assert.Equal(time.Time{}, resp.StartedAt)
	assert.Equal(0, resp.Streak)
}

func TestGithubService_GetTotalContributionsByUsername(t *testing.T) {
	assert := assert.New(t)
	githubClient := &MockGithubClient{}
	githubService := NewGithubService(githubClient)

	ctx := context.Background()

	username := "devy"
	to := time.Now()

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionYears"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"]
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionYears)
		(*a) = contributionYears{
			User: userContributionYears{
				ContributionsCollection: contributionsCollectionContributionYears{
					ContributionYears: []int{
						to.Year(),
					},
				},
			},
		}
	}).Once()

	year := time.Date(to.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func(params map[string]interface{}) bool {
			return githubv4.String(username) == params["username"] &&
				date.WithinDuration(year, params["from"].(githubv4.DateTime).Time, time.Millisecond) &&
				date.WithinDuration(time.Now(), params["to"].(githubv4.DateTime).Time, time.Millisecond)
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{
			User: user{
				ContributionsCollection: contributionsCollection{
					ContributionCalendar: contributionCalendar{
						TotalContributions: 1000,
					},
				},
			},
		}
	}).Once()

	resp, err := githubService.GetTotalContributionsByUsername(ctx, username)

	assert.NoError(err)
	assert.Equal(resp.String(), "total github contributions: 1000")
	assert.Equal(resp.Total, 1000)
}
