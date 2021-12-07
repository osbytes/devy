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
	from := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
	to := time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)

	options := GetContributionsByUsernameOptions{
		Username: "devy",
		From: from,
		To: to,
	}
	
	githubClient.On(
		"Query",
		ctx,
		mock.AnythingOfType("*github.contributionsQuery"),
		mock.MatchedBy(func (params map[string]interface{}) bool {
			assert.Equal(githubv4.String(options.Username), params["username"])
			assert.Equal(githubv4.DateTime{ Time: from }, params["from"])
			assert.Equal(githubv4.DateTime{ Time: to }, params["to"])			
			return true
		}),
	).Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*contributionsQuery)
		(*a) = contributionsQuery{}
	}).Once()

	resp, err := githubService.GetContributionsByUsername(ctx, options)

	assert.NoError(err)

	fmt.Println(resp)

	githubClient.AssertExpectations(t)
}

func TestGetFirstContributionYearByUsername(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(assert)
}