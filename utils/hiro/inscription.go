package hiro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetInscriptionByID(id string) (*InscriptionByIDRespone, error) {
	url := "https://api.hiro.so/ordinals/v1/inscriptions/" + id
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

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
	var result InscriptionByIDRespone
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &result, nil
}

type InscriptionByIDRespone struct {
	ID                 string `json:"id"`
	Number             int    `json:"number"`
	Address            string `json:"address"`
	GenesisAddress     string `json:"genesis_address"`
	GenesisBlockHeight int    `json:"genesis_block_height"`
	GenesisBlockHash   string `json:"genesis_block_hash"`
	GenesisTxID        string `json:"genesis_tx_id"`
	GenesisFee         string `json:"genesis_fee"`
	GenesisTimestamp   int64  `json:"genesis_timestamp"`
	TxID               string `json:"tx_id"`
	Location           string `json:"location"`
	Output             string `json:"output"`
	Value              string `json:"value"`
	Offset             string `json:"offset"`
	SatOrdinal         string `json:"sat_ordinal"`
	SatRarity          string `json:"sat_rarity"`
	SatCoinbaseHeight  int    `json:"sat_coinbase_height"`
	MimeType           string `json:"mime_type"`
	ContentType        string `json:"content_type"`
	ContentLength      int    `json:"content_length"`
	Timestamp          int64  `json:"timestamp"`
	CurseType          any    `json:"curse_type"`
	Recursive          bool   `json:"recursive"`
	RecursionRefs      any    `json:"recursion_refs"`
}
