package client

import (
	"fmt"
	"net/url"
	"strings"
)

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) Opt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u

		return nil
	}
}

// SetAPIKey is a client option for setting the APIKey token.
func SetAPIKey(apiKey string) Opt {
	return func(c *Client) error {
		tokenPartsCount := 2
		parts := strings.SplitN(apiKey, " ", tokenPartsCount)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			apiKey = parts[1]
		}
		c.APIKey = apiKey
		c.headers["Authorization"] = fmt.Sprintf("apikey %s", c.APIKey)

		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) Opt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) Opt {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// WithRetryAndBackoffs sets retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
// that fail with 429 or 500-level response codes using the go-retryablehttp client.
func WithRetryAndBackoffs(retryConfig RetryConfig) Opt {
	return func(c *Client) error {
		c.RetryConfig.RetryMax = retryConfig.RetryMax
		c.RetryConfig.RetryWaitMax = retryConfig.RetryWaitMax
		c.RetryConfig.RetryWaitMin = retryConfig.RetryWaitMin
		c.RetryConfig.Logger = retryConfig.Logger
		return nil
	}
}
