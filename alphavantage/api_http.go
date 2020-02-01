package alphavantage

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	ApiCallsPerMinLimit = 5 // Free API key only allows for 5 API calls per minute
	apiCallInterval = time.Minute / ApiCallsPerMinLimit + time.Second * 3 // Add 3 additional seconds to prevent API call failures
	ApiKeyEnvName = "AV_API_KEY"
)

var lastApiCall = time.Now().Add( -1 *apiCallInterval)

func ApiGetKeyFromEnv() string {
	return os.Getenv(ApiKeyEnvName)
}

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
