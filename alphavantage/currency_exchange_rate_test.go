package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCurrencyExchangeRateUrl(t *testing.T) {
	url, err := createCurrencyExchangeRateUrl("", "", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=&to_currency=&apikey=", url)

	url, err = createCurrencyExchangeRateUrl("USD", "", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=USD&to_currency=&apikey=", url)

	url, err = createCurrencyExchangeRateUrl("", "CAD", "")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=&to_currency=CAD&apikey=", url)

	url, err = createCurrencyExchangeRateUrl("USD", "CAD", "test")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=USD&to_currency=CAD&apikey=test", url)
}

func TestCurrencyExchangeRate(t *testing.T) {
	rate, err := CurrencyExchangeRate("", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, rate)

	rate, err = CurrencyExchangeRate("USD", "CAD", "test")
	assert.Nil(t, err)
	assert.NotNil(t, rate)
	//assert.Equal(t, "USD", rate.FromCode)
	//assert.Equal(t, "United States Dollar", rate.FromName)
	//assert.Equal(t, "CAD", rate.ToCode)
}
