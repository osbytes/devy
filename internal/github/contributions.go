package github

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

type Day struct {
	ContributionCount int
	Weekday           int
	Date              time.Time
}
type Contributions struct {
	TotalContributions int
	Days               []Day
}

func (g *GithubService) GetContributionsByUsername(ctx context.Context, username string) (*Contributions, error) {
	var contributionsQuery struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int
					Weeks              []struct {
						ContributionDays []struct {
							ContributionCount int
							Weekday           int
							Date              string
						}
					}
				}
			}
		} `graphql:"user(login: $username)"`
	}

	err := g.githubClient.Query(ctx, &contributionsQuery, map[string]interface{}{
		"username": githubv4.String(username),
	})
	if err != nil {
		return nil, errors.Wrap(err, "github client query")
	}

	contributions := &Contributions{
		TotalContributions: contributionsQuery.User.ContributionsCollection.ContributionCalendar.TotalContributions,
		Days:               []Day{},
	}

	for _, w := range contributionsQuery.User.ContributionsCollection.ContributionCalendar.Weeks {

		for _, d := range w.ContributionDays {

			date, err := time.Parse("2006-01-02", d.Date)
			if err != nil {
				return nil, errors.Wrap(err, "parsing date")
			}

			contributions.Days = append(contributions.Days, Day{
				ContributionCount: d.ContributionCount,
				Weekday:           d.Weekday,
				Date:              date,
			})

		}

	}

	// TODO: figure out how to sort by date DESC in graphql so we don't have to do it here
	sort.Slice(contributions.Days, func(i, j int) bool {
		return contributions.Days[i].Date.After(contributions.Days[j].Date)
	})

	return contributions, nil
}

type CurrentContributionStreak struct {
	Streak    int
	StartedAt time.Time
}

func (c CurrentContributionStreak) String() string {
	return fmt.Sprintf("current contribution streak: %d days started at: %s", c.Streak, c.StartedAt.Format("2006-01-02"))
}

func (g *GithubService) GetCurrentContributionStreakByUsername(ctx context.Context, username string) (*CurrentContributionStreak, error) {

	// TODO: we could decrease bandwidth by making a custom graphql request here that doesn't retrieve some of the unnecessary fields that this retrieves
	contributions, err := g.GetContributionsByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "get contributions by username")
	}

	currentContributionStreak := &CurrentContributionStreak{Streak: 0}

	if len(contributions.Days) == 0 {
		return currentContributionStreak, nil
	}

	// if the current day does not have a contribution, start counting from previous day as the user is still technically on a streak if they contribute today
	days := contributions.Days
	if contributions.Days[0].ContributionCount == 0 {
		days = contributions.Days[1:]
	}

	for _, d := range days {
		if d.ContributionCount == 0 {
			break
		}
		currentContributionStreak.StartedAt = d.Date

		currentContributionStreak.Streak++
	}

	return currentContributionStreak, nil
}

type LongestContributionStreak struct {
	Streak    int
	StartedAt time.Time
	EndedAt   time.Time
}

func (c LongestContributionStreak) String() string {
	return fmt.Sprintf("longest contribution streak: %d days started at: %s ended at: %s", c.Streak, c.StartedAt.Format("2006-01-02"), c.EndedAt.Format("2006-01-02"))
}

func (g *GithubService) GetLongestContributionStreakByUsername(ctx context.Context, username string) (*LongestContributionStreak, error) {

	// TODO: we could decrease bandwidth by making a custom graphql request here that doesn't retrieve some of the unnecessary fields that this retrieves
	contributions, err := g.GetContributionsByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "get contributions by username")
	}

	longestContributionStreak := &LongestContributionStreak{Streak: 0}

	if len(contributions.Days) == 0 {
		return longestContributionStreak, nil
	}

	// if the current day does not have a contribution, start counting from previous day as the user is still technically on a streak if they contribute today
	days := contributions.Days
	if contributions.Days[0].ContributionCount == 0 {
		days = contributions.Days[1:]
	}

	startedAt := time.Time{}
	endedAt := time.Time{}
	streak := 0
	for _, d := range days {
		if d.ContributionCount == 0 {
			if streak > longestContributionStreak.Streak {
				longestContributionStreak.Streak = streak
				longestContributionStreak.StartedAt = startedAt
				longestContributionStreak.EndedAt = endedAt
			}

			streak = 0

			continue
		}

		if streak == 0 {
			endedAt = d.Date
		}
		startedAt = d.Date

		streak++
	}

	return longestContributionStreak, nil
}
