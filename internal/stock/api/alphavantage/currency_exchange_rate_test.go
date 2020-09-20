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

func TestGetCurrencyExchangeRateInfo(t *testing.T) {
	rate, err := GetCurrencyExchangeRateInfo("", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, rate)
	assert.Equal(t, "", rate.FromCode)
	assert.Equal(t, "", rate.FromName)
	assert.Equal(t, "", rate.ToCode)
	assert.Equal(t, "", rate.ToName)
	assert.Equal(t, "", rate.ExchangeRate)
	assert.Equal(t, "", rate.LastRefreshed)
	assert.Equal(t, "", rate.Timezone)
	assert.Equal(t, "", rate.BidPrice)
	assert.Equal(t, "", rate.AskPrice)

	rate, err = GetCurrencyExchangeRateInfo("USD", "CAD", "test")
	assert.Nil(t, err)
	assert.NotNil(t, rate)
	assert.NotEqual(t, "", rate.FromCode)
	assert.NotEqual(t, "", rate.FromName)
	assert.NotEqual(t, "", rate.ToCode)
	assert.NotEqual(t, "", rate.ToName)
	assert.NotEqual(t, "", rate.ExchangeRate)
	assert.NotEqual(t, "", rate.LastRefreshed)
	assert.NotEqual(t, "", rate.Timezone)
	assert.NotEqual(t, "", rate.BidPrice)
	assert.NotEqual(t, "", rate.AskPrice)
}
