package questrade

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"
)

const (
	apiSymbolSearch = `https://{{.ApiServer}}/v1/symbols/search?prefix={{.Keywords}}`
)

type tplSymbolSearch struct {
	ApiKey    string
	ApiServer string
	Keywords  string
}

type SymbolSearchMatch struct {
	Symbol          string `json:"symbol"`
	SymbolId        int    `json:"symbolId"`
	Description     string `json:"description"`
	SecurityType    string `json:"securityType"`
	ListingExchange string `json:"listingExchange"`
	IsTradable      bool   `json:"isTradable"`
	IsQuotable      bool   `json:"isQuotable"`
	Currency        string `json:"currency"`
}

type symbolSearch struct {
	Symbols []SymbolSearchMatch `json:"symbols"`
}

func createSymbolSearchUrl(symbol, apiKey, apiServer string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplSymbolSearch{ApiKey: apiKey, ApiServer: apiServer, Keywords: symbol}

	if tpl, err = template.New("api").Parse(apiSymbolSearch); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetSymbolSearch(symbol, apiKey, apiServer string) (*SymbolSearchMatch, error) {
	var match *SymbolSearchMatch

	url, err := createSymbolSearchUrl(symbol, apiKey, apiServer)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url, apiKey); err == nil {
			search := symbolSearch{}
			if err = json.Unmarshal(body, &search); err == nil {
				if len(search.Symbols) > 0 {
					// TODO: look for exact match?
					match = &search.Symbols[0]
				} else {
					err = errors.New("SymbolSearch: no matches found for " + symbol)
				}
			} else {
				fmt.Println("raw body: ", string(body))
			}
		}
	}

	return match, err
}
