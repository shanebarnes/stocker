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
	ALPHA_VANTAGE_API = "https://www.alphavantage.co/query?"
)

type Asset struct {
	Symbol      string  `json:"symbol"`
	ohlcAverage float64 `json:"avgPrice"`
	MarketValue float64 `json:"marketValue"`
	Allocation  float64 `json:"%portfolio"`
	Quantity    int64   `json:"quantity"`
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
		getStockTimeSeries(*symbol, *apiKey)
	} else if len(*portfolio) > 0 {
		getCurrencyExchangeRate("USD", "CAD", *apiKey)
		p, _ := getPortfolio(*portfolio)
		allocation := float64(0)
		for i, asset := range p.Assets {
			p.Assets[i].ohlcAverage, _ = getStockTimeSeries(asset.Symbol, *apiKey)

			if p.Assets[i].Quantity > 0 {
				p.Cash += float64(p.Assets[i].Quantity) * p.Assets[i].ohlcAverage
				p.Assets[i].Quantity = 0
			}

			if asset.Allocation >= 0 {
				allocation += asset.Allocation
			}
		}

		if allocation == 100 {
			cash := p.Cash
			// Rebalance portfolio
			for i, asset := range p.Assets {
				cashAllowance := float64(0)
				if asset.Allocation > 0 {
					cashAllowance = p.Cash * asset.Allocation / 100
				}

				if asset.ohlcAverage > 0 {
					p.Assets[i].Quantity = int64(cashAllowance / asset.ohlcAverage)
					p.Assets[i].MarketValue = float64(p.Assets[i].Quantity) * asset.ohlcAverage
					cash -= p.Assets[i].MarketValue
				}
			}

			p.Cash = cash

			fmt.Println("Balanced Portfolio:", getPrettyString(p))
		} else {
			log.Fatal("Portfolio asset allocations do not add up to 100%")
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

func getCurrencyExchangeRate(from, to, key string) (float64, error) {
	function := "CURRENCY_EXCHANGE_RATE"
	api := ALPHA_VANTAGE_API + "function=" + function + "&from_currency=" + from + "&to_currency=" + to + "&apikey=" + key

	res, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	jsonBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBody))

	return -1, err
}

func getStockTimeSeries(symbol, key string) (float64, error) {
	var avg float64

	function := "TIME_SERIES_INTRADAY"
	interval := "5min"
	api := ALPHA_VANTAGE_API + "function=" + function + "&symbol=" + symbol + "&interval=" + interval + "&apikey=" + key

	res, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	jsonBody, err := ioutil.ReadAll(res.Body)

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
		avg, _ = getOhlcAverage(val)
		fmt.Println("OHLC Average:", avg)
	} else {
		fmt.Println(ts)
	}

	return avg, err
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
