package client

const (
	libraryVersion = "1.0.0"
	defaultBaseURL = "https://api.edgecenter.ru/cloud/"
	userAgent      = "edgecloud/" + libraryVersion
	mediaType      = "application/json"

	internalHeaderRetryAttempts = "X-Edgecloud-Retry-Attempts"

	defaultRetryMax     = 3
	defaultRetryWaitMax = 30
	defaultRetryWaitMin = 1
)
