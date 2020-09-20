# stocker
Rebalance your financial assets with realtime stock market data acquired on demand from [Alpha Vantage](https://www.alphavantage.co/) or [Questrade](https://www.questrade.com) APIs. For example, get the latest stock prices or foreign exchange rates.

## Build Instructions

```shell
$ git clone https://github.com/shanebarnes/stocker.git
$ cd stocker
$ ./build/build.sh
```

## Examples

Try rebalancing a sample portfolio! An [Alpha Vantage](https://www.alphavantage.co/) API key is required. A free API key can be claimed [here](https://www.alphavantage.co/support/#api-key).

```shell
$ # Pass API key on command line
$ ./stocker -apiKey <your_api_key> -apiServer www.alphavantage.co -rebalance ./examples/portfolio.json
$
$ # Alternatively, load API key from environment
$ STOCKER_API_KEY=<your_api_key>; STOCKER_API_SERVER=www.alphavantage.co; ./stocker -rebalance ./examples/portfolio.json
```

Here is an example of currency conversion.

```shell
$ ./stocker -rebalance ../../examples/currency_conversion.json
WARN[2020-02-01T21:34:22.148-05:00] Rebalancing requires making Alpha Vantage API calls
WARN[2020-02-01T21:34:22.148-05:00] Only 5 API calls to Alpha Vantage will be performed each minute
INFO[2020-02-01T21:34:22.149-05:00] Validating source assets
DEBU[2020-02-01T21:34:22.149-05:00] USD: searching for symbol information
INFO[2020-02-01T21:34:22.149-05:00] Validating target assets
DEBU[2020-02-01T21:34:22.149-05:00] EUR: searching for symbol information
DEBU[2020-02-01T21:34:22.149-05:00] EUR: searching for exchange rate from EUR to USD
INFO[2020-02-01T21:34:22.574-05:00] Liquidating source assets into USD funds
DEBU[2020-02-01T21:34:22.574-05:00] USD: searching for symbol information
INFO[2020-02-01T21:34:22.574-05:00] source portfolio:{
  "USD": {
    "allocation": "100.0000%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "10000.00USD",
    "name": "USD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  }
}
INFO[2020-02-01T21:34:22.574-05:00] Re-allocating source assets to match target allocations
DEBU[2020-02-01T21:34:22.574-05:00] EUR: searching for symbol information
DEBU[2020-02-01T21:34:22.574-05:00] EUR: searching for exchange rate from EUR to USD
DEBU[2020-02-01T21:34:22.574-05:00] EUR: found cached exchange rate to USD: 1.1094
INFO[2020-02-01T21:34:22.574-05:00] target portfolio:{
  "EUR": {
    "allocation": "100.0000%",
    "currency": "EUR",
    "exchangeRate": "1.1094",
    "marketValue": "10000.00USD",
    "name": "EUR",
    "order": {
      "marketValue": "+10000.00USD",
      "quantity": "+9013.88"
    },
    "price": "1.00",
    "quantity": "9013.88",
    "type": "currency"
  },
  "USD": {
    "allocation": "0.0000%",
    "currency": "USD",
    "exchangeRate": "1.0000",
    "marketValue": "0.00USD",
    "name": "USD",
    "order": {
      "marketValue": "-10000.00USD",
      "quantity": "-10000.00"
    },
    "price": "1.00",
    "quantity": "0.00",
    "type": "currency"
  }
}
```

Here is an example of stock portfolio rebalancing in Canadian dollars.

```shell
$ ./stocker -rebalance ../../examples/portfolio.json -currency CAD
WARN[2020-02-01T21:42:01.759-05:00] Rebalancing requires making Alpha Vantage API calls
WARN[2020-02-01T21:42:01.760-05:00] Only 5 API calls to Alpha Vantage will be performed each minute
INFO[2020-02-01T21:42:01.760-05:00] Validating source assets
DEBU[2020-02-01T21:42:01.760-05:00] CAD: searching for symbol information
DEBU[2020-02-01T21:42:01.760-05:00] MSFT: searching for symbol information
DEBU[2020-02-01T21:42:02.219-05:00] MSFT: searching for symbol quote information
DEBU[2020-02-01T21:42:16.828-05:00] MSFT: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:42:31.807-05:00] AAPL: searching for symbol information
DEBU[2020-02-01T21:42:46.864-05:00] AAPL: searching for symbol quote information
DEBU[2020-02-01T21:43:02.154-05:00] AAPL: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:43:02.154-05:00] USD: found cached exchange rate to CAD: 1.3229
INFO[2020-02-01T21:43:02.154-05:00] Validating target assets
DEBU[2020-02-01T21:43:02.154-05:00] TSLA: searching for symbol information
DEBU[2020-02-01T21:43:16.899-05:00] TSLA: searching for symbol quote information
DEBU[2020-02-01T21:43:31.839-05:00] TSLA: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:43:31.839-05:00] USD: found cached exchange rate to CAD: 1.3229
DEBU[2020-02-01T21:43:31.839-05:00] AMZN: searching for symbol information
DEBU[2020-02-01T21:43:46.936-05:00] AMZN: searching for symbol quote information
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:44:01.850-05:00] USD: found cached exchange rate to CAD: 1.3229
DEBU[2020-02-01T21:44:01.850-05:00] CAD: searching for symbol information
INFO[2020-02-01T21:44:01.850-05:00] Liquidating source assets into CAD funds
DEBU[2020-02-01T21:44:01.850-05:00] CAD: searching for symbol information
DEBU[2020-02-01T21:44:01.850-05:00] MSFT: searching for symbol information
DEBU[2020-02-01T21:44:01.850-05:00] MSFT: found cached symbol information
DEBU[2020-02-01T21:44:01.850-05:00] MSFT: searching for symbol quote information
DEBU[2020-02-01T21:44:01.850-05:00] MSFT: found cached symbol quote
DEBU[2020-02-01T21:44:01.850-05:00] MSFT: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:44:01.850-05:00] USD: found cached exchange rate to CAD: 1.3229
DEBU[2020-02-01T21:44:01.850-05:00] AAPL: searching for symbol information
DEBU[2020-02-01T21:44:01.850-05:00] AAPL: found cached symbol information
DEBU[2020-02-01T21:44:01.850-05:00] AAPL: searching for symbol quote information
DEBU[2020-02-01T21:44:01.850-05:00] AAPL: found cached symbol quote
DEBU[2020-02-01T21:44:01.850-05:00] AAPL: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:44:01.850-05:00] USD: found cached exchange rate to CAD: 1.3229
INFO[2020-02-01T21:44:01.850-05:00] source portfolio:{
  "AAPL": {
    "allocation": "32.5001%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "10236.27CAD",
    "name": "Apple Inc.",
    "price": "309.51",
    "quantity": "25.00",
    "type": "Equity"
  },
  "CAD": {
    "allocation": "31.7499%",
    "currency": "CAD",
    "exchangeRate": "1.0000",
    "marketValue": "10000.00CAD",
    "name": "CAD",
    "price": "1.00",
    "quantity": "10000.00",
    "type": "currency"
  },
  "MSFT": {
    "allocation": "35.7500%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "11259.86CAD",
    "name": "Microsoft Corporation",
    "price": "170.23",
    "quantity": "50.00",
    "type": "Equity"
  }
}
INFO[2020-02-01T21:44:01.850-05:00] Re-allocating source assets to match target allocations
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: searching for symbol information
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: found cached symbol information
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: searching for symbol quote information
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: found cached symbol quote
DEBU[2020-02-01T21:44:01.850-05:00] AMZN: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:44:01.850-05:00] USD: found cached exchange rate to CAD: 1.3229
DEBU[2020-02-01T21:44:01.850-05:00] TSLA: searching for symbol information
DEBU[2020-02-01T21:44:01.851-05:00] TSLA: found cached symbol information
DEBU[2020-02-01T21:44:01.851-05:00] TSLA: searching for symbol quote information
DEBU[2020-02-01T21:44:01.851-05:00] TSLA: found cached symbol quote
DEBU[2020-02-01T21:44:01.851-05:00] TSLA: searching for exchange rate from USD to CAD
DEBU[2020-02-01T21:44:01.851-05:00] USD: found cached exchange rate to CAD: 1.3229
INFO[2020-02-01T21:44:01.852-05:00] target portfolio:{
  "AAPL": {
    "allocation": "0.0000%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "0.00CAD",
    "name": "Apple Inc.",
    "order": {
      "marketValue": "-10236.27CAD",
      "quantity": "-25.00"
    },
    "price": "309.51",
    "quantity": "0.00",
    "type": "Equity"
  },
  "AMZN": {
    "allocation": "59.0592%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "18601.35CAD",
    "name": "Amazon.com Inc.",
    "order": {
      "marketValue": "+18601.35CAD",
      "quantity": "+7.00"
    },
    "price": "2008.72",
    "quantity": "7.00",
    "type": "Equity"
  },
  "CAD": {
    "allocation": "8.1506%",
    "currency": "CAD",
    "exchangeRate": "1.0000",
    "marketValue": "2567.11CAD",
    "name": "CAD",
    "order": {
      "marketValue": "-7432.89CAD",
      "quantity": "-7432.89"
    },
    "price": "1.00",
    "quantity": "2567.11",
    "type": "Currency"
  },
  "MSFT": {
    "allocation": "0.0000%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "0.00CAD",
    "name": "Microsoft Corporation",
    "order": {
      "marketValue": "-11259.86CAD",
      "quantity": "-50.00"
    },
    "price": "170.23",
    "quantity": "0.00",
    "type": "Equity"
  },
  "TSLA": {
    "allocation": "32.7903%",
    "currency": "USD",
    "exchangeRate": "1.3229",
    "marketValue": "10327.67CAD",
    "name": "Tesla Inc.",
    "order": {
      "marketValue": "+10327.67CAD",
      "quantity": "+12.00"
    },
    "price": "650.57",
    "quantity": "12.00",
    "type": "Equity"
  }
}
```
