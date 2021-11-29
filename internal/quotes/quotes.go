package quotes

import (
	"bot/pkg/universalinspirationalquotes"
	"context"

	"github.com/pkg/errors"
)

type Quote struct {
	Title    string
	Author   string
	URL      string
	Media    string
	Category string
}

func (qs *QuoteService) GetByTopic(ctx context.Context, topic string) (*Quote, error) {
	quotes, err := qs.quoteClient.Quotes(ctx, universalinspirationalquotes.QuotesOpts{
		Topics: []string{topic},
		Limit:  1,
	})
	if err != nil {
		return nil, errors.Wrap(err, "quote client get quotes")
	}

	if len(quotes) == 0 {
		return nil, ErrNotFound
	}

	q := quotes[0]

	return &Quote{
		Title:    q.Title,
		Author:   q.Author,
		URL:      q.URL,
		Media:    q.Media,
		Category: q.Category,
	}, nil
}
