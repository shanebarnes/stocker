package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSymbolSearchUrl(t *testing.T) {
	url, err := createSymbolSearchUrl("", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=&apikey=", url)

	url, err = createSymbolSearchUrl("AAPL", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=AAPL&apikey=", url)

	url, err = createSymbolSearchUrl("", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=&apikey=test", url)

	url, err = createSymbolSearchUrl("AAPL", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=AAPL&apikey=test", url)
}

func TestSymbolSearch(t *testing.T) {
	info, err := SymbolSearch("", "test")
	assert.NotNil(t, err)
	assert.Nil(t, info)

	info, err = SymbolSearch("AAPL", "test")
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "AAPL", info.Symbol)
	assert.Equal(t, "Apple Inc.", info.Name)
	assert.Equal(t, "Equity", info.Type)
	assert.Equal(t, "United States", info.Region)
	assert.Equal(t, "09:30", info.MarketOpen)
	assert.Equal(t, "16:00", info.MarketClose)
	assert.Equal(t, "UTC-04", info.Timezone)
	assert.Equal(t, "USD", info.Currency)
	assert.Equal(t, "1.0000", info.MatchScore)
}
