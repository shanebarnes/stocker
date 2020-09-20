package questrade

import (
	"bytes"
	"encoding/json"
	"errors"
	"text/template"
)

const (
	apiSymbolQuote = `https://{{.ApiServer}}/v1/markets/quotes?ids={{.SymbolId}}`
)

type tplSymbolQuote struct {
	ApiKey    string
	ApiServer string
	SymbolId  string
}

type SymbolQuote struct {
	Symbol              string  `json:"symbol"`
	SymbolId            int     `json:"symbolId"`
	Tier                string  `json:"tier"`
	BidPrice            float64 `json:"bidPrice"`
	BidSize             int     `json:"bidSize"`
	AskPrice            float64 `json:"askPrice"`
	AskSize             int     `json:"askSize"`
	LastTradePriceTrHrs float64 `json:"lastTradePriceTrHrs"`
	LastTradePrice      float64 `json:"lastTracePrice"`
	LastTradeSize       int     `json:"lastTradeSize"`
	LastTradeTick       string  `json:"lastTradeTick"`
	LastTradeTime       string  `json:"lastTradeTime"`
	Volume              int     `json:"volume"`
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	Delay               int     `json:"delay"`
	IsHalted            bool    `json:"isHalted"`
}

type symbolQuote struct {
	Quotes []SymbolQuote `json:"quotes"`
}

func createSymbolQuoteUrl(symbolId, apiKey, apiServer string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplSymbolQuote{ApiKey: apiKey, ApiServer: apiServer, SymbolId: symbolId}

	if tpl, err = template.New("api").Parse(apiSymbolQuote); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetSymbolQuote(symbolId, apiKey, apiServer string) (*SymbolQuote, error) {
	var quote *SymbolQuote

	url, err := createSymbolQuoteUrl(symbolId, apiKey, apiServer)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url, apiKey); err == nil {
			sq := symbolQuote{}
			if err = json.Unmarshal(body, &sq); err == nil {
				if len(sq.Quotes) > 0 {
					// TODO: look for exact match?
					quote = &sq.Quotes[0]
				} else {
					err = errors.New("SymbolQuote: no matches found for " + symbolId)
				}
			}
		}
	}

	return quote, err
}
