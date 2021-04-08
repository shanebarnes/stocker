package exchangerate

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/shanebarnes/stocker/internal/stock/api"
)

const (
	apiCurrencyExchangeRate = `https://api.exchangerate.host/latest?base={{.FromCurrency}}&source=imf&places=4&symbols={{.ToCurrency}}`
)

type tplCurrencyExchangeRate struct {
	FromCurrency string
	ToCurrency   string
}

type motd struct {
	Message string `json:"msg"`
	Url     string `json:"url"`
}

type ExchangeRate struct {
	BaseSymbol string             `json:"base"`
	Date       string             `json:"date"`
	Motd       motd               `json:"motd"`
	Rates      map[string]float64 `json:"rates"`
	Success    bool               `json:"success"`
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
		if body, err = api.GetApiResponseBody(url, "", nil); err == nil {
			er := ExchangeRate{}
			if err = json.Unmarshal(body, &er); err == nil {
				rate = &er
			}
		}
	}

	return rate, err
}
