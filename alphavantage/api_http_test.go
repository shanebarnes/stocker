package alphavantage

import (
	"os"
	"testing"

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

func TestApiGetResponseBody(t *testing.T) {
	body, err := ApiGetResponseBody("")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(body))

	body, err = ApiGetResponseBody("https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=AAPL&apikey=" + ApiGetKeyFromEnv())
	assert.Nil(t, err)
	assert.Less(t, 0, len(body))
}
