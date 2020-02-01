package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSymbolQuoteUrl(t *testing.T) {
	url, err := createSymbolQuoteUrl("", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=&apikey=", url)

	url, err = createSymbolQuoteUrl("AAPL", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=AAPL&apikey=", url)

	url, err = createSymbolQuoteUrl("", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=&apikey=test", url)

	url, err = createSymbolQuoteUrl("AAPL", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=AAPL&apikey=test", url)
}

func TestGetSymbolQuote(t *testing.T) {
	quote, err := GetSymbolQuote("", ApiGetKeyFromEnv())
	assert.Nil(t, err)
	assert.NotNil(t, quote)
	assert.Equal(t, "", quote.Symbol)
	assert.Equal(t, "", quote.Open)
	assert.Equal(t, "", quote.High)
	assert.Equal(t, "", quote.Low)
	assert.Equal(t, "", quote.Price)
	assert.Equal(t, "", quote.Volume)
	assert.Equal(t, "", quote.LatestTradingDay)
	assert.Equal(t, "", quote.PreviousClose)
	assert.Equal(t, "", quote.Change)
	assert.Equal(t, "", quote.ChangePercent)

	quote, err = GetSymbolQuote("AAPL", ApiGetKeyFromEnv())
	assert.Nil(t, err)
	assert.NotNil(t, quote)
	assert.NotEqual(t, "", quote.Symbol)
	assert.NotEqual(t, "", quote.Open)
	assert.NotEqual(t, "", quote.High)
	assert.NotEqual(t, "", quote.Low)
	assert.NotEqual(t, "", quote.Price)
	assert.NotEqual(t, "", quote.Volume)
	assert.NotEqual(t, "", quote.LatestTradingDay)
	assert.NotEqual(t, "", quote.PreviousClose)
	assert.NotEqual(t, "", quote.Change)
	assert.NotEqual(t, "", quote.ChangePercent)
}
