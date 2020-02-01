# stocker
Rebalance your finanical assets with realtime stock market data acquired on demand from [Alpha Vantage](https://www.alphavantage.co/) APIs.

## Build Instructions

```shell
$ git clone https://github.com/shanebarnes/stocker.git
$ cd stocker
$ go vet -v ./...  # Optional
$ export AV_API_KEY=<your_api_key>; go test -v ./... # Optional, takes some time with free API Key due to API call limit
$ go build -v
```

## Examples

Try rebalancing a sample portfolio! An [Alpha Vantage](https://www.alphavantage.co/) API key is required. A free API key can be claimed [here](https://www.alphavantage.co/support/#api-key).

```shell
$ ./stocker -apiKey <your_api_key> -portfolio examples/portfolio.json
$ export AV_API_KEY=<your_api_key>; ./stocker -portfolio examples/portfolio.json # Alternatively, load API key from environment
```

The [portfolio.json](https://github.com/shanebarnes/stocker/blob/master/examples/portfolio.json) contains source and target assets.

```json
{
    "assets": {
        "source": {
            "AAPL": {
                "quantity": "10"
            },
            "CAD": {
                "type": "Currency",
                "quantity": "10000.00"
            },
            "USD": {
                "type": "Currency",
                "quantity": "10000.00"
            }
        },
        "target": {
            "CAD": {
                "type": "Currency",
                "allocation": "20.00"
            },
            "MSFT": {
                "allocation": "70.00"
            },
            "USD": {
                "type": "Currency",
                "allocation": "10.00"
            }
        }
    }
}
```

Here is what the asset rebalancing output will look like for the sample portfolio.

```shell
WARN[2019-08-26T22:38:46.647-04:00] Rebalancing requires making Alpha Vantage API calls
WARN[2019-08-26T22:38:46.648-04:00] Only 5 API calls to Alpha Vantage will be performed each minute
INFO[2019-08-26T22:38:46.648-04:00] Liquidating source assets into CAD funds
DEBU[2019-08-26T22:38:46.648-04:00] USD: searching for symbol information
DEBU[2019-08-26T22:38:46.648-04:00] USD: searching for exchange rate from USD to CAD
DEBU[2019-08-26T22:38:47.038-04:00] AAPL: searching for symbol information
DEBU[2019-08-26T22:39:01.851-04:00] AAPL: searching for symbol quote information
DEBU[2019-08-26T22:39:16.706-04:00] AAPL: searching for exchange rate from USD to CAD
DEBU[2019-08-26T22:39:16.706-04:00] USD: found cached exchange rate to CAD: 1.3239
DEBU[2019-08-26T22:39:16.706-04:00] CAD: searching for symbol information
INFO[2019-08-26T22:39:16.706-04:00] source portfolio:{
  "AAPL": {
    "allocation": "10.5253%",
    "currency": "USD",
    "exchangeRate": "1.3239",
    "marketValue": "2733.72CAD",
    "name": "Apple Inc.",
    "price": "206.49",
    "quantity": "10.00",
    "type": "Equity"
  },
  "CAD": {
    "allocation": "38.5019%",
    "currency": "CAD",
    "exchangeRate": "1.0000",
    "marketValue": "10000.00CAD",
    "name": "CAD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  },
  "USD": {
    "allocation": "50.9727%",
    "currency": "USD",
    "exchangeRate": "1.3239",
    "marketValue": "13239.00CAD",
    "name": "USD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  }
}
INFO[2019-08-26T22:39:16.706-04:00] Re-allocating source assets to match target allocations
DEBU[2019-08-26T22:39:16.706-04:00] MSFT: searching for symbol information
DEBU[2019-08-26T22:39:31.864-04:00] MSFT: searching for symbol quote information
DEBU[2019-08-26T22:39:46.754-04:00] MSFT: searching for exchange rate from USD to CAD
DEBU[2019-08-26T22:39:46.755-04:00] USD: found cached exchange rate to CAD: 1.3239
DEBU[2019-08-26T22:39:46.755-04:00] USD: searching for symbol information
DEBU[2019-08-26T22:39:46.755-04:00] USD: searching for exchange rate from USD to CAD
DEBU[2019-08-26T22:39:46.755-04:00] USD: found cached exchange rate to CAD: 1.3239
INFO[2019-08-26T22:39:46.755-04:00] target portfolio:{
  "CAD": {
    "allocation": "20.2713%",
    "currency": "CAD",
    "exchangeRate": "1.0000",
    "marketValue": "5265.00CAD",
    "name": "CAD",
    "price": "1.00",
    "quantity": "5265.00",
    "type": "Currency"
  },
  "MSFT": {
    "allocation": "69.7330%",
    "currency": "USD",
    "exchangeRate": "1.3239",
    "marketValue": "18111.55CAD",
    "name": "Microsoft Corporation",
    "price": "135.45",
    "quantity": "101.00",
    "type": "Equity"
  },
  "USD": {
    "allocation": "9.9958%",
    "currency": "USD",
    "exchangeRate": "1.3239",
    "marketValue": "2596.17CAD",
    "name": "USD",
    "price": "1.00",
    "quantity": "1961.00",
    "type": "currency"
  }
}
```

Please note that a free API key is limited to 5 API calls per minute.
