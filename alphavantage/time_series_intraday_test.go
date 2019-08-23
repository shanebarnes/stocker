package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTimeSeriesIntradayUrl(t *testing.T) {
	url, err := createTimeSeriesIntradayUrl("", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=&interval=5minapikey=", url)

	url, err = createTimeSeriesIntradayUrl("AAPL", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=AAPL&interval=5minapikey=", url)

	url, err = createTimeSeriesIntradayUrl("", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=&interval=5minapikey=test", url)

	url, err = createTimeSeriesIntradayUrl("AAPL", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=AAPL&interval=5minapikey=test", url)
}
