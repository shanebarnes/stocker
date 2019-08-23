package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	av "github.com/shanebarnes/stocker/alphavantage"
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
		//GetCurrencyExchangeRate("USD", "CAD", *apiKey)
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

func getStockTimeSeries(symbol, key string) (float64, error) {
	var avg float64

	ts, err := av.GetTimeSeriesIntraday(symbol, key)
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

func getOhlcAverage(ts av.TimeSeries) (float64, error) {
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
