package alphavantage

import (
	"bytes"
	"encoding/json"
	"strconv"
	"text/template"
)

const (
	apiSymbolQuote = `https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol={{.Symbol}}&apikey={{.ApiKey}}`
)

type tplSymbolQuote struct {
	Symbol string
	ApiKey string
}

type SymbolQuote struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

type symbolQuote struct {
	Quote SymbolQuote `json:"Global Quote"`
}

func createSymbolQuoteUrl(symbol, apiKey string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplSymbolQuote{Symbol: symbol, ApiKey: apiKey}

	if tpl, err = template.New("api").Parse(apiSymbolQuote); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetSymbolQuote(symbol, apiKey string) (*SymbolQuote, error) {
	var quote *SymbolQuote

	url, err := createSymbolQuoteUrl(symbol, apiKey)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url); err == nil {
			sq := symbolQuote{}
			if err = json.Unmarshal(body, &sq); err == nil {
				quote = &sq.Quote
			}
		}
	}

	return quote, err
}

func GetSymbolQuotePrice(symbol, apiKey string) (float64, error) {
	var price float64

	quote, err := GetSymbolQuote(symbol, apiKey)
	if err == nil {
		price, err = strconv.ParseFloat(quote.Price, 64)
	}

	return price, err
}
