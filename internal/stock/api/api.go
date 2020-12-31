package api

import (
	"net/http"
	"os"
	"time"

	"github.com/shanebarnes/stocker/internal/stock"
)

const (
	ApiKeyEnvName = "STOCKER_API_KEY"
	ApiServerEnvName = "STOCKER_API_SERVER"

	DefaultClientTimeout = time.Second * 4

	DefaultRequestBackoffDelay = time.Millisecond * 125
	DefaultRequestBackoffLimit = time.Second
	DefaultRequestRetryLimit = 10
)

type OAuthCredentials struct {
	AccessToken  string
	ApiServer    string
	ExpiresIn    int
	RefreshToken string
	TokenType    string
}

type StockApi interface {
	GetCurrency(currency, currencyTo string) (stock.Currency, error)
	GetQuote(symbol string) (stock.Quote, error)
	GetSymbol(symbol string) (stock.Symbol, error)
	RefreshCredentials() (*OAuthCredentials, error)
}

func GetApiKeyFromEnv() string {
	return os.Getenv(ApiKeyEnvName)
}

func GetApiServerFromEnv() string {
	return os.Getenv(ApiServerEnvName)
}

func MakeApiRequestWithRetry(client *http.Client, req *http.Request, retryCb func(res *http.Response, err error) bool) {
	backoff := DefaultRequestBackoffDelay
	backoffLimit := DefaultRequestBackoffLimit
	retry := 0
	retryLimit := DefaultRequestRetryLimit

	if retryCb == nil {
		retry = retryLimit
	}

	for retry < retryLimit {
		res, err := client.Do(req)
		defer func() {
			if err == nil {
				res.Body.Close()
			}
		}()

		if retryCb(res, err) {
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > backoffLimit {
				backoff = backoffLimit
			}
			retry++
		} else {
			break
		}
	}
}
