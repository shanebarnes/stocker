package portfolio

import (
	"encoding/json"
	"fmt"
	fp "github.com/robaho/fixed"
	"github.com/shanebarnes/stocker/internal/stock"
	"github.com/shanebarnes/stocker/internal/stock/api"
	av "github.com/shanebarnes/stocker/internal/stock/api/alphavantage"
	qt "github.com/shanebarnes/stocker/internal/stock/api/questrade"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
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
	PriceDiff   fp.Fixed
	Qty         fp.Fixed
	QtyDiff     fp.Fixed
}

type order struct {
	MarketValue string `json:"marketValue"`
	Qty         string `json:"quantity"`
}

type Asset struct {
	Alloc       string  `json:"allocation"`
	Currency    string  `json:"currency"`
	fp          fpAsset
	Fxr         string  `json:"exchangeRate"`
	MarketValue string  `json:"marketValue"`
	Name        string  `json:"name"`
	Order       *order  `json:"order,omitempty"`
	Price       string  `json:"price"`
	Qty         string  `json:"quantity"`
	Type        string  `json:"type"`
}

// TODO: Convert to struct and include total market value field
type AssetGroup map[string]Asset

type AssetRebalance struct {
	Source AssetGroup `json:"source"`
	Target AssetGroup `json:"target"`
}

type Portfolio struct {
	Api      api.StockApi
	Assets   AssetRebalance `json:"assets"`
	currency string
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

					// Currency quantities do not need to be integers
					if asset.Type == typeCurrency {
						asset.fp.Qty = qty
					} else {
						asset.fp.Qty = fp.NewI(qty.Int(), 0)
					}

					// asset.MarketValue = math.Round(asset.Qty * asset.Price * asset.Fxr)
					asset.fp.MarketValue = asset.fp.Qty.Mul(asset.fp.Price.Mul(asset.fp.Fxr))

					// asset.Alloc = math.Round(asset.MarketValue * 100. / cash.Qty)
					asset.fp.Alloc = asset.fp.MarketValue.Mul(fp.NewF(100))
					asset.fp.Alloc = asset.fp.Alloc.Div(cash.fp.Qty)

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
		cash.fp.Qty = cashLeft
		p.Assets.Target[p.currency] = cash
		p.diffAssets(&p.Assets.Source, &p.Assets.Target)
		p.copyAssetFixedToStrings(&p.Assets.Target)
		log.Info("target portfolio:", getPrettyString(p.Assets.Target))
	} else {
		err = fmt.Errorf("Invalid portfolio allocation total: %s", allocation.Round(2).StringN(2))
	}

	return err
}

func (p *Portfolio) copyAssetFixedToStrings(group *AssetGroup) {
	for symbol, asset := range *group {
		asset.Alloc       = asset.fp.Alloc.Round(4).StringN(4) + "%"
		asset.Fxr         = asset.fp.Fxr.Round(4).StringN(4)
		asset.MarketValue = asset.fp.MarketValue.Round(2).StringN(2) + p.currency
		asset.Price       = asset.fp.Price.Round(2).StringN(2)
		asset.Qty         = asset.fp.Qty.Round(2).StringN(2)
		(*group)[symbol]  = asset
	}
}

// Find order quantities (currency/share buys/sells)
func (p *Portfolio) diffAssets(source, target *AssetGroup) {
	// Add any symbols in source assets that are not found in target assets
	for symbol, srcAsset := range *source {
		if _, ok := (*target)[symbol]; !ok {
			srcAsset.fp.QtyDiff = srcAsset.fp.Qty.Mul(fp.NewF(-1))
			srcAsset.fp.Alloc = fp.NewF(0)
			srcAsset.fp.MarketValue = fp.NewF(0)
			srcAsset.fp.Qty = fp.NewF(0)
			(*target)[symbol] = srcAsset
		}
	}

	// Find difference between symbols common to source and target assets
	for symbol, tgtAsset := range *target {
		if srcAsset, ok := (*source)[symbol]; ok {
			tgtAsset.fp.QtyDiff = tgtAsset.fp.Qty.Sub(srcAsset.fp.Qty)
		} else {
			tgtAsset.fp.QtyDiff = tgtAsset.fp.Qty
		}
		tgtAsset.fp.PriceDiff = tgtAsset.fp.QtyDiff.Mul(tgtAsset.fp.Price).Mul(tgtAsset.fp.Fxr)

		sign := ""
		if tgtAsset.fp.QtyDiff.Sign() != -1 {
			sign = "+"
		}

		tgtAsset.Order = &order{}
		tgtAsset.Order.MarketValue = sign + tgtAsset.fp.PriceDiff.Round(2).StringN(2) + p.currency
		tgtAsset.Order.Qty = sign + tgtAsset.fp.QtyDiff.Round(2).StringN(2)
		(*target)[symbol] = tgtAsset
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

func getPrettyString(v interface{}) string {
	str := ""
	buf, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		str = string(buf)
	}
	return str
}

func (p *Portfolio) initializeAsset(symbol string, asset *Asset) error {
	var err error
	var search stock.Symbol

	log.Debug(symbol, ": searching for symbol information")

	asset.Type = strings.ToLower(asset.Type)
	if asset.Type == typeCurrency {
		search = stock.Symbol{Currency: symbol}
		asset.fp.Price = fp.NewF(1)
		asset.Currency = symbol
		asset.Name = symbol
		asset.Type = typeCurrency
		log.Debug(symbol, ": ", asset)
	} else if search, err = p.Api.GetSymbol(symbol); err == nil {
		log.Debug(symbol, ": searching for symbol quote information")
		var quote stock.Quote
		if quote, err = p.Api.GetQuote(symbol); err == nil {
			asset.fp.Price = fp.NewF(quote.Prices.Latest)
		}

		asset.Currency = search.Currency
		asset.Name = search.Description
		asset.Type = search.Type
		log.Debug(symbol, ": ", asset)
	}

	if err == nil {
		if search.Currency == p.currency {
			asset.fp.Fxr = fp.NewF(1)
		} else {
			log.Debug(symbol, ": searching for exchange rate from ", search.Currency, " to ", p.currency)
			var ccy stock.Currency
			ccy, err = p.Api.GetCurrency(search.Currency, p.currency)
			if err == nil {
				asset.fp.Fxr = ccy.Rates[p.currency]
				//asset.fp.Fxr, err = p.getExchangeRate(search.Currency)
			}
		}
	} else {
		log.Debug("Error getting symbol ", symbol, ": ", err)
		// TODO: Wrap error?
		//err = fmt.Errorf("%v: %w", err, syscall.EACCES)
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

func NewPortfolio(filename, apiKey, apiServer, currency string) (*Portfolio, error) {
	var api api.StockApi
	if av.IsApiAlphavantage(apiServer) {
		api = av.NewApiAlphavantage(apiKey)
	} else if qt.IsApiQuestrade(apiServer) {
		api = qt.NewApiQuestrade(apiKey, apiServer)
	} else {
		log.Fatal("Invalid API server: ", apiServer)
	}

	portfolio := Portfolio{
		Api: api,
		currency: strings.ToUpper(currency),
	}
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		if err = json.Unmarshal([]byte(file), &portfolio); err == nil {
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
		log.Fatal("Validation failed: ", err)
	} else if cash, err = p.liquidate(); err != nil {
		log.Fatal("Liquidation failed: ", err)
	} else if err = p.allocate(cash); err != nil {
		log.Fatal("Allocation failed: ", err)
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
