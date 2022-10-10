package exchange

import (
	coingecko "github.com/superoo7/go-gecko/v3"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"

	"net/http"
	"strings"
	"time"
)

type coingeckoImpl struct {
	client *coingecko.Client
}

func (c coingeckoImpl) GetRate(currency string) (float64, error) {
	currency = strings.TrimSpace(currency)
	coidID, hasID := cryptocurrency.CoingeckoCoinIdByCurrency[currency]
	if !hasID {
		return 0, ErrCurrencyNotSupported
	}

	results, err := c.client.SimplePrice([]string{coidID}, []string{cryptocurrency.USD})
	if err != nil {
		return 0, err
	}

	resultsMap := *results
	return float64(resultsMap[coidID][cryptocurrency.USD]), nil
}

func NewCoingeckoClient() Client {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)
	return &coingeckoImpl{
		client: cg,
	}
}
