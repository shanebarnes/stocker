package stock

import fp "github.com/robaho/fixed"

type Currency struct {
	Currency string
	Name     string
	Rates    map[string]fp.Fixed // map[currencyTo]ExchangeRate
}

type price struct {
	Ask         float64
	Bid         float64
	Close       float64
	High        float64
	Low         float64
	Open        float64
	Latest      float64
	LatestTrHrs float64
}

type Quote struct {
	Prices price
	Symbol string
	Volume string
}

type Symbol struct {
	Currency    string
	Description string
	Id          string
	Symbol      string
	Type        string
}
