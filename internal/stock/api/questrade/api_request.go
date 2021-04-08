package questrade

import (
	"net/http"
	"strconv"
)

const ApiRateLimitRemaining = "X-RateLimit-Remaining"

func isApiResponseRetryable(res *http.Response) bool {
	return (isResponseApiLimit(res) || isResponseRetryable(res))
}

// see https://www.questrade.com/api/documentation/rate-limiting
func isResponseApiLimit(res *http.Response) bool {
	limitReached := false
	if val := res.Header.Get(ApiRateLimitRemaining); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			limitReached = (i == 0)
		}
	}
	return limitReached
}

func isResponseRetryable(res *http.Response) bool {
	return (res.StatusCode < http.StatusBadRequest || res.StatusCode >= http.StatusInternalServerError)
}
