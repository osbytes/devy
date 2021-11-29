package universalinspirationalquotes

import "context"

type Client interface {
	Quotes(ctx context.Context, opts QuotesOpts) ([]*Quote, error)
}
