package explorer_trustless

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"rederinghub.io/utils/redis"
)

type ExplorerTrustless struct {
	serverURL string
	cache     redis.IRedisCache
}

func NewExplorerTrustless(cache redis.IRedisCache) *ExplorerTrustless {

	serverURL := os.Getenv("EXPLORER_TRUSTLESS_API")
	if serverURL == "" {
		serverURL = "https://explorer.trustless.computer"
	}
	return &ExplorerTrustless{
		serverURL: serverURL,
		cache:     cache,
	}
}

type metadataChan struct {
	Key int
	Err error
}

func (m ExplorerTrustless) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}

func (m ExplorerTrustless) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Response with status %d", res.StatusCode))
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (m ExplorerTrustless) TokenHolders(tokenAddress string) (*InscriptioResponse, error) {
	url := fmt.Sprintf(TOKEN_HOLDERS, tokenAddress)
	fullUrl := m.generateUrl(url)
	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &InscriptioResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
