package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	ALPHA_VANTAGE_API_STOCK_TIME_SERIES = "https://www.alphavantage.co/query?"
)

type Asset struct {
	Symbol    string  `json:"symbol"`
	Portfolio float64 `json:"%portfolio"`
	Quantity  int64   `json:"quantity"`
}

type Portfolio struct {
	Assets []Asset `json:"assets"`
	Cash   float64 `json:"cash"`
}

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
	portfolio := flag.String("portfolio", "", "portfolio file")
	symbol := flag.String("symbol", "", "stock symbol")
	apiKey := flag.String("apiKey", "", "Alpha Vantage API key")
	flag.Parse()

	if len(*apiKey) == 0 {
		flag.PrintDefaults()
	} else if len(*symbol) > 0 {
		getStockTimeSeries("TIME_SERIES_INTRADAY", *symbol, "5min", *apiKey)
	} else if len(*portfolio) > 0 {
		p, _ := getPortfolio(*portfolio)
		for _, asset := range p.Assets {
			getStockTimeSeries("TIME_SERIES_INTRADAY", asset.Symbol, "5min", *apiKey)
		}
	} else {
		flag.PrintDefaults()
	}
}

func getPortfolio(filename string) (*Portfolio, error) {
	portfolio := Portfolio{}
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		if err = json.Unmarshal([]byte(file), &portfolio); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return &portfolio, err
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
		fmt.Println("Symbol:", ts.MetaData.Symbol)
		fmt.Println("Last Refreshed:", ts.MetaData.LastRefreshed)
		fmt.Println("Time Series:", getPrettyString(val))
		avg, _ := getOhlcAverage(val)
		fmt.Println("OHLC Average:", avg)
	} else {
		fmt.Println(ts)
	}
}

func getOhlcAverage(ts TimeSeries) (float64, error) {
	var avg, o, h, l, c float64
	var err error

	if o, err = strconv.ParseFloat(ts.Open, 64); err != nil {
		// Error
	} else if h, err = strconv.ParseFloat(ts.High, 64); err != nil {
		// Error
	} else if l, err = strconv.ParseFloat(ts.Low, 64); err != nil {
		// Error
	} else if c, err = strconv.ParseFloat(ts.Close, 64); err != nil {
		// Error
	} else {
		avg = (o + h + l + c) / 4.
	}

	return avg, err
}

func getPrettyString(v interface{}) string {
	str := ""
	buf, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		str = string(buf)
	}
	return str
}
