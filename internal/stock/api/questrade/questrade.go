package questrade

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"syscall"

	fp "github.com/robaho/fixed"
	"github.com/shanebarnes/stocker/internal/stock"
	"github.com/shanebarnes/stocker/internal/stock/api"
)

const (
	qtDomain = "questrade.com"
)

type qt struct {
	apiKey     string
	apiServer  string
	cache     *stock.Cache
	creds     api.OAuthCredentials
}

type redeemTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ApiServer    string `json:"api_server"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
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

func getServerHostname(server string) (string, error) {
	hostname := ""
	u, err := url.Parse(server)
	if err == nil {
		hostname = u.Hostname()
	}
	return hostname, err
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
	return strings.HasSuffix(apiServer, qtDomain)
}

func NewApiQuestrade(apiKey, apiServer string, creds api.OAuthCredentials) api.StockApi {
	if len(creds.AccessToken) > 0 && len(creds.ApiServer) > 0 {
		apiKey = creds.AccessToken
		apiServer, _ = getServerHostname(creds.ApiServer)
	}

	return &qt{
		apiKey: apiKey,
		apiServer: apiServer,
		cache: stock.NewCache(),
		creds: creds,
	}
}

// References:
//   https://www.questrade.com/api/documentation/getting-started
//   https://www.questrade.com/api/documentation/security
func (q *qt) RefreshCredentials() (*api.OAuthCredentials, error) {
	var creds *api.OAuthCredentials
	body, err := ApiGetResponseBody("https://login.questrade.com/oauth2/token?grant_type=refresh_token&refresh_token=" + q.creds.RefreshToken, "")
	if err == nil {
		response := redeemTokenResponse{}
		if err = json.Unmarshal(body, &response); err == nil {
			q.creds = api.OAuthCredentials{
				AccessToken: response.AccessToken,
				ApiServer: response.ApiServer,
				ExpiresIn: response.ExpiresIn,
				RefreshToken: response.RefreshToken,
				TokenType: response.TokenType,
			}
			creds = &q.creds

			q.apiKey = creds.AccessToken
			q.apiServer, err = getServerHostname(creds.ApiServer)
		}
	}

	return creds, err
}
