package stock

import (
	fp "github.com/robaho/fixed"
	"sync"
	"syscall"
)

type Cache struct {
	mpCcy  map[string]Currency
	mpQte  map[string]Quote
	mpSym  map[string]Symbol
	mtxCcy sync.RWMutex
	mtxQte sync.RWMutex
	mtxSym sync.RWMutex
}

func (c *Cache) AddCurrency(currency Currency) error {
	c.mtxCcy.Lock()
	defer c.mtxCcy.Unlock()
	if ccy, exists := c.mpCcy[currency.Currency]; exists {
		if currency.Rates == nil {
			currency.Rates = make(map[string]fp.Fixed)
		} else {
			for key, val := range currency.Rates {
				ccy.Rates[key] = val
			}
		}
		c.mpCcy[currency.Currency] = ccy
	} else {
		c.mpCcy[currency.Currency] = currency
	}
	return nil
}

func (c *Cache) AddQuote(quote Quote) error {
	c.mtxQte.Lock()
	defer c.mtxQte.Unlock()
	c.mpQte[quote.Symbol] = quote
	return nil
}

func (c *Cache) AddSymbol(symbol Symbol) error {
	c.mtxSym.Lock()
	defer c.mtxSym.Unlock()
	c.mpSym[symbol.Symbol] = symbol
	return nil
}

func (c *Cache) GetCurrency(currency, currencyTo string) (Currency, error) {
	var err error = syscall.ENOENT
	c.mtxCcy.RLock()
	defer c.mtxCcy.RUnlock()
	ccy, exists := c.mpCcy[currency]
	if exists {
		if _, exists = ccy.Rates[currencyTo]; exists {
			err = nil
		}
	}
	return ccy, err
}

func (c *Cache) GetQuote(quote string) (Quote, error) {
	var err error
	c.mtxQte.RLock()
	defer c.mtxQte.RUnlock()
	qte, exists := c.mpQte[quote]
	if !exists {
		err = syscall.ENOENT
	}
	return qte, err
}

func (c *Cache) GetSymbol(symbol string) (Symbol, error) {
	var err error
	c.mtxSym.RLock()
	defer c.mtxSym.RUnlock()
	sym, exists := c.mpSym[symbol]
	if !exists {
		err = syscall.ENOENT
	}
	return sym, err
}

func NewCache() *Cache {
	return &Cache{
		mpCcy: make(map[string]Currency),
		mpQte: make(map[string]Quote),
		mpSym: make(map[string]Symbol),
	}
}
