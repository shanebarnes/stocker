package main

import (
	"flag"
	"fmt"
	"os"

	port "github.com/shanebarnes/stocker/internal/portfolio"
	"github.com/shanebarnes/stocker/internal/stock/api"
	ver "github.com/shanebarnes/stocker/internal/version"
	log "github.com/sirupsen/logrus"
)

var apiKey string
var apiServer string

func init() {
	format := new(log.TextFormatter)
	format.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	format.FullTimestamp = true

	log.SetFormatter(format)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	initEnvVars()
}

func initEnvVars() {
	apiKey = api.GetApiKeyFromEnv()
	apiServer = api.GetApiServerFromEnv()
}

func main() {
	key := flag.String("apiKey", "", "Stock API key")
	server := flag.String("apiServer", "", "Stock API server")
	currency := flag.String("currency", "USD", "Currency")
	debug := flag.Bool("debug", true, "Debug mode")
	//requests := flag.Int("requests", 5, "Maximum API requests per minute. The free API key only allows for 5 API requests per minute")
	help := flag.Bool("help", false, "Display help information")
	version := flag.Bool("version", false, "Display version information")
	portfolio := flag.String("rebalance", "", "Portfolio file containing source assets to rebalance against target assets")
	flag.Parse()

	if len(apiKey) == 0 {
		apiKey = *key
	}

	if len(apiServer) == 0 {
		apiServer = *server
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	//av.ApiRequestsPerMinLimit = *requests

	if *help {
		flag.PrintDefaults()
	} else if *version {
		fmt.Println("stocker version", ver.String())
	} else if len(apiKey) == 0 || len(apiServer) == 0 {
		flag.PrintDefaults()
	} else if len(*portfolio) > 0 {
		log.Warn("Rebalancing requires making stock API calls")
		//log.Warn("Only ", av.ApiRequestsPerMinLimit, " API calls to Alpha Vantage will be performed each minute")
		if p, err := port.NewPortfolio(*portfolio, apiKey, apiServer, *currency); err == nil {
			p.Rebalance()
		}
	} else {
		flag.PrintDefaults()
	}
}
