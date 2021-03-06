package github

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

var (
	ErrMissingUsername      = errors.New("missing github username")
	ErrToDateBeforeFromDate = errors.New("to date is before from date")
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

type GetContributionsByUsernameOptions struct {
	Username string
	From     time.Time
	To       time.Time
}

type contributionsQuery struct {
	User user `graphql:"user(login: $username)"`
}

type user struct {
	ContributionsCollection contributionsCollection `graphql:"contributionsCollection(from: $from, to: $to)"`
}

type contributionsCollection struct {
	ContributionCalendar contributionCalendar
}

type contributionCalendar struct {
	TotalContributions int
	Weeks              []week
}

type week struct {
	ContributionDays []contributionDays
}

type contributionDays struct {
	ContributionCount int
	Weekday           int
	Date              string
}

func (g *GithubService) GetContributionsByUsername(ctx context.Context, options GetContributionsByUsernameOptions) (*Contributions, error) {
	if len(options.Username) == 0 {
		return nil, ErrMissingUsername
	}

	from, to := options.From, options.To

	if to.IsZero() {
		to = time.Now()
	}

	if from.IsZero() {
		from = to.AddDate(-1, 0, 0)
	}

	if to.Before(from) {
		return nil, ErrToDateBeforeFromDate
	}

	contributions := &Contributions{
		TotalContributions: 0,
		Days:               []Day{},
	}

	originalFrom := from

	for {

		from = to.AddDate(-1, 0, 0)
		if from.Before(originalFrom) {
			from = originalFrom
		}

		var cq contributionsQuery

		err := g.githubClient.Query(ctx, &cq, map[string]interface{}{
			"username": githubv4.String(options.Username),
			"from":     githubv4.DateTime{Time: from},
			"to":       githubv4.DateTime{Time: to},
		})
		if err != nil {
			return nil, errors.Wrap(err, "github client query")
		}

		contributions.TotalContributions += cq.User.ContributionsCollection.ContributionCalendar.TotalContributions

		for _, w := range cq.User.ContributionsCollection.ContributionCalendar.Weeks {

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

		// offset to date by one day because github's graphql api returns inclusive days between the from and to
		to = from.AddDate(0, 0, -1)

		if from.Equal(originalFrom) {
			break
		}
	}

	// TODO: figure out how to sort by date DESC in graphql so we don't have to do it here
	sort.Slice(contributions.Days, func(i, j int) bool {
		return contributions.Days[i].Date.After(contributions.Days[j].Date)
	})

	return contributions, nil
}

type contributionYears struct {
	User userContributionYears `graphql:"user(login: $username)"`
}

type userContributionYears struct {
	ContributionsCollection contributionsCollectionContributionYears
}

type contributionsCollectionContributionYears struct {
	ContributionYears []int
}

func (g *GithubService) GetFirstContributionYearByUsername(ctx context.Context, username string) (*time.Time, error) {
	var cy contributionYears

	err := g.githubClient.Query(ctx, &cy, map[string]interface{}{
		"username": githubv4.String(username),
	})
	if err != nil {
		return nil, errors.Wrap(err, "github client query")
	}

	years := cy.User.ContributionsCollection.ContributionYears

	firstYear := years[len(years)-1]

	t := time.Date(firstYear, 1, 1, 0, 0, 0, 0, time.UTC)

	return &t, nil
}

type CurrentContributionStreak struct {
	Streak        int
	Contributions int
	StartedAt     time.Time
}

func (c CurrentContributionStreak) String() string {
	msg := fmt.Sprintf("current contribution streak\nstreak: %d days", c.Streak)
	if c.Streak > 0 {
		msg = fmt.Sprintf("%s\nstarted at: %s", msg, c.StartedAt.Format("2006-01-02"))
	}

	msg = fmt.Sprintf("%s\ncontributions: %d", msg, c.Contributions)

	return msg
}

func (g *GithubService) GetCurrentContributionStreakByUsername(ctx context.Context, username string) (*CurrentContributionStreak, error) {

	options := GetContributionsByUsernameOptions{
		Username: username,
	}

	// TODO: we could decrease bandwidth by making a custom graphql request here that doesn't retrieve some of the unnecessary fields that this retrieves
	contributions, err := g.GetContributionsByUsername(ctx, options)
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

		currentContributionStreak.Contributions += d.ContributionCount

		currentContributionStreak.Streak++
	}

	return currentContributionStreak, nil
}

type LongestContributionStreak struct {
	Streak        int
	Contributions int
	StartedAt     time.Time
	EndedAt       time.Time
}

func (c LongestContributionStreak) String() string {
	startAt := c.StartedAt.Format("2006-01-02")
	msg := "longest contribution streak"

	if c.EndedAt.IsZero() {
		return fmt.Sprintf("%s\nstreak: %d days\nstarted at: %s", msg, c.Streak, startAt)
	}

	return fmt.Sprintf("%s\nstreak: %d days\nstarted at: %s\nended at: %s\ncontributions: %d", msg, c.Streak, startAt, c.EndedAt.Format("2006-01-02"), c.Contributions)
}

func (g *GithubService) GetLongestContributionStreakByUsername(ctx context.Context, username string) (*LongestContributionStreak, error) {
	year, err := g.GetFirstContributionYearByUsername(ctx, username)

	if err != nil {
		return nil, errors.Wrap(err, "get first contribution year by username")
	}

	options := GetContributionsByUsernameOptions{
		From:     *year,
		To:       time.Now(),
		Username: username,
	}

	// TODO: we could decrease bandwidth by making a custom graphql request here that doesn't retrieve some of the unnecessary fields that this retrieves
	contributions, err := g.GetContributionsByUsername(ctx, options)
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
	contributionCount := 0
	for _, d := range days {
		if d.ContributionCount == 0 {
			if streak > longestContributionStreak.Streak {
				longestContributionStreak.Streak = streak
				longestContributionStreak.StartedAt = startedAt
				longestContributionStreak.EndedAt = endedAt
				longestContributionStreak.Contributions = contributionCount
			}

			contributionCount = 0
			streak = 0

			continue
		}

		if streak == 0 {
			endedAt = d.Date
		}
		startedAt = d.Date

		contributionCount += d.ContributionCount
		streak++

	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if longestContributionStreak.EndedAt == now {
		longestContributionStreak.EndedAt = time.Time{}
	}

	return longestContributionStreak, nil
}

type TotalContribution struct {
	Total int
}

func (c TotalContribution) String() string {
	return fmt.Sprintf("total github contributions: %d", c.Total)
}

func (g *GithubService) GetTotalContributionsByUsername(ctx context.Context, username string) (*TotalContribution, error) {
	year, err := g.GetFirstContributionYearByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "get first contribution year by username")
	}

	options := GetContributionsByUsernameOptions{
		From:     *year,
		Username: username,
	}

	contributions, err := g.GetContributionsByUsername(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "get contributions by username")
	}

	totalContributions := &TotalContribution{Total: contributions.TotalContributions}

	return totalContributions, nil
}
