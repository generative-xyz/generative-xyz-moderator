package adapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"rederinghub.io/pkg/config"
)

type MoralisNFTResponse struct {
	Result []*MoralisNFT `json:"result"`
	Total  int           `json:"total"`
}

type MoralisNFT struct {
	TokenID  string `json:"token_id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Amount   string `json:"amount"`
	Metadata string `json:"metadata"`
}

type MoralisAdapter interface {
	ListNFTs(contract string, chainID string) (*MoralisNFTResponse, error)
	ResyncNFTMetadata(contract string, chainID string, tokenID string) (error)
}

type moralisAdapter struct {
	url    string
	apiKey string
}

func (m *moralisAdapter) ListNFTs(contract string, chainID string) (*MoralisNFTResponse, error) {
	url := m.url + contract + "?chain=" + chainID + "&format=decimal"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-Key", m.apiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp MoralisNFTResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (m *moralisAdapter) ResyncNFTMetadata(contract string, chainID string, tokenID string) (error) {
	url := fmt.Sprintf("%v%v/%v/metadata/resync?chain=%v&mode=sync", m.url, contract, tokenID, chainID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-Key", m.apiKey)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func NewMoralisAdapter() MoralisAdapter {
	appConfig := config.AppConfig()
	if appConfig.MoralisURL == "" {
		panic(any("missing Moralis URL"))
	}
	return &moralisAdapter{
		url:    appConfig.MoralisURL,
		apiKey: appConfig.MoralisAPIKey,
	}
}
