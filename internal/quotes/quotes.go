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

func (q *QuoteService) GetByTopic(ctx context.Context, topic string) (*Quote, error) {
	quotes, err := q.quoteClient.Quotes(ctx, universalinspirationalquotes.QuotesOpts{
		Topics: []string{topic},
		Limit:  1,
	})
	if err != nil {
		return nil, errors.Wrap(err, "quote client get quotes")
	}

	if len(quotes) == 0 {
		return nil, ErrNotFound
	}

	return &Quote{
		Title:    quotes[0].Title,
		Author:   quotes[0].Author,
		URL:      quotes[0].URL,
		Media:    quotes[0].Media,
		Category: quotes[0].Category,
	}, nil
}
