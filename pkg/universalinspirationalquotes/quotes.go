package universalinspirationalquotes

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

const (
	ImageSizeThumbnail imageSize = "thumbnail"
	ImageSizeMedium    imageSize = "medium"
)

type imageSize string

func (i imageSize) String() string {
	return string(i)
}

type QuotesOpts struct {
	Topics    []string
	Limit     int
	ImageSize imageSize
	ID        int
}

type Quote struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	URL      string `json:"url"`
	Media    string `json:"media"`
	Category string `json:"cat"`
}

func (h *HTTPClient) Quotes(ctx context.Context, opts QuotesOpts) ([]*Quote, error) {

	query := url.Values{}

	if opts.ID != 0 {
		query.Add("id", strconv.Itoa(opts.ID))
	}

	if len(opts.ImageSize) > 0 {
		query.Add("size", opts.ImageSize.String())
	}

	if opts.Limit > 0 {
		query.Add("maxR", strconv.Itoa(opts.Limit))
	}

	if len(opts.Topics) > 0 {
		query.Add("t", strings.Join(opts.Topics, ","))
	}

	quotes := []*Quote{}

	_, err := h.Get(ctx, "/quotes/", query, &quotes)
	if err != nil {
		return nil, err
	}

	return quotes, nil
}
