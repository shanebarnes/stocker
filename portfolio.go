package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	av "github.com/shanebarnes/stocker/alphavantage"
	fp "github.com/robaho/fixed"
	log "github.com/sirupsen/logrus"
)

const (
	typeCurrency = "currency"
)

// Used for internal fixed point representation of assets
type fpAsset struct {
	Alloc       fp.Fixed
	Currency    string
	Fxr         fp.Fixed
	MarketValue fp.Fixed
	Name        string
	Price       fp.Fixed
	Qty         fp.Fixed
	Type        string
}

type Asset struct {
	Alloc       float64  `json:"allocation"`
	Currency    string   `json:"currency"`
	Fxr         float64  `json:"exchangeRate"`
	MarketValue float64  `json:"marketValue"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Qty         float64  `json:"quantity"`
	Type        string   `json:"type"`
}

type AssetRebalance struct {
	Source map[string]Asset `json:"source"`
	Target map[string]Asset `json:"target"`
}

type FxrCache    map[string]map[string]float64
type QuoteCache  map[string]*av.SymbolQuote
type SymbolCache map[string]*av.SymbolSearchMatch

type Portfolio struct {
	apiKey      string         `json:"-"`
	Assets      AssetRebalance `json:"assets"`
	currency    string         `json:"-"`
	fxrCache    FxrCache       `json:"-"`
	quoteCache  QuoteCache     `json:"-"`
	symbolCache SymbolCache    `json:"-"`
}

func (p *Portfolio) allocate(funds float64) error {
	var err error
	allocation := float64(0)

	// Create new target cash asset if necessary
	cash, exists := p.Assets.Target[p.currency]
	if exists {
		cash.Price = 1.
		cash.Fxr = 1.
	} else {
		cash = Asset{Alloc: 0, Price: 1., Fxr: 1., Type: typeCurrency}
	}

	cash.Currency = p.currency
	cash.Name = p.currency
	cash.Qty = funds
	p.Assets.Target[p.currency] = cash

	for _, v := range p.Assets.Target {
		allocation += v.Alloc
	}

	cashLeft := cash.Qty

	if allocation == 100. {
		for symbol, asset := range p.Assets.Target {
			if symbol != p.currency && asset.Alloc > 0. {
				if asset.Price <= 0 {
					err = p.initializeAsset(symbol, &asset)
				}

				if err == nil && asset.Price > 0 {
					asset.Qty = math.Floor(cash.Qty * asset.Alloc / 100. / asset.Price)
					asset.MarketValue = asset.Qty * asset.Price
					asset.Alloc = asset.MarketValue * 100. / cash.Qty
					cashLeft -= asset.MarketValue
				}
			}

			p.Assets.Target[symbol] = asset
		}

		cash.MarketValue = cashLeft
		cash.Alloc = cash.MarketValue * 100. / cash.Qty
		cash.Qty = cashLeft
		p.Assets.Target[p.currency] = cash

		log.Info("target portfolio:", getPrettyString(p.Assets.Target))
	} else {
		err = fmt.Errorf("Invalid portfolio allocation total: %f", allocation)
	}

	return err
}

func (p *Portfolio) getExchangeRate(fromCcy string) (float64, error) {
	var fxr float64
	var err error

	toCcyCache, exists := p.fxrCache[fromCcy]
	if exists {
		fxr, exists = toCcyCache[p.currency]
		log.Debug(fromCcy, ": found cached exchange rate to ", p.currency, ": ", fxr)
	} else {
		p.fxrCache[fromCcy] = map[string]float64{}
	}

	if !exists {
		if fxr, err = av.GetCurrencyExchangeRate(fromCcy, p.currency, p.apiKey); err == nil {
			p.fxrCache[fromCcy][p.currency] = fxr
		}
	}

	return fxr, err
}

func (p *Portfolio) getSymbolQuotePrice(symbol string) (float64, error) {
	var price float64
	var quote *av.SymbolQuote
	var exists bool
	var err error

	if quote, exists = p.quoteCache[symbol]; exists {
		log.Debug(symbol, ": found cached symbol quote")
	} else {
		quote, err = av.GetSymbolQuote(symbol, p.apiKey)
	}

	if err == nil {
		if price, err = strconv.ParseFloat(quote.Price, 64); err == nil && !exists {
			p.quoteCache[symbol] = quote
		}
	}

	return price, err
}

func (p *Portfolio) getSymbolSearch(symbol string) (*av.SymbolSearchMatch, error) {
	var match *av.SymbolSearchMatch
	var exists bool
	var err error

	if match, exists = p.symbolCache[symbol]; exists {
		log.Debug(symbol, ": found cached symbol information")
	} else {
		if match, err = av.GetSymbolSearch(symbol, p.apiKey); err == nil {
			p.symbolCache[symbol] = match
		}
	}

	return match, err
}

func (p *Portfolio) initializeAsset(symbol string, asset *Asset) error {
	var err error
	var search *av.SymbolSearchMatch

	log.Debug(symbol, ": searching for symbol information")

	asset.Type = strings.ToLower(asset.Type)
	if asset.Type == typeCurrency {
		search = &av.SymbolSearchMatch{Currency: symbol}
		asset.Price = 1
		asset.Currency = symbol
		asset.Name = symbol
		asset.Type = typeCurrency
	} else if search, err = p.getSymbolSearch(symbol); err == nil {
		log.Debug(symbol, ": searching for symbol quote information")
		asset.Price, err = p.getSymbolQuotePrice(symbol)
		asset.Currency = search.Currency
		asset.Name = search.Name
		asset.Type = search.Type
	}

	if err == nil {
		if search.Currency == p.currency {
			asset.Fxr = 1.
		} else {
			log.Debug(symbol, ": searching for exchange rate from ", search.Currency, " to ", p.currency)
			asset.Fxr, err = p.getExchangeRate(search.Currency)
		}
	}

	return err
}

func (p *Portfolio) liquidate() (float64, error) {
	var err error
	var cash float64

	for symbol, asset := range p.Assets.Source {
		if err = p.initializeAsset(symbol, &asset); err == nil {
			cash += asset.Price * asset.Fxr * float64(asset.Qty)

			asset.MarketValue = asset.Qty * asset.Price * asset.Fxr
			p.Assets.Source[symbol] = asset
		} else {
			log.Fatal("Error liquidating source assets:", err)
			break
		}
	}

	// Calculate source asset allocation
	for symbol, asset := range p.Assets.Source {
		if asset.Qty > 0 && cash > 0 {
			asset.Alloc = asset.MarketValue * 100. / cash
			p.Assets.Source[symbol] = asset
		}
	}

	if cash < 0 {
		log.Fatal("Source assets cannot be liquidated")
	}

	log.Info("source portfolio:", getPrettyString(p.Assets.Source))

	return cash, err
}

func NewPortfolio(filename, apiKey, currency string) (*Portfolio, error) {
	portfolio := Portfolio{
		currency: strings.ToUpper(currency),
		fxrCache: map[string]map[string]float64{},
		quoteCache: map[string]*av.SymbolQuote{},
		symbolCache: map[string]*av.SymbolSearchMatch{},
	}
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
	cash, err := p.liquidate()
	if err == nil {
		log.Info("Re-allocating source assets to match target allocations")
		err = p.allocate(cash)
	}

	if err != nil {
		log.Fatal(err)
	}

	return err
}
