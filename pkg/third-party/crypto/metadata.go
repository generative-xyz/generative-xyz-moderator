package crypto

import (
	"strings"
	"time"

	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/types"
)

type Metadata struct {
	Rate           float64      `json:"rate"`
	Currency       string       `json:"currency"`
	TotalAmount    *types.Price `json:"totalAmount"`
	CurrencySymbol string       `json:"currencySymbol"`
	USDAmount      float64      `json:"usdAmount"`
	RefreshTime    time.Time    `json:"refreshTime"`
}

func NewMetadata() *Metadata {
	return &Metadata{
		RefreshTime: time.Now().Add(time.Minute * 5), // 5 minutes from now
	}
}

func (c *Metadata) WithTotalAmount(cartTotalAmount float64) *Metadata {
	cryptoAmount := cartTotalAmount / c.Rate
	coinRoundPlaces := cryptocurrency.RoundPlacesByCurrency[c.Currency]
	c.TotalAmount = types.NewPrice(
		cryptoAmount,
		types.WithPrecision(cryptocurrency.DefaultRoundPlaces),
		types.WithCryptoPrecision(coinRoundPlaces),
	)
	c.USDAmount = cartTotalAmount
	return c
}

func (c *Metadata) WithCurrency(currency string) *Metadata {
	c.Currency = strings.ToLower(currency)
	return c
}

func (c *Metadata) WithRate(rate float64) *Metadata {
	c.Rate = rate
	return c
}

func (c *Metadata) WithCurrencySymbol(currencySymbol string) *Metadata {
	c.CurrencySymbol = currencySymbol
	return c
}

func (c *Metadata) GetRate() float64 {
	return c.Rate
}

func (c *Metadata) GetCurrencySymbol() string {
	return c.CurrencySymbol
}

func (c *Metadata) GetCurrency() string {
	return c.Currency
}

func (c *Metadata) GetTotalAmount() *types.Price {
	return c.TotalAmount
}
