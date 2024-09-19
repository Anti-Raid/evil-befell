// Package `fetch` provides functions for fetching data from the Anti-Raid API while handling all the painful stuff such as error handling, properly reading+closing responses etc.

package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/evil-befall/types"
	"github.com/anti-raid/evil-befall/types/silverpelt"
)

var (
	// To allow reusing the same client, just define a global one
	FetchHttpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

var (
	ErrUnmarshalError    = errors.New("fetch: failed to unmarshal response")
	ErrServerMaintenance = errors.New("fetch: server currently undergoing maintenance")
)

type ExtraFetchOptions struct {
	// Whether to wait if we are ratelimited
	NoWait bool

	// Whether or not any extra-but-possibly-needed headers should be added to the sent headers
	NoExtraHeaders bool

	// Any other headers to add to the request
	Headers map[string]string

	// What session, if any, to use for the request
	Session *state.StateSessionAuth

	// Function to call on ratelimit
	OnRatelimit func(fo FetchOptions, retryAfter float64, err error, sfo *state.StateFetchOptions, sess *state.StateSessionAuth)

	// Whether or not to error on fail
	NoErrorOnFail bool
}

var DefaultFetchOptions = ExtraFetchOptions{
	OnRatelimit: func(fo FetchOptions, retryAfter float64, err error, sfo *state.StateFetchOptions, sess *state.StateSessionAuth) {
		slog.Info("Ratelimited", slog.String("req", fo.String()), slog.Float64("retryAfter", retryAfter), slog.String("err", err.Error()), slog.Bool("isAuthorized", sess != nil && sess.IsAuthorized()))
	},
}

type ClientResponse struct {
	resp      *http.Response
	errorType string
}

func NewClientResponse(resp *http.Response) *ClientResponse {
	return &ClientResponse{
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

// Unmarshals the body, wrapping ErrUnmarshalError if an error
func (c *ClientResponse) unmarshalBody(t any) error {
	//nolint:errcheck
	defer c.resp.Body.Close()
	//nolint:errcheck
	defer io.Copy(io.Discard, c.resp.Body)

	bytes, err := io.ReadAll(c.resp.Body)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalError, err)
	}

	if err := json.Unmarshal(bytes, t); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalError, err)
	}

	return nil
}

func (c *ClientResponse) formatApiError(base string, err types.ApiError) string {
	if len(err.Context) > 0 {
		return fmt.Sprintf("%v: %v [%v]", base, err.Message, fmtMap(err.Context))
	}

	return fmt.Sprintf("%v: %v", base, err.Message)
}

func (c *ClientResponse) Err() error {
	if !c.Ok() {
		panic("fetch: tried to get error from non-error response")
	}

	switch c.errorType {
	case "permission_check":
		var pr *silverpelt.PermissionResult

		if err := c.unmarshalBody(&pr); err != nil {
			return err
		}

		prf := NewPermissionResultFormatter(*pr)

		return errors.New(prf.ToMarkdown())
	case "settings_error":
		var se *silverpelt.CanonicalSettingsError

		if err := c.unmarshalBody(&se); err != nil {
			return err
		}

		sef := NewSettingsErrorFormatter(*se)

		return errors.New(sef.ToMarkdown())
	}

	var apiErr types.ApiError

	if err := c.unmarshalBody(&apiErr); err != nil {
		return err
	}

	return errors.New(c.formatApiError("API error", apiErr))
}

func (c *ClientResponse) Json(t any) error {
	if !c.Ok() {
		panic("fetch: tried to get JSON from non-OK response")
	}

	return c.unmarshalBody(t)
}

type FetchOptions struct {
	Method string
	URL    string
	Body   io.ReadSeeker // Must be a read seeker for ratelimits etc to work
}

func (fo FetchOptions) String() string {
	return fmt.Sprintf("FetchOptions{Method: %v, URL: %v}", fo.Method, fo.URL)
}

func Fetch(
	ctx context.Context,
	sfo *state.StateFetchOptions,
	efo ExtraFetchOptions,
	opts FetchOptions,
) (*ClientResponse, error) {
	for {
		var headers = map[string]string{}

		if !efo.NoExtraHeaders {
			headers["Content-Type"] = "application/json"
		}

		if len(efo.Headers) > 0 {
			for k, v := range efo.Headers {
				headers[k] = v
			}
		}

		if efo.Session != nil {
			sess, err := efo.Session.GetCurrentSession()

			if err != nil {
				return nil, err
			}

			headers["Authorization"] = fmt.Sprintf("User %v", sess.Token)
		}

		req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, opts.Body)

		if err != nil {
			return nil, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := FetchHttpClient.Do(req)

		if err != nil {
			return nil, err
		}

		if slices.Contains([]int{408, 502, 503, 504}, resp.StatusCode) {
			return nil, ErrServerMaintenance
		}

		retryAfterStr := resp.Header.Get("Retry-After")

		if retryAfterStr != "" {
			// NOTE: We use milliseconds here even though the API *currently* returns seconds
			// to make it easier to change in the future
			retryAfter, err := strconv.ParseFloat(retryAfterStr, 64)

			if err != nil {
				retryAfter = 3000
			} else {
				retryAfter *= 1000
			}

			if efo.OnRatelimit != nil {
				efo.OnRatelimit(opts, retryAfter, err, sfo, efo.Session)
			}

			// Wait for the time specified by the server
			if !efo.NoWait {
				time.Sleep(time.Duration(retryAfter) * time.Second)

				if opts.Body != nil {
					if _, err := opts.Body.Seek(0, io.SeekStart); err != nil {
						return nil, err
					}
				}

				if efo.OnRatelimit != nil {
					efo.OnRatelimit(opts, 0, err, sfo, efo.Session)
				}

				continue
			}
		}

		ncr := NewClientResponse(resp)

		if !efo.NoErrorOnFail && !ncr.Ok() {
			return ncr, ncr.Err()
		}

		return ncr, nil
	}
}

func JsonBody(v any) (io.ReadSeeker, error) {
	b, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
