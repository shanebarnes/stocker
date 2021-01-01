# stocker
Practice rebalancing your financial assets with realtime stock market data acquired on demand from [Alpha Vantage](https://www.alphavantage.co/) or [Questrade](https://www.questrade.com) APIs. For example, get the latest stock prices or foreign exchange rates.

Please note that this project is not production ready and should only be used with practice portfolios where no real money is at risk.

![build workflow](https://github.com/shanebarnes/stocker/workflows/stocker/badge.svg)

## Build Instructions

```shell
$ git clone https://github.com/shanebarnes/stocker.git
$ cd stocker
$ ./build/build.sh
```

## Examples

Try rebalancing a sample portfolio!

### Alpha Vantage

An [Alpha Vantage](https://www.alphavantage.co/) API key is required. A free API key can be claimed [here](https://www.alphavantage.co/support/#api-key).

```shell
$ # Pass API key on command line
$ ./bin/stocker-darwin -apiKey <your_api_key> -apiServer alphavantage.co -rebalance ./examples/portfolio.json
$
$ # Alternatively, load API key from environment
$ STOCKER_API_KEY=<your_api_key> STOCKER_API_SERVER=alphavantage.co ./stocker -debug -rebalance ./examples/portfolio.json
```

### Questrade

The stocker app must be [registered](https://www.questrade.com/api/documentation/getting-started) with Questrade. Here is an example of using OAuth credentials with refresh for use with Questrade APIs.

```shell
$ ./bin/stocker-darwin -apiServer questrade.com -credentials ./examples/credentials.json -rebalance ./examples/portfolio.json -refresh
```

### Miscellaneous

Here is an example of currency conversion.

```shell
$ STOCKER_API_KEY=<your_api_key> STOCKER_API_SERVER=alphavantage.co ./bin/stocker-darwin -rebalance ./examples/currency_conversion.json -currency EUR
```

Here is an example of stock portfolio rebalancing in Canadian dollars.

```shell
$ STOCKER_API_KEY=<your_api_key> STOCKER_API_SERVER=alphavantage.co ./bin/stocker-darwin -rebalance ./examples/portfolio.json -currency CAD
```