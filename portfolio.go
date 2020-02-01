package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Fxr         fp.Fixed
	MarketValue fp.Fixed
	Price       fp.Fixed
	Qty         fp.Fixed
}

type Asset struct {
	Alloc       string  `json:"allocation"`
	Currency    string  `json:"currency"`
	fp          fpAsset
	Fxr         string  `json:"exchangeRate"`
	MarketValue string  `json:"marketValue"`
	Name        string  `json:"name"`
	Price       string  `json:"price"`
	Qty         string  `json:"quantity"`
	Type        string  `json:"type"`
}

type AssetGroup map[string]Asset

type AssetRebalance struct {
	Source AssetGroup `json:"source"`
	Target AssetGroup `json:"target"`
}

type FxrCache    map[string]map[string]fp.Fixed
type QuoteCache  map[string]*av.SymbolQuote
type SymbolCache map[string]*av.SymbolSearchMatch

type Portfolio struct {
	apiKey      string
	Assets      AssetRebalance `json:"assets"`
	currency    string
	fxrCache    FxrCache
	quoteCache  QuoteCache
	symbolCache SymbolCache
}

func (p *Portfolio) allocate(funds fp.Fixed) error {
	var err error
	allocation := fp.NewF(0)

	log.Info("Re-allocating source assets to match target allocations")
	// Create new target cash asset if it does not already exist
	cash, exists := p.Assets.Target[p.currency]
	if exists {
		cash.fp.Price = fp.NewF(1)
		cash.fp.Fxr = fp.NewF(1)
	} else {
		cash = Asset{
			fp: fpAsset{Alloc: fp.NewF(0), Price: fp.NewF(1), Fxr: fp.NewF(1)},
			Type: typeCurrency,
		}
	}

	cash.Currency = p.currency
	cash.Name = p.currency
	cash.fp.Qty = funds
	p.Assets.Target[p.currency] = cash

	for _, v := range p.Assets.Target {
		allocation = allocation.Add(v.fp.Alloc)
	}

	cashLeft := cash.fp.Qty

	if allocation.Equal(fp.NewF(100)) {
		for symbol, asset := range p.Assets.Target {
			if symbol != p.currency && asset.fp.Alloc.GreaterThan(fp.NewF(0)) {
				if asset.fp.Price.LessThanOrEqual(fp.NewF(0)) {
					err = p.initializeAsset(symbol, &asset)
				}

				if err == nil && asset.fp.Price.GreaterThan(fp.NewF(0)) {
					// asset.Qty = math.Floor(cash.Qty * asset.Alloc / 100. / (asset.Price * asset.Fxr))
					qty := cash.fp.Qty.Mul(asset.fp.Alloc)
					qty = qty.Div(fp.NewF(100))
					qty = qty.Div(asset.fp.Price.Mul(asset.fp.Fxr))
					asset.fp.Qty = fp.NewI(qty.Int(), 0)

					// asset.MarketValue = math.Round(asset.Qty * asset.Price * asset.Fxr)
					asset.fp.MarketValue = asset.fp.Qty.Mul(asset.fp.Price.Mul(asset.fp.Fxr))
					asset.fp.MarketValue = asset.fp.MarketValue.Round(2)

					// asset.Alloc = math.Round(asset.MarketValue * 100. / cash.Qty)
					asset.fp.Alloc = asset.fp.MarketValue.Mul(fp.NewF(100))
					asset.fp.Alloc = asset.fp.Alloc.Div(cash.fp.Qty)
					asset.fp.Alloc = asset.fp.Alloc.Round(4)

					// cashLeft = cashLeft - asset.MarketValue
					cashLeft = cashLeft.Sub(asset.fp.MarketValue)
				}
			}

			p.Assets.Target[symbol] = asset
		}

		cash.fp.MarketValue = cashLeft
		// cash.Alloc = math.Round(cash.MarketValue * 100. / cash.Qty)
		cash.fp.Alloc = cash.fp.MarketValue.Mul(fp.NewF(100))
		cash.fp.Alloc = cash.fp.Alloc.Div(cash.fp.Qty)
		cash.fp.Alloc = cash.fp.Alloc.Round(4)
		cash.fp.Qty = cashLeft
		p.Assets.Target[p.currency] = cash
		p.copyAssetFixedToStrings(&p.Assets.Target)
		log.Info("target portfolio:", getPrettyString(p.Assets.Target))
	} else {
		err = fmt.Errorf("Invalid portfolio allocation total: %s", allocation.StringN(2))
	}

	return err
}

func (p *Portfolio) copyAssetFixedToStrings(group *AssetGroup) {
	for i, asset := range *group {
		asset.Alloc       = asset.fp.Alloc.StringN(4) + "%"
		asset.Fxr         = asset.fp.Fxr.StringN(4)
		asset.MarketValue = asset.fp.MarketValue.StringN(2) + p.currency
		asset.Price       = asset.fp.Price.StringN(2)
		asset.Qty         = asset.fp.Qty.StringN(2)
		(*group)[i]       = asset
	}
}

func (p *Portfolio) copyAssetStringsToFixed(group *AssetGroup) {
	for i, asset := range *group {
		asset.fp.Alloc       = newFixedFromString("alloc", asset.Alloc)
		asset.fp.Fxr         = newFixedFromString("fxr", asset.Fxr)
		asset.fp.MarketValue = newFixedFromString("mvp", asset.MarketValue)
		asset.fp.Price       = newFixedFromString("price", asset.Price)
		asset.fp.Qty         = newFixedFromString("qty", asset.Qty)
		(*group)[i]          = asset
	}
}

