package questrade

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/shanebarnes/stocker/internal/stock/api"
)

const ApiRateLimitRemaining = "X-RateLimit-Remaining"

var apiClient *http.Client = &http.Client{
	Timeout: api.DefaultClientTimeout,
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, network, address)
		},
	},
}

func ApiGetResponseBody(url, accessToken string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		if len(accessToken) > 0 {
			req.Header.Set("Authorization", "Bearer " + accessToken)
		}

		api.MakeApiRequestWithRetry(apiClient, req, func(res *http.Response, rerr error) bool {
			retry := false
			err = rerr
			if err == nil {
				body, err = ioutil.ReadAll(res.Body)
				if res.StatusCode != http.StatusOK {
					err = errors.New(fmt.Sprintf("API response status code: %d, details: %s", res.StatusCode, string(body)))
					if isResponseApiLimit(res) || isResponseRetryable(res) {
						retry = true
					}
				}
			} else {
				retry = api.IsErrorRetryable(err)
			}
			return retry
		})
	}

	return body, err
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
