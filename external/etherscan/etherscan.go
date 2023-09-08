package etherscan

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

type EtherscanService struct {
	conf      *config.Config
	serverURL string
	apiKey    string
	cache     redis.IRedisCache
}

func NewEtherscanService(conf *config.Config, cache redis.IRedisCache) *EtherscanService {
	return &EtherscanService{
		conf:      conf,
		serverURL: os.Getenv("ETHERSCAN_API_URL"),
		apiKey:    os.Getenv("ETHERSCAN_API_KEY"),
		cache:     cache,
	}
}

func (m EtherscanService) generateUrl(params url.Values) string {

	fullUrl := fmt.Sprintf("%s?%s", m.serverURL, params.Encode())
	return fullUrl
}

func (m EtherscanService) AddressBalance(address string) (*WalletAddressResponse, error) {
	queries := url.Values{}
	queries.Set("module", ModuleAccount)
	queries.Set("action", ActionBalance)
	queries.Set("address", address)
	queries.Set("tag", "latest")
	queries.Set("apikey", m.apiKey)

	fullUrl := m.generateUrl(queries)
	data, _, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &WalletAddressResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m EtherscanService) AddressTransactions(address string) (*WalletAddressTxResponse, error) {
	queries := url.Values{}
	queries.Set("module", ModuleAccount)
	queries.Set("action", ActionTxList)
	queries.Set("address", address)
	queries.Set("tag", "latest")
	queries.Set("apikey", m.apiKey)

	fullUrl := m.generateUrl(queries)
	data, _, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &WalletAddressTxResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m EtherscanService) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, int, error) {

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
