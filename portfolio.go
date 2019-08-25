package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"strings"

	av "github.com/shanebarnes/stocker/alphavantage"
	log "github.com/sirupsen/logrus"
)

const (
	typeCurrency = "currency"
)

type Asset struct {
	AvgPrice    float64 `json:"avgPrice"`
	Fxr         float64 `json:"exchangeRate"`
	MarketValue float64 `json:"marketValue"`
	Alloc       float64 `json:"allocation"`
	Qty         float64 `json:"quantity"`
	Type        string  `json:"type"`
}

type AssetRebalance struct {
	Source map[string]Asset `json:"source"`
	Target map[string]Asset `json:"target"`
}

type Portfolio struct {
	apiKey   string           `json:"-"`
	Assets   AssetRebalance   `json:"assets"`
	currency string           `json:"-"`
}

func (p *Portfolio) allocate() error {
	var err error
	allocation := float64(0)

	// if err == nil and check that cash is greater than zero

	for _, v := range p.Assets.Target {
		allocation += v.Alloc
	}

	cash, _ := p.Assets.Target[p.currency]
	cashLeft := cash.Qty

	if allocation == 100. {
		for symbol, asset := range p.Assets.Target {
			if asset.Alloc > 0. && symbol != p.currency {
				if asset.AvgPrice <= 0 {
					err = p.initializeAsset(symbol, &asset)
				}

				if err == nil && asset.AvgPrice > 0 {
					asset.Qty = math.Floor(cash.Qty * asset.Alloc / 100. / asset.AvgPrice)
					asset.MarketValue = asset.Qty * asset.AvgPrice
					asset.Alloc = asset.MarketValue * 100. / cash.Qty
					cashLeft -= asset.MarketValue
					p.Assets.Target[symbol] = asset
				}
			}
		}

		cash.MarketValue = cashLeft
		cash.Alloc = cash.MarketValue * 100. / cash.Qty
		cash.Qty = cashLeft
		p.Assets.Target[p.currency] = cash
		log.Info("target portfolio:", getPrettyString(p.Assets.Target))
	} else {
		err = errors.New("Portfolio asset allocations do not add up to 100%")
	}

	return err
}

func (p *Portfolio) initializeAsset(symbol string, asset *Asset) error {
	var err error
	var search *av.SymbolSearchMatch

	log.Debug(symbol, ": searching for symbol information")

	asset.Type = strings.ToLower(asset.Type)
	if asset.Type == typeCurrency {
		search = &av.SymbolSearchMatch{Currency: symbol}
		asset.AvgPrice = 1
		asset.Type = typeCurrency
	} else if search, err = av.GetSymbolSearch(symbol, p.apiKey); err == nil {
		log.Debug(symbol, ": searching for symbol quote information")
		asset.AvgPrice, err = av.GetSymbolQuotePrice(symbol, p.apiKey)
		asset.Type = search.Type
	}

	if err == nil {
		if search.Currency == p.currency {
			asset.Fxr = 1.
		} else {
			log.Debug(symbol, ": searching for exchange rate from ", search.Currency, " to ", p.currency)
			asset.Fxr, err = av.GetCurrencyExchangeRate(search.Currency, p.currency, p.apiKey)
		}
	}

	return err
}

func (p *Portfolio) liquidate() error {
	var err error

	// Create new target cash asset if necessary 
	targetCurrency, ok := p.Assets.Target[p.currency]
	if ok {
		targetCurrency.AvgPrice = 1.
		targetCurrency.Fxr = 1.
	} else {
		targetCurrency = Asset{Alloc: 0, AvgPrice: 1., Fxr: 1., Type: typeCurrency}
	}

	// Liquidate source assets
	for symbol, asset := range p.Assets.Source {
		if err = p.initializeAsset(symbol, &asset); err == nil {
			// Copy overlapping stock info from source to target
			if tmp, ok := p.Assets.Target[symbol]; ok {
				tmp.Fxr = asset.Fxr
				tmp.AvgPrice = asset.AvgPrice
				p.Assets.Target[symbol] = tmp
			}

			targetCurrency.Qty += asset.AvgPrice * asset.Fxr * float64(asset.Qty)

			asset.MarketValue = math.Floor(asset.Qty * asset.AvgPrice * asset.Fxr)
			p.Assets.Source[symbol] = asset
		} else {
			log.Fatal("Error liquidating source assets:", err)
			break
		}
	}

	// Calculate source asset allocation
	for symbol, asset := range p.Assets.Source {
		if asset.Qty > 0 {
			asset.Alloc = asset.MarketValue * 100. / targetCurrency.Qty
			p.Assets.Source[symbol] = asset
		}
	}

	if targetCurrency.Qty >= 0 {
		p.Assets.Target[p.currency] = targetCurrency
	} else {
		log.Fatal("Source assets cannot be liquidated")
	}

	log.Info("source portfolio:", getPrettyString(p.Assets.Source))

	return err
}

func NewPortfolio(filename, apiKey, currency string) (*Portfolio, error) {
	portfolio := Portfolio{currency: strings.ToUpper(currency)}
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		if err = json.Unmarshal([]byte(file), &portfolio); err == nil {
			portfolio.apiKey = apiKey
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return &portfolio, err
}

func (p *Portfolio) Rebalance() error {
	log.Info("Liquidating source assets into ", p.currency, " funds")
	err := p.liquidate()
	if err == nil {
		log.Info("Re-allocating source assets to match target allocations")
		err = p.allocate()
	} else {
		log.Fatal(err)
	}

	return err
}
