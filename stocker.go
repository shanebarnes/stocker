package main

import (
	"encoding/json"
	"flag"
	"os"

	av "github.com/shanebarnes/stocker/alphavantage"
	log "github.com/sirupsen/logrus"
)

var apiKey = av.ApiGetKeyFromEnv()

func init() {
	format := new(log.TextFormatter)
	format.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	format.FullTimestamp = true

	log.SetFormatter(format)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	key := flag.String("apiKey", "", "Alpha Vantage API key")
	if len(apiKey) == 0 {
		apiKey = *key
	}
	currency := flag.String("currency", "USD", "Currency")
	debug := flag.Bool("debug", true, "Debug mode")
	av.ApiRequestsPerMinLimit = *(flag.Int("requests", 5, "Maximum API requests per minute. The free API key only allows for 5 API requests per minute"))
	help := flag.Bool("help", false, "Display help information")
	portfolio := flag.String("rebalance", "", "Portfolio file containing source assets to rebalance against target assets")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *help {
		flag.PrintDefaults()
	} else if len(apiKey) == 0 {
		flag.PrintDefaults()
	} else if len(*portfolio) > 0 {
		log.Warn("Rebalancing requires making Alpha Vantage API calls")
		log.Warn("Only ", av.ApiRequestsPerMinLimit, " API calls to Alpha Vantage will be performed each minute")
		if p, err := NewPortfolio(*portfolio, apiKey, *currency); err == nil {
			p.Rebalance()
		}
	} else {
		flag.PrintDefaults()
	}
}

func getPrettyString(v interface{}) string {
	str := ""
	buf, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		str = string(buf)
	}
	return str
}
