package alphavantage

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApiGetKeyFromEnv(t *testing.T) {
	// Save current environment variable before modifying
	cache := ApiGetKeyFromEnv()

	os.Setenv(ApiKeyEnvName, "")
	assert.Equal(t, "", ApiGetKeyFromEnv())

	os.Setenv(ApiKeyEnvName, "ABCD1234")
	assert.Equal(t, "ABCD1234", ApiGetKeyFromEnv())

	os.Setenv(ApiKeyEnvName, "")
	assert.Equal(t, "", ApiGetKeyFromEnv())

	// Restore original environment variable value
	os.Setenv(ApiKeyEnvName, cache)
}

func TestApiGetRequestInterval(t *testing.T) {
	assert.Equal(t, 0, ApiRequestsPerMinLimit)

	ApiRequestsPerMinLimit = -1
	assert.Equal(t, 0 * time.Second, apiGetRequestInterval())

	ApiRequestsPerMinLimit = 0
	assert.Equal(t, 0 * time.Second, apiGetRequestInterval())

	ApiRequestsPerMinLimit = 1
	assert.Equal(t, 60 * time.Second + 3 * time.Second, apiGetRequestInterval())

	ApiRequestsPerMinLimit = 5
	assert.Equal(t, 12 * time.Second + 3 * time.Second, apiGetRequestInterval())
}

func TestApiGetResponseBody(t *testing.T) {
	body, err := ApiGetResponseBody("")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(body))

	body, err = ApiGetResponseBody("https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=AAPL&apikey=" + ApiGetKeyFromEnv())
	assert.Nil(t, err)
	assert.Less(t, 0, len(body))
}

func TestApiIsRequestLimitError(t *testing.T) {
	note := apiNote{
		Note: "Thank you for using Alpha Vantage! Our standard API call frequency is 5 calls per minute and 500 calls per day. Please visit https://www.alphavantage.co/premium/ if you would like to target a higher API call frequency.",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(note)
	err := apiIsRequestLimitError(buf.Bytes())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), note.Note)

	match := SymbolSearchMatch{}
	buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(match)
	err = apiIsRequestLimitError(buf.Bytes())
	assert.Nil(t, err)
}
