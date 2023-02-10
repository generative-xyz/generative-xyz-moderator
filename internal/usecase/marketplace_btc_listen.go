package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type Txs struct {
	Tx    string `json:"tx_hash"`
	Value string `json:"value" binding:"required"`
}

type Txo struct {
	Address string `json:"address"`
	Txs     []Txs  `json:"txrefs"`
}

func (u Usecase) buildBTCClient() (*rpcclient.Client, error) {
	host := u.Config.BTC_FULLNODE
	user := u.Config.BTC_RPCUSER
	pass := u.Config.BTC_RPCPASSWORD

	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   false, //!(os.Getenv("BTC_NODE_HTTPS") == "true"), // Bitcoin core does not provide TLS by default
	}
	return rpcclient.New(connCfg, nil)
}

// check nft of the nft:
func (u Usecase) BtcChecktListNft() error {

	btcClient, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingListings()
	if len(listPending) == 0 {
		return nil
	}

	for _, item := range listPending {

		txs, _ := u.getLastTxs(item.HoldOrdAddress)

		if len(txs) == 0 {
			continue
		}
		found := false
		for _, tx := range txs {
			detail, err := chainhash.NewHashFromStr(tx.Tx)
			if err != nil {
				fmt.Println("can not NewHashFromStr with err:", err)
				continue
			}
			result, _ := btcClient.GetRawTransactionVerboseAsync(detail).Receive()

			for _, vin := range result.Vin {
				if strings.Contains(vin.Txid, item.InscriptionID) {
					found = true
					item.IsConfirm = true
					_, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
					if err != nil {
						fmt.Println("UpdateBTCNFTConfirmListings", err.Error(), err)
					}
					break
				}
			}
			if found {
				break
			}
		}
	}

	return nil
}

func (u *Usecase) getLastTxs(address string) ([]Txs, error) {

	var txs []Txs

	// get last tx:
	url := u.Config.BlockcypherAPI + "addrs/" + address
	fmt.Println("url", url)

	req, err := http.NewRequest("GET", url, nil)
	var (
		result *Txo
	)

	if err != nil {
		fmt.Println("eth get list txs failed", address, err.Error())
		return txs, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Call eth get list txs failed", err.Error())
		return txs, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return txs, errors.New("get eth list txs Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Read body failed", err.Error())
		return txs, errors.New("Read body failed")
	}

	return result.Txs, nil

}
