package questrade

import (
	"strconv"
	"strings"
	"syscall"

	fp "github.com/robaho/fixed"
	"github.com/shanebarnes/stocker/internal/stock"
	"github.com/shanebarnes/stocker/internal/stock/api"
)

type qt struct {
	apiKey     string
	apiServer  string
	cache     *stock.Cache
}

func (q *qt) GetCurrency(currency, currencyTo string) (stock.Currency, error) {
	ccy, err := q.cache.GetCurrency(currency, currencyTo)
	if err != nil {
		var xr *ExchangeRate
		if xr, err = GetCurrencyExchangeRateInfo(currency, currencyTo, q.apiKey); err == nil {
			if _, exists := xr.Rates[currencyTo]; exists {
				ccy.Currency = xr.BaseSymbol
				ccy.Name = xr.BaseSymbol
				ccy.Rates = make(map[string]fp.Fixed)
				for key, val := range xr.Rates {
					ccy.Rates[key] = fp.NewF(val)
				}

				q.cache.AddCurrency(ccy)
			} else {
				err = syscall.ENOENT
			}
		}
	}
	return ccy, err
}

func (q *qt) GetQuote(symbol string) (stock.Quote, error) {
	sym, _ := q.GetSymbol(symbol)
	qte, err := q.cache.GetQuote(sym.Id)
	if err != nil {
		var quote *SymbolQuote
		if quote, err = GetSymbolQuote(sym.Id, q.apiKey, q.apiServer); err == nil {
			qte.Symbol = quote.Symbol
			qte.Prices.Ask = quote.AskPrice
			qte.Prices.Bid = quote.BidPrice
			qte.Prices.High = quote.HighPrice
			qte.Prices.Low = quote.LowPrice
			qte.Prices.Open = quote.OpenPrice
			qte.Prices.Latest = quote.LastTradePrice
			qte.Prices.Latest = quote.LastTradePriceTrHrs
			qte.Prices.LatestTrHrs = quote.LastTradePriceTrHrs
			qte.Volume = strconv.FormatInt(int64(quote.Volume), 10)

			q.cache.AddQuote(qte)
		}
	}
	return qte, err
}

func (q *qt) GetSymbol(symbol string) (stock.Symbol, error) {
	sym, err := q.cache.GetSymbol(symbol)
	if err != nil {
		var match *SymbolSearchMatch
		if match, err = GetSymbolSearch(symbol, q.apiKey, q.apiServer); err == nil {
			sym.Currency = match.Currency
			sym.Description = match.Description
			sym.Id = strconv.FormatInt(int64(match.SymbolId), 10)
			sym.Symbol = match.Symbol
			sym.Type = match.SecurityType

			q.cache.AddSymbol(sym)
		}
	}
	return sym, err
}

func IsApiQuestrade(apiServer string) bool {
	return strings.HasSuffix(apiServer, "questrade.com")
}

func NewApiQuestrade(apiKey, apiServer string) api.StockApi {
	return &qt{
		apiKey: apiKey,
		apiServer: apiServer,
		cache: stock.NewCache(),
	}
}
