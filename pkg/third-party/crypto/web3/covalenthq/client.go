package covalenthq

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl string
	apiKey  string
}

func NewClient(baseUrl string, apiKey string) *Client {
	return &Client{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}

func (c Client) GetBalance(address string, chainID int, noFetchNFTMetadata bool) (*Balance, error) {
	// set chainID to solona in case it equal 0
	if chainID == 0 {
		chainID = chainIDSolana
	}

	covalenthqURL := fmt.Sprintf("%s/%d/address/%s/balances_v2/?nft=true&key=%s", c.baseUrl, chainID, address, c.apiKey)
	if noFetchNFTMetadata {
		covalenthqURL += "&no-nft-fetch=true"
	}
	req, _ := http.NewRequest(http.MethodGet, covalenthqURL, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	result := &Balance{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c Client) GetNFTFeeds(address string, chainID int) ([]interface{}, error) {
	b, err := c.GetBalance(address, chainID, false)
	if err != nil {
		return nil, err
	}

	return b.GetNFTItems(), nil
}