func (p *Portfolio) getExchangeRate(fromCcy string) (fp.Fixed, error) {
	var fxr fp.Fixed
	var err error

	toCcyCache, exists := p.fxrCache[fromCcy]
	if exists {
		fxr, exists = toCcyCache[p.currency]
		log.Debug(fromCcy, ": found cached exchange rate to ", p.currency, ": ", fxr.StringN(4))
	} else {
		p.fxrCache[fromCcy] = map[string]fp.Fixed{}
	}

	if !exists {
		var f float64
		if f, err = av.GetCurrencyExchangeRate(fromCcy, p.currency, p.apiKey); err == nil {
			fxr = fp.NewF(f)
			p.fxrCache[fromCcy][p.currency] = fxr
		}
	}

	return fxr, err
}

func (p *Portfolio) getSymbolQuotePrice(symbol string) (fp.Fixed, error) {
	var price fp.Fixed
	var quote *av.SymbolQuote
	var exists bool
	var err error

	if quote, exists = p.quoteCache[symbol]; exists {
		log.Debug(symbol, ": found cached symbol quote")
	} else {
		quote, err = av.GetSymbolQuote(symbol, p.apiKey)
	}

	if err == nil {
		if price, err = fp.NewSErr(quote.Price); err == nil && !exists {
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
		asset.fp.Price = fp.NewF(1)
		asset.Currency = symbol
		asset.Name = symbol
		asset.Type = typeCurrency
	} else if search, err = p.getSymbolSearch(symbol); err == nil {
		log.Debug(symbol, ": searching for symbol quote information")
		asset.fp.Price, err = p.getSymbolQuotePrice(symbol)
		asset.Currency = search.Currency
		asset.Name = search.Name
		asset.Type = search.Type
	}

	if err == nil {
		if search.Currency == p.currency {
			asset.fp.Fxr = fp.NewF(1)
		} else {
			log.Debug(symbol, ": searching for exchange rate from ", search.Currency, " to ", p.currency)
			asset.fp.Fxr, err = p.getExchangeRate(search.Currency)
		}
	}

	return err
}

func (p *Portfolio) liquidate() (fp.Fixed, error) {
	var err error
	var cash fp.Fixed

	log.Info("Liquidating source assets into ", p.currency, " funds")
	for symbol, asset := range p.Assets.Source {
		if err = p.initializeAsset(symbol, &asset); err == nil {
			// cash = cash + asset.Price * asset.Fxr * float64(asset.Qty)
			mvp := asset.fp.Price.Mul(asset.fp.Qty)
			mvp = mvp.Mul(asset.fp.Fxr)
			cash = cash.Add(mvp)

			asset.fp.MarketValue = mvp
			p.Assets.Source[symbol] = asset
		} else {
			log.Fatal("Error liquidating source assets:", err)
			break
		}
	}

	// Calculate source asset allocation
	for symbol, asset := range p.Assets.Source {
		if asset.fp.Qty.GreaterThan(fp.NewF(0)) && cash.GreaterThan(fp.NewF(0)) {
			// asset.Alloc = asset.MarketValue * 100. / cash
			alloc := asset.fp.MarketValue.Mul(fp.NewF(100))
			alloc = alloc.Div(cash)
			asset.fp.Alloc = alloc
			p.Assets.Source[symbol] = asset
		}
	}

	if cash.LessThan(fp.NewF(0)) {
		log.Fatal("Source assets cannot be liquidated")
	}

	p.copyAssetFixedToStrings(&p.Assets.Source)
	log.Info("source portfolio:", getPrettyString(p.Assets.Source))

	return cash, err
}

func newFixedFromString(key, val string) fp.Fixed {
	if len(val) == 0 {
		return fp.NewF(0)
	}
	fp, err := fp.NewSErr(val)
	if err != nil {
		log.Fatal(key, ": invalid value ", val)
	}

	return fp
}

func NewPortfolio(filename, apiKey, currency string) (*Portfolio, error) {
	portfolio := Portfolio{
		currency: strings.ToUpper(currency),
		fxrCache: map[string]map[string]fp.Fixed{},
		quoteCache: map[string]*av.SymbolQuote{},
		symbolCache: map[string]*av.SymbolSearchMatch{},
	}
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		if err = json.Unmarshal([]byte(file), &portfolio); err == nil {
			portfolio.apiKey = apiKey
			portfolio.copyAssetStringsToFixed(&portfolio.Assets.Source)
			portfolio.copyAssetStringsToFixed(&portfolio.Assets.Target)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return &portfolio, err
}

func (p *Portfolio) Rebalance() error {
	var err error
	var cash fp.Fixed

	if err = p.validate(); err != nil {
		log.Fatal("Validation failed:", err)
	} else if cash, err = p.liquidate(); err != nil {
		log.Fatal("Liquidation failed:", err)
	} else if err = p.allocate(cash); err != nil {
		log.Fatal("Allocation failed:", err)
	}

	return err
}

func (p *Portfolio) validate() error {
	var err error

	log.Info("Validating source assets")
	for symbol, asset := range p.Assets.Source {
		if err = p.initializeAsset(symbol, &asset); err != nil {
			break
		}
	}

	if err == nil {
		log.Info("Validating target assets")
		for symbol, asset := range p.Assets.Target {
			if err = p.initializeAsset(symbol, &asset); err != nil {
				break
			}
		}
	}

	return err
}
