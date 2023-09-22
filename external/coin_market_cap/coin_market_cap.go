package coin_market_cap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
)

type CoinMarketCap struct {
	conf      *config.Config
	serverURL string
	apiKey    string
	cache     redis.IRedisCache
}

func NewCoinMarketCap(conf *config.Config, cache redis.IRedisCache) *CoinMarketCap {
	apiURL := os.Getenv("COINMAKET_CAP_API_URL")
	if apiURL == "" {
		apiURL = "https://pro-api.coinmarketcap.com/v2"
	}
	return &CoinMarketCap{
		conf:      conf,
		serverURL: apiURL,
		apiKey:    os.Getenv("COINMAKET_CAP_API_KEY"),
		cache:     cache,
	}
}

func (m CoinMarketCap) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}

func (m CoinMarketCap) PriceConversion(coinID int) (*PriceConversionResponse, error) {
	urlQueries := url.Values{}
	urlQueries.Set("amount", "1")
	urlQueries.Set("id", fmt.Sprintf("%d", coinID))

	path := fmt.Sprintf("%s?%s", PriceConversion, urlQueries.Encode())
	fullUrl := m.generateUrl(path)
	data, _, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &PriceConversionResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m CoinMarketCap) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, int, error) {

	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, 0, err
	}

	if len(headers) > 0 {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	// remove this for error response:
	// if res.StatusCode != http.StatusOK {
	// 	err = errors.New(fmt.Sprintf("Response with status %d", res.StatusCode))
	// 	return nil, statusCode, err
	// }

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return body, res.StatusCode, nil
}
