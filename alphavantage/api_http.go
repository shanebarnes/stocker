package alphavantage

import (
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiCallsPerMinLimit = 5 // Free API key only allows for 5 API calls per minute
	apiCallInterval = time.Minute / apiCallsPerMinLimit + time.Second * 3 // Add 3 additional seconds to prevent API call failures
)

var lastApiCall = time.Now().Add( -1 *apiCallInterval)

func ApiGetResponseBody(url string) ([]byte, error) {
	var body []byte

	if elapsed := time.Since(lastApiCall); elapsed < apiCallInterval {
		time.Sleep(apiCallInterval - elapsed)
	}
	lastApiCall = time.Now()

	res, err := http.Get(url)
	// TODO: check that res.Status == http.StatusOK?
	if err == nil {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
	}

	return body, err
}
