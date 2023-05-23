package token_explorer

import (
	"fmt"
	"net/url"
	"os"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

type TokenExplorer struct {
	serverURL string
	cache     redis.IRedisCache
}

func NewTokenExplorer(cache redis.IRedisCache) *TokenExplorer {
	sUrl := os.Getenv("TOKEN_EXPLORER_URL")
	if sUrl == "" {
		sUrl = "https://api-token.trustless.computer/api/v1"
	}
	return &TokenExplorer{
		serverURL: sUrl,
		cache:     cache,
	}
}

func (q *TokenExplorer) Tokens(params url.Values) ([]Token, error) {
	headers := make(map[string]string)
	fUrl := fmt.Sprintf("%s/%s", q.serverURL, "tokens")
	data, _, _, err := helpers.JsonRequest(fmt.Sprintf("%s?%s", fUrl, params.Encode()), "GET", headers, nil)
	if err != nil {
		return nil, err
	}

	resp, err := q.ParseData(data)
	if err != nil {
		return nil, err
	}

	return resp.ToTokens()
}

func (q *TokenExplorer) Search(params url.Values) (*SearchToken, error) {
	headers := make(map[string]string)
	fUrl := fmt.Sprintf("%s/%s", q.serverURL, "search")
	data, _, _, err := helpers.JsonRequest(fmt.Sprintf("%s?%s", fUrl, params.Encode()), "GET", headers, nil)
	if err != nil {
		return nil, err
	}

	resp, err := q.ParseData(data)
	if err != nil {
		return nil, err
	}

	return resp.ToSearchTokens()
}

func (q *TokenExplorer) Token(address string) (*Token, error) {
	data, _, _, err := helpers.JsonRequest(fmt.Sprintf("%s/%s/%s", q.serverURL, "token", address), "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := q.ParseData(data)
	if err != nil {
		return nil, err
	}
	return resp.ToToken()
}

func (q *TokenExplorer) WalletAddressTokens(walletAddress string, params url.Values) ([]WalletAddressToken, error) {
	headers := make(map[string]string)
	url := fmt.Sprintf("%s/%s/%s", q.serverURL, walletAddress, "tokens")
	data, _, _, err := helpers.JsonRequest(fmt.Sprintf("%s?%s", url, params.Encode()), "GET", headers, nil)
	if err != nil {
		return nil, err
	}

	resp, err := q.ParseData(data)
	if err != nil {
		return nil, err
	}

	return resp.ToWalletAddressTokens()
}

func (q *TokenExplorer) ParseData(data []byte) (*Response, error) {
	resp := &Response{}
	err := helpers.ParseData(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != "OK" {
		return nil, resp.Error
	}

	return resp, nil
}
