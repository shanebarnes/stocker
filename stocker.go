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
	// TODO: Add a flag for the number of API calls permitted per minute
	if len(apiKey) == 0 {
		apiKey = *(flag.String("apiKey", "", "Alpha Vantage API key"))
	}
	currency := flag.String("currency", "USD", "Currency")
	debug := flag.Bool("debug", true, "Debug mode")
	help := flag.Bool("help", false, "Display help information")
	portfolio := flag.String("portfolio", "", "Portfolio file containing source and target assets")
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
		log.Warn("Only ", av.ApiCallsPerMinLimit, " API calls to Alpha Vantage will be performed each minute")
		p, _ := NewPortfolio(*portfolio, apiKey, *currency)
		p.Rebalance()
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
