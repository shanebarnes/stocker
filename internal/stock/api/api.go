package api

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/shanebarnes/stocker/internal/stock"
	log "github.com/sirupsen/logrus"
)

const (
	ApiKeyEnvName    = "STOCKER_API_KEY"
	ApiServerEnvName = "STOCKER_API_SERVER"

	DefaultClientTimeout = time.Second * 4

	DefaultRequestBackoffDelay = time.Millisecond * 125
	DefaultRequestBackoffLimit = time.Second
	DefaultRequestRetryLimit   = 10
)

var Client *http.Client = &http.Client{
	Timeout: DefaultClientTimeout,
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, network, address)
		},
	},
}

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

func GetApiResponseBody(url, accessToken string, isRetryable func(*http.Response) bool) ([]byte, error) {
	var body []byte

	//fmt.Println("Making request to: ", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		if len(accessToken) > 0 {
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}

		MakeApiRequestWithRetry(Client, req, func(res *http.Response, rerr error) bool {
			retry := false
			err = rerr
			if err == nil {
				body, err = ioutil.ReadAll(res.Body)
				if res.StatusCode != http.StatusOK {
					err = errors.New(fmt.Sprintf("API response status code: %d, details: %s", res.StatusCode, string(body)))
					retry = (isRetryable != nil && isRetryable(res))
				}
			} else {
				retry = IsErrorRetryable(err)
			}
			return retry
		})
	}

	return body, err
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
		if buf, err := httputil.DumpRequest(req, true); err == nil {
			log.Debug(string(buf))
		}

		res, err := client.Do(req)
		if err == nil {
			if buf, err := httputil.DumpResponse(res, true); err == nil {
				log.Debug(string(buf))
			}
		}

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
