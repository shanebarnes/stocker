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
	ApiAuthTokenEnvName = "STOCKER_API_AUTH_TOKEN"

	DefaultClientTimeout = time.Second * 4

	DefaultRequestBackoffDelay = time.Millisecond * 125
	DefaultRequestBackoffLimit = time.Second
	DefaultRequestRetryLimit = 10
)

type AuthResponse struct {
	AccessToken string
	ApiServer   string
}

type StockApi interface {
	GetCurrency(currency, currencyTo string) (stock.Currency, error)
	GetQuote(symbol string) (stock.Quote, error)
	GetSymbol(symbol string) (stock.Symbol, error)
	RedeemAuthToken(token string) (*AuthResponse, error)
}

func GetApiKeyFromEnv() string {
	return os.Getenv(ApiKeyEnvName)
}

func GetApiServerFromEnv() string {
	return os.Getenv(ApiServerEnvName)
}

func GetApiAuthTokenFromEnv() string {
	return os.Getenv(ApiAuthTokenEnvName)
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

func RedeemApiAuthToken(authToken, apiServer string) (string, string, error) {
	return "", "", nil
}
