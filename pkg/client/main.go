package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewWithRetries(httpClient *http.Client, opts ...Opt) (*Client, error) {
	opts = append(opts, WithRetryAndBackoffs(
		RetryConfig{
			RetryMax:     defaultRetryMax,
			RetryWaitMin: PtrTo(float64(defaultRetryWaitMin)),
			RetryWaitMax: PtrTo(float64(defaultRetryWaitMax)),
		},
	))
	return New(httpClient, opts...)
}

func New(httpClient *http.Client, opts ...Opt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.RetryConfig.RetryMax > 0 {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryMax = c.RetryConfig.RetryMax

		if c.RetryConfig.RetryWaitMin != nil {
			retryableClient.RetryWaitMin = time.Duration(*c.RetryConfig.RetryWaitMin * float64(time.Second))
		}
		if c.RetryConfig.RetryWaitMax != nil {
			retryableClient.RetryWaitMax = time.Duration(*c.RetryConfig.RetryWaitMax * float64(time.Second))
		}

		// By default, this is nil and does not log.
		retryableClient.Logger = c.RetryConfig.Logger

		// if timeout is set, it is maintained before overwriting client with StandardClient()
		retryableClient.HTTPClient.Timeout = c.HTTPClient.Timeout

		retryableClient.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			if resp != nil {
				resp.Header.Add(internalHeaderRetryAttempts, strconv.Itoa(numTries))

				return resp, err
			}
			return resp, err
		}
		c.HTTPClient = retryableClient.StandardClient()
	}
	return c, nil
}

// NewClient returns a new EdgecenterCloud API, using the given
// http.Client to perform all requests.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{HTTPClient: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Service = &CloudOrchestratorOp{client: c}
	c.headers = make(map[string]string)
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(_ context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.HTTPClient, req)
	if err != nil {
		return &Response{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
			},
		}, err
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		if rErr := resp.Body.Close(); err == nil {
			err = rErr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusNoContent && v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return &Response{
				Response: &http.Response{
					Status:     http.StatusText(http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
			}, err
		}
	}

	return response, err
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	return client.Do(req)
}

func (r *ResponseError) Error() string {
	var attempted string
	if r.Attempts > 0 {
		attempted = fmt.Sprintf("; giving up after %d attempt(s)", r.Attempts)
	}

	return fmt.Sprintf("%v %v: %d %v%s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, attempted)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ResponseError. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ResponseError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	attempts, strconvErr := strconv.Atoi(r.Header.Get(internalHeaderRetryAttempts))
	if strconvErr == nil {
		errorResponse.Attempts = attempts
	}

	return errorResponse
}
