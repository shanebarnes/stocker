package alphavantage

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	apiCurrencyExchangeRate = `https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency={{.FromCurrency}}&to_currency={{.ToCurrency}}&apikey={{.ApiKey}}`
)

type tplCurrencyExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	ApiKey       string
}

type ExchangeRate struct {
	FromCode      string `json:"1. From_Currency Code"`
	FromName      string `json:"2. From_Currency Name"`
	ToCode        string `json:"3. To_Currency Code"`
	ToName        string `json:"4. To_Currency Name"`
	ExchangeRate  string `json:"5. Exchange Rate"`
	LastRefreshed string `json:"6. Last Refreshed"`
	Timezone      string `json:"7. Time Zone"`
	BidPrice      string `json:"8. Bid Price"`
	AskPrice      string `json:"9. Ask Price"`
}

type exchangeRate struct {
	Rate ExchangeRate `json:"Realtime Currency Exchange Rate"`
}

func createCurrencyExchangeRateUrl(fromCurrency, toCurrency, apiKey string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplCurrencyExchangeRate{FromCurrency: fromCurrency, ToCurrency: toCurrency, ApiKey: apiKey}

	if tpl, err = template.New("api").Parse(apiCurrencyExchangeRate); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetCurrencyExchangeRate(fromCurrency, toCurrency, apiKey string) (*ExchangeRate, error) {
	var rate *ExchangeRate

	url, err := createCurrencyExchangeRateUrl(fromCurrency, toCurrency, apiKey)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url); err == nil {
			er := exchangeRate{}
			if err = json.Unmarshal(body, &er); err == nil {
				rate = &er.Rate
			}
		}
	}

	return rate, err
}
