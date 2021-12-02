package github

import (
	"context"
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

	return contributions, nil
}
