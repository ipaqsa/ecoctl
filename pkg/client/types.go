package client

import (
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// APIKey token for client
	APIKey string

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string

	Service CloudOrchestrator
	// Optional retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
	// that fail with 429 or 500-level response codes
	RetryConfig RetryConfig
}

// Opt are options for New.
type Opt func(*Client) error

type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *float64    // Minimum time to wait
	RetryWaitMax *float64    // Maximum time to wait
	Logger       interface{} // Customer logger instance. Must implement either go-retryablehttp.Logger or go-retryablehttp.LeveledLogger
}

type Response struct {
	*http.Response
}

// An ResponseError reports the error caused by an API request.
type ResponseError struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// Attempts is the number of times the request was attempted when retries are enabled.
	Attempts int
}
