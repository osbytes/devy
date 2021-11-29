package quotes

import (
	"bot/pkg/universalinspirationalquotes"
	"context"
)

type QuoteServicer interface {
	GetByTopic(ctx context.Context, topic string) (*Quote, error)
}

type QuoteService struct {
	quoteClient universalinspirationalquotes.Client
}

var _ QuoteServicer = (*QuoteService)(nil)

func NewQuoteService(quoteClient universalinspirationalquotes.Client) *QuoteService {
	return &QuoteService{
		quoteClient: quoteClient,
	}
}
