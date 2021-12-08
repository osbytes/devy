package github

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

type Language struct {
	Size  int
	Name  string
	Color string
}

type Languages []Language

func (l Languages) String() string {
	str := ""

	for _, lang := range l {
		str += fmt.Sprintf("**%s** %d bytes\n", lang.Name, lang.Size)
	}

	return str
}

func (g *GithubService) GetLanguagesByUsername(ctx context.Context, username string) (Languages, error) {

	languageMap := map[string]*Language{}

	var afterCursor string

	for {

		var query struct {
			User struct {
				Repositories struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
					Nodes []struct {
						ID        string
						Name      string
						Languages struct {
							Edges []struct {
								Size int
								Node struct {
									Color string
									Name  string
								}
							}
						} `graphql:"languages(first: 100, orderBy: {field: SIZE, direction: DESC})"` // If someone is using over 100 different languages in a repo, the results will be incorrect. I'm okay with this because you might be a monster if you have over 100 languages in a single repo.
					}
				} `graphql:"repositories(ownerAffiliations: OWNER, isFork: false, after: $after, first: 100)"`
			} `graphql:"user(login: $username)"`
		}

		params := map[string]interface{}{
			"username": githubv4.String(username),
		}

		if afterCursor == "" {
			params["after"] = (*githubv4.String)(nil)
		} else {
			params["after"] = githubv4.String(afterCursor)
		}

		err := g.githubClient.Query(ctx, &query, params)
		if err != nil {
			return nil, errors.Wrap(err, "github client query")
		}

		for _, repository := range query.User.Repositories.Nodes {

			for _, lang := range repository.Languages.Edges {

				if _, ok := languageMap[lang.Node.Name]; !ok {
					languageMap[lang.Node.Name] = &Language{
						Name:  lang.Node.Name,
						Color: lang.Node.Color,
					}
				}

				languageMap[lang.Node.Name].Size += lang.Size

			}

		}

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}

		afterCursor = query.User.Repositories.PageInfo.EndCursor

	}

	languages := []Language{}
	for _, l := range languageMap {
		languages = append(languages, *l)
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Size > languages[j].Size
	})

	return languages, nil
}
