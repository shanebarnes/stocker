package alphavantage

import (
	"bytes"
	"encoding/json"
	"errors"
	"text/template"
)

const (
	apiSymbolSearch = `https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords={{.Keywords}}&apikey={{.ApiKey}}`
)

type tplSymbolSearch struct {
	Keywords string
	ApiKey   string
}

type SymbolSearchMatch struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	MarketOpen  string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
	MatchScore  string `json:"9. matchScore"`
}

type symbolSearch struct {
	BestMatches []SymbolSearchMatch `json:"bestMatches"`
}

func createSymbolSearchUrl(symbol, apiKey string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplSymbolSearch{Keywords: symbol, ApiKey: apiKey}

	if tpl, err = template.New("api").Parse(apiSymbolSearch); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetSymbolSearch(symbol, apiKey string) (*SymbolSearchMatch, error) {
	var match *SymbolSearchMatch

	url, err := createSymbolSearchUrl(symbol, apiKey)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url); err == nil {
			search := symbolSearch{}
			if err = json.Unmarshal(body, &search); err == nil {
				if len(search.BestMatches) > 0 {
					// TODO: check that match score is 1.000?
					match = &search.BestMatches[0]
				} else {
					err = errors.New("SymbolSearch: no matches found for " + symbol)
				}
			}
		}
	}

	return match, err
}
