package quotes

import "context"

type NOOPQuoteService struct {
}

var _ QuoteServicer = (*NOOPQuoteService)(nil)

func (n *NOOPQuoteService) GetByTopic(ctx context.Context, topic string) (*Quote, error) {
	return nil, nil
}
