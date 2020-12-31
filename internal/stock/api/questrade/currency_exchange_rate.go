package questrade

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	apiCurrencyExchangeRate = `https://api.exchangeratesapi.io/latest?base={{.FromCurrency}}&symbols={{.ToCurrency}}`
)

type tplCurrencyExchangeRate struct {
	FromCurrency string
	ToCurrency   string
}

type ExchangeRate struct {
	BaseSymbol string             `json:"base"`
	Date       string             `json:"date"`
	Rates      map[string]float64 `json:"rates"`
}

func createCurrencyExchangeRateUrl(fromCurrency, toCurrency string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplCurrencyExchangeRate{FromCurrency: fromCurrency, ToCurrency: toCurrency}

	if tpl, err = template.New("api").Parse(apiCurrencyExchangeRate); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetCurrencyExchangeRate(fromCurrency, toCurrency, apiKey string) (float64, error) {
	var xr float64

	xri, err := GetCurrencyExchangeRateInfo(fromCurrency, toCurrency, apiKey)
	if err == nil {
		if f, ok := xri.Rates[toCurrency]; ok {
			xr = f
		} else {
			// add error
		}
	}

	return xr, err
}

func GetCurrencyExchangeRateInfo(fromCurrency, toCurrency, apiKey string) (*ExchangeRate, error) {
	var rate *ExchangeRate

	url, err := createCurrencyExchangeRateUrl(fromCurrency, toCurrency)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url, ""); err == nil {
			er := ExchangeRate{}
			if err = json.Unmarshal(body, &er); err == nil {
				rate = &er
			}
		}
	}

	return rate, err
}
