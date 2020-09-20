package alphavantage

import (
	"strconv"
	"strings"

	fp "github.com/robaho/fixed"
	"github.com/shanebarnes/stocker/internal/stock"
	"github.com/shanebarnes/stocker/internal/stock/api"
)

type av struct {
	apiKey string
	cache  *stock.Cache
}

func (a *av) GetCurrency(currency, currencyTo string) (stock.Currency, error) {
	ccy, err := a.cache.GetCurrency(currency, currencyTo)
	if err != nil {
		var xr *ExchangeRate
		if xr, err = GetCurrencyExchangeRateInfo(currency, currencyTo, a.apiKey); err == nil {
			var rate float64
			if rate, err = strconv.ParseFloat(xr.ExchangeRate, 64); err == nil {
				ccy.Currency = xr.FromCode
				ccy.Name = xr.FromName
				ccy.Rates = make(map[string]fp.Fixed)
				ccy.Rates[currencyTo] = fp.NewF(rate)

				a.cache.AddCurrency(ccy)
			}
		}
	}
	return ccy, err
}

func (a *av) GetQuote(symbol string) (stock.Quote, error) {
	qte, err := a.cache.GetQuote(symbol)
	if err != nil {
		var quote *SymbolQuote
		if quote, err = GetSymbolQuote(symbol, a.apiKey); err == nil {
			qte.Symbol = quote.Symbol
			qte.Prices.Close, _ = strconv.ParseFloat(quote.PreviousClose, 64)
			qte.Prices.High, _ = strconv.ParseFloat(quote.High, 64)
			qte.Prices.Low, _ = strconv.ParseFloat(quote.Low, 64)
			qte.Prices.Open, _ = strconv.ParseFloat(quote.Open, 64)
			qte.Prices.Latest, _ = strconv.ParseFloat(quote.Price, 64)
			qte.Prices.LatestTrHrs, _ = strconv.ParseFloat(quote.LatestTradingDay, 64)
			qte.Volume = quote.Volume

			a.cache.AddQuote(qte)
		}
	}
	return qte, err
}

func (a *av) GetSymbol(symbol string) (stock.Symbol, error) {
	sym, err := a.cache.GetSymbol(symbol)
	if err != nil {
		var match *SymbolSearchMatch
		if match, err = GetSymbolSearch(symbol, a.apiKey); err == nil {
			sym.Currency = match.Currency
			sym.Description = match.Name
			sym.Symbol = match.Symbol
			sym.Type = match.Type

			a.cache.AddSymbol(sym)
		}
	}
	return sym, err
}

	func IsApiAlphavantage(apiServer string) bool {
		return strings.HasSuffix(apiServer, "alphavantage.co")
	}

func NewApiAlphavantage(apiKey string) api.StockApi {
	return &av{
		apiKey: apiKey,
		cache: stock.NewCache(),
	}
}
