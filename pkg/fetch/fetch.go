package fetch

import (
	"log/slog"
	"net/http"

	"github.com/anti-raid/evil-befall/pkg/state"
)

type ExtraFetchOptions struct {
	// Whether to wait if we are ratelimited
	NoWait bool

	// Function to call on ratelimit
	OnRatelimit func(route string, retryAfter int, err error)
}

var DefaultFetchOptions = ExtraFetchOptions{
	OnRatelimit: func(route string, retryAfter int, err error) {
		slog.Info("Ratelimited", slog.String("route", route), slog.Int("retryAfter", retryAfter), slog.String("err", err.Error()))
	},
}

type ClientResponse struct {
	resp      *http.Response
	errorType string
}

func NewClientResponse(resp *http.Response) ClientResponse {
	return ClientResponse{
		resp:      resp,
		errorType: resp.Header.Get("X-Error-Type"),
	}
}

// Returns the error type. This is a function to make the inner value read-only
func (c *ClientResponse) ErrorType() string {
	return c.errorType
}

func (c *ClientResponse) Status() int {
	return c.resp.StatusCode
}

func (c *ClientResponse) Ok() bool {
	if c.errorType != "" {
		return false // Even if status is mistakenly 2xx, it is still an error
	}

	return c.resp.StatusCode >= 200 && c.resp.StatusCode < 300
}

func Fetch(
	sfo state.StateFetchOptions,
	efo ExtraFetchOptions,
	url string,
) (*ClientResponse, error) {
	return nil, nil
}
