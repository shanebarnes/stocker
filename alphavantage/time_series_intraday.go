package alphavantage

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	apiTimeSeriesIntraday = `https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol={{.Symbol}}&interval={{.Interval}}apikey={{.ApiKey}}`
)

type tplTimeSeriesIntraday struct {
	Symbol   string
	Interval string
	ApiKey   string
}

type TimeSeriesMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Interval      string `json:"4. Interval"`
	OutputSize    string `json:"5. Output Size"`
	TimeZone      string `json:"6. TimeZone'`
}

type TimeSeries struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TsIntraday struct {
	MetaData   TimeSeriesMetaData    `json:"Meta Data"`
	Ts         map[string]TimeSeries `json:"Time Series (5min)"`
}

func createTimeSeriesIntradayUrl(symbol, apiKey string) (string, error) {
	var url bytes.Buffer
	var err error

	var tpl *template.Template
	t := tplTimeSeriesIntraday{Symbol: symbol, Interval: "5min", ApiKey: apiKey}

	if tpl, err = template.New("api").Parse(apiTimeSeriesIntraday); err == nil {
		err = tpl.Execute(&url, t)
	}

	return url.String(), err
}

func GetTimeSeriesIntraday(symbol, apiKey string) (*TsIntraday, error) {
	var tsIntraday *TsIntraday

	url, err := createTimeSeriesIntradayUrl(symbol, apiKey)
	if err == nil {
		var body []byte
		if body, err = ApiGetResponseBody(url); err == nil {
			tsi := TsIntraday{}
			if err = json.Unmarshal(body, &tsi); err == nil {
				tsIntraday = &tsi
			}
		}
	}

	return tsIntraday, err
}
