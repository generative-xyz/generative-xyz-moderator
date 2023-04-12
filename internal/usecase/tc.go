package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (u Usecase) GetNftsByAddress(address string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/nft-explorer/nfts?limit=1&page=2&owner=%s", address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result struct {
		Data []*struct {
			Collection        string `json:"collection"`
			CollectionAddress string `json:"collection_address"`
			TokenID           string `json:"token_id"`
			Name              string `json:"name"`
			ContentType       string `json:"content_type"`
			Image             string `json:"image"`
			Explorer          string `json:"explorer"`
			ArtistName        string `json:"artist_name"`
		} `json:"data"`
	}

	if true {
		return result.Data, nil
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for _, nft := range result.Data {
		if len(nft.Image) > 0 {
			nft.Image += "/content"
		}
		nft.Explorer = fmt.Sprintf("https://trustless.computer/inscription?contract=%s&id=%s", nft.CollectionAddress, nft.TokenID)
		nft.ArtistName = "" // todo update late
	}
	return result.Data, err
}
