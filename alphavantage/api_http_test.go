package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiGetResponseBody(t *testing.T) {
	body, err := ApiGetResponseBody("")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(body))

	body, err = ApiGetResponseBody("https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=AAPL&apikey=test")
	assert.Nil(t, err)
	assert.Less(t, 0, len(body))
}
