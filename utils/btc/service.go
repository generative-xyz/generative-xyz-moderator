package btc

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type VIn struct {
	TxID string `json:"txid"`
	Vout int    `json:"vout"`
}

type Status struct {
	Confirmed bool `json:"confirmed"`
}

type GetUTXOResp struct {
	TxID   string `json:"txid"`
	Vin    []VIn  `json:"vin"`
	Status Status `json:"status"`
}

func FilterPendingUTXOs(utxos []UTXOType, address string) (pendingUTXOs []UTXOType, spendableUTXOs []UTXOType, err error) {
	pendingUTXOs = []UTXOType{}
	spendableUTXOs = []UTXOType{}

	client := resty.New()
	url := "https://blockstream.info/api/address/" + address + "/txs"

	response, err := client.R().Get(url)
	if err != nil {
		return
	}
	if response.StatusCode() != 200 {
		err = fmt.Errorf("Get pending txs status code %v", response.StatusCode())
		return
	}

	res := []GetUTXOResp{}
	err = json.Unmarshal(response.Body(), &res)
	if err != nil {
		err = fmt.Errorf("Unmarshal response error: %v", err)
		return
	}

	pendingVins := []VIn{}
	for _, tx := range res {
		if !tx.Status.Confirmed {
			pendingVins = append(pendingVins, tx.Vin...)
		}
	}

	for _, in := range pendingVins {
		for i, utxo := range utxos {
			if utxo.TxHash == in.TxID && utxo.TxOutIndex == in.Vout {
				pendingUTXOs = append(pendingUTXOs, utxo)
				utxos = append(utxos[:i], utxos[i+1:]...)
				break
			}
		}
	}
	spendableUTXOs = utxos
	return
}
