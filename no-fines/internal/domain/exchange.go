package domain

type CurrencyPair struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

type ExchangeRate struct {
	CurrencyPair
	Rate float64 `json:"rate"`
}

type ExchangeProvider interface {
	GetRates(pair CurrencyPair) ([]ExchangeRate, error)
}

// mock
type MockProvider struct{}

func (p *MockProvider) GetRates(pair CurrencyPair) ([]ExchangeRate, error) {
	return []ExchangeRate{
		{CurrencyPair: pair, Rate: 0.0123},
	}, nil
}
