# stocker
Rebalance your finanical assets with realtime stock market data acquired on demand from [Alpha Vantage](https://www.alphavantage.co/) APIs.

## Build Instructions

```shell
$ git clone https://github.com/shanebarnes/stocker.git
$ cd stocker
$
$ # Optional
$ go vet -v ./...
$
$ # Optional: unit tests could take a while with free API key due to API call limit
$ export AV_API_KEY=<your_api_key>; go test -cover -v ./...
$
$ cd cmd/stocker
$ go build -v
```

## Examples

Try rebalancing a sample portfolio! An [Alpha Vantage](https://www.alphavantage.co/) API key is required. A free API key can be claimed [here](https://www.alphavantage.co/support/#api-key).

```shell
$ # Pass API key on command line
$ ./stocker -apiKey <your_api_key> -rebalance ../../examples/portfolio.json
$
$ # Alternatively, load API key from environment
$ export AV_API_KEY=<your_api_key>; ./stocker -rebalance ../../examples/portfolio.json
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
WARN[2020-02-01T18:39:21.237-05:00] Rebalancing requires making Alpha Vantage API calls
WARN[2020-02-01T18:39:21.237-05:00] Only 5 API calls to Alpha Vantage will be performed each minute
INFO[2020-02-01T18:39:21.238-05:00] Validating source assets
DEBU[2020-02-01T18:39:21.238-05:00] AAPL: searching for symbol information
DEBU[2020-02-01T18:39:21.687-05:00] AAPL: searching for symbol quote information
DEBU[2020-02-01T18:39:36.300-05:00] CAD: searching for symbol information
DEBU[2020-02-01T18:39:36.300-05:00] CAD: searching for exchange rate from CAD to USD
DEBU[2020-02-01T18:39:51.294-05:00] USD: searching for symbol information
INFO[2020-02-01T18:39:51.294-05:00] Validating target assets
DEBU[2020-02-01T18:39:51.294-05:00] CAD: searching for symbol information
DEBU[2020-02-01T18:39:51.294-05:00] CAD: searching for exchange rate from CAD to USD
DEBU[2020-02-01T18:39:51.294-05:00] CAD: found cached exchange rate to USD: 0.7554
DEBU[2020-02-01T18:39:51.294-05:00] MSFT: searching for symbol information
DEBU[2020-02-01T18:40:06.376-05:00] MSFT: searching for symbol quote information
DEBU[2020-02-01T18:40:21.317-05:00] USD: searching for symbol information
INFO[2020-02-01T18:40:21.317-05:00] Liquidating source assets into USD funds
DEBU[2020-02-01T18:40:21.317-05:00] AAPL: searching for symbol information
DEBU[2020-02-01T18:40:21.317-05:00] AAPL: found cached symbol information
DEBU[2020-02-01T18:40:21.317-05:00] AAPL: searching for symbol quote information
DEBU[2020-02-01T18:40:21.317-05:00] AAPL: found cached symbol quote
DEBU[2020-02-01T18:40:21.317-05:00] CAD: searching for symbol information
DEBU[2020-02-01T18:40:21.317-05:00] CAD: searching for exchange rate from CAD to USD
DEBU[2020-02-01T18:40:21.317-05:00] CAD: found cached exchange rate to USD: 0.7554
DEBU[2020-02-01T18:40:21.317-05:00] USD: searching for symbol information
INFO[2020-02-01T18:40:21.317-05:00] source portfolio:{
  "AAPL": {
    "allocation": "14.9890%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "3095.10USD",
    "name": "Apple Inc.",
    "price": "309.51",
    "quantity": "10.00",
    "type": "Equity"
  },
  "CAD": {
    "allocation": "36.5827%",
    "currency": "CAD",
    "exchangeRate": "0.7554",
    "marketValue": "7554.00USD",
    "name": "CAD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  },
  "USD": {
    "allocation": "48.4282%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "10000.00USD",
    "name": "USD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  }
}
INFO[2020-02-01T18:40:21.317-05:00] Re-allocating source assets to match target allocations
DEBU[2020-02-01T18:40:21.317-05:00] CAD: searching for symbol information
DEBU[2020-02-01T18:40:21.317-05:00] CAD: searching for exchange rate from CAD to USD
DEBU[2020-02-01T18:40:21.318-05:00] CAD: found cached exchange rate to USD: 0.7554
DEBU[2020-02-01T18:40:21.318-05:00] MSFT: searching for symbol information
DEBU[2020-02-01T18:40:21.318-05:00] MSFT: found cached symbol information
DEBU[2020-02-01T18:40:21.318-05:00] MSFT: searching for symbol quote information
DEBU[2020-02-01T18:40:21.318-05:00] MSFT: found cached symbol quote
INFO[2020-02-01T18:40:21.318-05:00] target portfolio:{
  "AAPL": {
    "allocation": "0.0000%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "0.00USD",
    "name": "Apple Inc.",
    "order": {
      "marketValue": "-3095.10USD",
      "quantity": "-10.00"
    },
    "price": "309.51",
    "quantity": "0.00",
    "type": "Equity"
  },
  "CAD": {
    "allocation": "19.9998%",
    "currency": "CAD",
    "exchangeRate": "0.7554",
    "marketValue": "4129.77USD",
    "name": "CAD",
    "order": {
      "marketValue": "-3424.22USD",
      "quantity": "-4533.00"
    },
    "price": "1.00",
    "quantity": "5467.00",
    "type": "currency"
  },
  "MSFT": {
    "allocation": "70.0735%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "14469.55USD",
    "name": "Microsoft Corporation",
    "order": {
      "marketValue": "+14469.55USD",
      "quantity": "+85.00"
    },
    "price": "170.23",
    "quantity": "85.00",
    "type": "Equity"
  },
  "USD": {
    "allocation": "9.9267%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "2049.78USD",
    "name": "USD",
    "order": {
      "marketValue": "-7950.22USD",
      "quantity": "-7950.22"
    },
    "price": "1.00",
    "quantity": "2049.78",
    "type": "Currency"
  }
}
```

Please note that a free API key is limited to 5 API calls per minute.
