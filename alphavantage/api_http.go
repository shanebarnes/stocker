package alphavantage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	ApiKeyEnvName = "AV_API_KEY"
)

var (
	ApiRequestsPerMinLimit = 0
	apiLastRequestTime = time.Time{}
)

type apiNote struct {
	Note string `json:"Note"`
}

func ApiGetKeyFromEnv() string {
	return os.Getenv(ApiKeyEnvName)
}

func apiGetRequestInterval() time.Duration {
	dur := time.Duration(0)
	if ApiRequestsPerMinLimit > 0 {
		dur = time.Minute / time.Duration(ApiRequestsPerMinLimit) + time.Second * 3 // Add 3 additional seconds to prevent API call failures
	}
	return dur
}

func ApiGetResponseBody(url string) ([]byte, error) {
	var body []byte

	interval := apiGetRequestInterval()
	if elapsed := time.Since(apiLastRequestTime); elapsed < interval {
		time.Sleep(interval - elapsed)
	}
	apiLastRequestTime = time.Now()

	res, err := http.Get(url)
	// TODO: check that res.Status == http.StatusOK?
	if err == nil {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)

		// A 200 status code is returned when the API call limit is reached.
		// Inspect response body for API call limit "note".
		err = apiIsRequestLimitError(body)
	}
	return body, err
}

func apiIsRequestLimitError(body []byte) error {
	note := apiNote{}
	err := json.Unmarshal(body, &note)
	if err == nil && len(note.Note) > 0 {
		err = errors.New(note.Note)
	} else {
		err = nil
	}
	return err
}
