package mempool_space

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
)

type MempoolService struct {
	conf      *config.Config
	serverURL string
	//apiKey    string
	cache redis.IRedisCache
}

func NewMempoolService(conf *config.Config, cache redis.IRedisCache) *MempoolService {
	return &MempoolService{
		conf:      conf,
		serverURL: os.Getenv("MEMPOOL_API_URL"),
		//apiKey:    os.Getenv("ETHERSCAN_API_KEY"),
		cache: cache,
	}
}

func (m MempoolService) generateUrl(path string) string {

	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}

func (m MempoolService) AddressTransactions(address string) ([]AddressTxItemResponse, error) {
	path := fmt.Sprintf(TxAddressURL, address)
	fullUrl := m.generateUrl(path)
	data, _, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := []AddressTxItemResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m MempoolService) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, int, error) {

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
