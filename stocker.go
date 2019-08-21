package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ALPHA_VANTAGE_API_STOCK_TIME_SERIES = "https://www.alphavantage.co/query?"
)

type TsMetaData struct {
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
	MetaData   TsMetaData            `json:"Meta Data"`
	Ts         map[string]TimeSeries `json:"Time Series (5min)"`
}

func main() {
	symbol := flag.String("symbol", "", "stock symbol")
	apiKey := flag.String("apiKey", "", "Alpha Vantage API key")
	flag.Parse()

	if len(*symbol) == 0 || len(*apiKey) == 0 {
		flag.PrintDefaults()
	} else {
		getStockTimeSeries("TIME_SERIES_INTRADAY", *symbol, "5min", *apiKey)
	}
}

func getStockTimeSeries(function, symbol, interval, key string) {
	api := ALPHA_VANTAGE_API_STOCK_TIME_SERIES + "function=" + function + "&symbol=" + symbol + "&interval=" + interval + "&apikey=" + key

	res, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	jsonBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	ts := TsIntraday{}
	err = json.Unmarshal([]byte(jsonBody), &ts)
	if err != nil {
		log.Fatal(err)
	}

	if val, ok := ts.Ts[ts.MetaData.LastRefreshed]; ok {
		fmt.Println(ts.MetaData.Symbol, ", ", ts.MetaData.LastRefreshed, ":", getPrettyString(val))
	} else {
		fmt.Println(ts)
	}
}

func getPrettyString(v interface{}) string {
	str := ""
	buf, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		str = string(buf)
	}
	return str
}
