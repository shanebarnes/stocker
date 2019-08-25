# stocker
Rebalance your finanical assets with realtime stock market data acquired on demand from [Alpha Vantage](https://www.alphavantage.co/) APIs.

## Build Instructions

```shell
$ git clone https://github.com/shanebarnes/stocker.git
$ cd stocker
$ go build -v
$ go test -v ./...
```

## Examples

Try rebalancing a sample portfolio! An [Alpha Vantage](https://www.alphavantage.co/) API key is required. A free API key can be claimed [here](https://www.alphavantage.co/support/#api-key).

```shell
$ ./stocker -apiKey <your_api_key> -portfolio examples/portfolio.json
```

The [portfolio.json](https://github.com/shanebarnes/stocker/blob/master/examples/portfolio.json) contains source and target assets.

```json
{
    "assets": {
        "source": {
            "AAPL": {
                "quantity": 10
            },
            "CAD": {
                "type": "Currency",
                "quantity": 10000.00
            },
            "USD": {
                "type": "Currency",
                "quantity": 10000.00
            }
        },
        "target": {
            "CAD": {
                "type": "Currency",
                "allocation": 20.00
            },
            "MSFT": {
                "allocation": 70.00
            },
            "USD": {
                "type": "Currency",
                "allocation": 10.00
            }
        }
    }
}
```

Here is what the asset rebalancing output will look like for the sample portfolio.

```shell
WARN[2019-08-24T22:00:58.577-04:00] Rebalancing requires making Alpha Vantage API calls
WARN[2019-08-24T22:00:58.577-04:00] Only 5 API calls to Alpha Vantage will be performed each minute
INFO[2019-08-24T22:00:58.578-04:00] Liquidating source assets into USD funds
DEBU[2019-08-24T22:00:58.578-04:00] AAPL: searching for symbol information
DEBU[2019-08-24T22:00:59.086-04:00] AAPL: searching for symbol quote information
DEBU[2019-08-24T22:01:13.646-04:00] CAD: searching for symbol information
DEBU[2019-08-24T22:01:13.646-04:00] CAD: searching for exchange rate from CAD to USD
DEBU[2019-08-24T22:01:28.720-04:00] USD: searching for symbol information
INFO[2019-08-24T22:01:28.720-04:00] source portfolio:{
  "AAPL": {
    "avgPrice": 202.64,
    "exchangeRate": 1,
    "marketValue": 2026,
    "allocation": 10.35925020708274,
    "quantity": 10,
    "type": "Equity"
  },
  "CAD": {
    "avgPrice": 1,
    "exchangeRate": 0.7531,
    "marketValue": 7531,
    "allocation": 38.50716352889443,
    "quantity": 10000,
    "type": "currency"
  },
  "USD": {
    "avgPrice": 1,
    "exchangeRate": 1,
    "marketValue": 10000,
    "allocation": 51.131541002382725,
    "quantity": 10000,
    "type": "currency"
  }
}
INFO[2019-08-24T22:01:28.720-04:00] Re-allocating source assets to match target allocations
DEBU[2019-08-24T22:01:28.720-04:00] MSFT: searching for symbol information
DEBU[2019-08-24T22:01:43.784-04:00] MSFT: searching for symbol quote information
INFO[2019-08-24T22:01:58.650-04:00] target portfolio:{
  "CAD": {
    "avgPrice": 1,
    "exchangeRate": 0.7531,
    "marketValue": 3911,
    "allocation": 19.997545686031884,
    "quantity": 3911,
    "type": "Currency"
  },
  "MSFT": {
    "avgPrice": 133.39,
    "exchangeRate": 1,
    "marketValue": 13605.779999999999,
    "allocation": 69.56844979393989,
    "quantity": 102,
    "type": "Equity"
  },
  "USD": {
    "avgPrice": 1,
    "exchangeRate": 1,
    "marketValue": 2040.6200000000026,
    "allocation": 10.434004520028237,
    "quantity": 2040.6200000000026,
    "type": "Currency"
  }
}
```

Please note that a free API key is limited to 5 API calls per minute.
