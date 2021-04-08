package questrade

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/shanebarnes/stocker/internal/stock/api"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestApiDefaultClient(t *testing.T) {
	assert.NotNil(t, api.Client)
	assert.Equal(t, api.DefaultClientTimeout, api.Client.Timeout)
}

func TestApiGetResponseBody_NoError(t *testing.T) {
	requestCount := 0
	saveClient := api.Client
	defer func() {
		api.Client = saveClient
	}()
	api.Client = NewTestClient(func(req *http.Request) *http.Response {
		requestCount++

		assert.Equal(t, "https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", req.URL.String())
		assert.Equal(t, "Bearer AccessToken01", req.Header.Get("Authorization"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
		}
	})
	body, err := api.GetApiResponseBody("https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", "AccessToken01", nil)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, 1, requestCount)
}

func TestApiGetResponseBody_NonretryableError(t *testing.T) {
	requestCount := 0
	saveClient := api.Client
	defer func() {
		api.Client = saveClient
	}()
	api.Client = NewTestClient(func(req *http.Request) *http.Response {
		requestCount++

		assert.Equal(t, "https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", req.URL.String())
		assert.Equal(t, "Bearer AccessToken01", req.Header.Get("Authorization"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("Access token is invalid")),
			Header:     make(http.Header),
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
		}
	})
	_, err := api.GetApiResponseBody("https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", "AccessToken01", nil)
	assert.NotNil(t, err)
	assert.Equal(t, "API response status code: 401, details: Access token is invalid", err.Error())
	assert.Equal(t, 1, requestCount)
}

func TestApiGetResponseBody_RetryableError(t *testing.T) {
	requestCount := 0
	saveClient := api.Client
	defer func() {
		api.Client = saveClient
	}()
	api.Client = NewTestClient(func(req *http.Request) *http.Response {
		requestCount++

		assert.Equal(t, "https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", req.URL.String())
		assert.Equal(t, "Bearer AccessToken01", req.Header.Get("Authorization"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("Something strange happened")),
			Header:     make(http.Header),
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
		}
	})
	_, err := api.GetApiResponseBody("https://api01.iq.questrade.com/v1/symbols/search?prefix=ACME", "AccessToken01", nil)
	assert.NotNil(t, err)
	assert.Equal(t, "API response status code: 500, details: Something strange happened", err.Error())
	assert.Equal(t, api.DefaultRequestRetryLimit, requestCount)
}

func TestIsResponseApiLimit_HeaderGtZero(t *testing.T) {
	res := http.Response{Header: make(map[string][]string)}
	res.Header.Add(ApiRateLimitRemaining, "1")
	assert.False(t, isResponseApiLimit(&res))
}

func TestIsResponseApiLimit_HeaderInvalid(t *testing.T) {
	res := http.Response{Header: make(map[string][]string)}
	res.Header.Add(ApiRateLimitRemaining, "abc")
	assert.False(t, isResponseApiLimit(&res))
}

func TestIsResponseApiLimit_HeaderMissing(t *testing.T) {
	res := http.Response{Header: make(map[string][]string)}
	res.Header.Add(ApiRateLimitRemaining, "")
	assert.False(t, isResponseApiLimit(&res))
}

func TestIsResponseApiLimit_HeaderZero(t *testing.T) {
	res := http.Response{Header: make(map[string][]string)}
	res.Header.Add(ApiRateLimitRemaining, "0")
	assert.True(t, isResponseApiLimit(&res))
}
