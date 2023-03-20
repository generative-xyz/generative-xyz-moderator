package btc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SpentVoutBlockStream struct {
	Spent  bool   `json:"spent"`
	Txid   string `json:"txid"`
	Vin    int    `json:"vin"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int    `json:"block_time"`
	} `json:"status"`
}

func CheckOutcoinSpentBlockStream(txhash string, vout uint) (string, error) {
	var result SpentVoutBlockStream
	resp, err := http.Get(fmt.Sprintf("https://blockstream.info/api/tx/%v/outspend/%v", txhash, vout))
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if result.Spent {
		return result.Txid, nil
	}
	return "", nil
}
