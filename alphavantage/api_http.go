package alphavantage

import (
	"io/ioutil"
	"net/http"
)

func ApiGetResponseBody(url string) ([]byte, error) {
	var body []byte
	res, err := http.Get(url)
	if err == nil {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
	}

	return body, err
}
