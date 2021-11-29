package universalinspirationalquotes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	baseURL = "https://healthruwords.p.rapidapi.com/v1"

	defaultRequestTimeout = 15 * time.Second
)

var _ Client = (*HTTPClient)(nil)

type HTTPClient struct {
	httpClient *http.Client

	apiKey  string
	baseURL string
}

type APIError struct {
	Err     error  `json:"-"`
	Message string `json:"message"`

	HTTPResponse *http.Response
}

func (a *APIError) Error() string {
	message := "no detailed message"
	if len(a.Message) > 0 {
		message = a.Message
	}

	var status string
	if a.HTTPResponse != nil {
		status = a.HTTPResponse.Status
	}

	return fmt.Sprintf("healthruwords API error (%s): %s", status, message)
}

func NewHTTPClient(httpClient *http.Client, apiKey string) *HTTPClient {
	return &HTTPClient{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

func (h *HTTPClient) request(ctx context.Context, method, path string, body io.Reader, out interface{}) (*http.Response, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultRequestTimeout)
		defer cancel()
	}

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	reqURL := fmt.Sprintf("%s%s", h.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Host", "healthruwords.p.rapidapi.com")
	req.Header.Set("X-RapidAPI-Key", h.apiKey)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return res, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	res.Body.Close()

	responseError := res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices

	if responseError {
		err := &APIError{}

		//nolint
		json.Unmarshal(resBody, err)

		err.HTTPResponse = res

		return res, err
	}

	if out != nil {
		err = json.Unmarshal(resBody, out)
		if err != nil {
			return res, errors.Wrap(err, "unmarshaling response body")
		}
	}

	return res, nil
}

func (h *HTTPClient) Get(ctx context.Context, path string, query url.Values, out interface{}) (*http.Response, error) {
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	return h.request(ctx, "GET", path, nil, out)
}
