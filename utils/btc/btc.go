package btc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/blockcypher/gobcy/v2"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

func NewBlockcypherService(chainEndpoint string, explorerEndPoint string, bcyToken string, env *chaincfg.Params) *BlockcypherService {

	return &BlockcypherService{
		chainEndpoint:    chainEndpoint,
		explorerEndPoint: explorerEndPoint,
		bcyToken:         bcyToken,
		network:          env,
		chain:            gobcy.API{bcyToken, "btc", "main"},
	}
}

// SendTX from Segwit address by lib gobcy, with preference, without manually setting fees :
//send from segwit to legacy address |
func (bs *BlockcypherService) SendTransactionWithPreferenceFromSegwitAddress(secret string, from string, destination string, amount int,
	preference string) (string, error) {
	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return "", err
	}

	pkHex := hex.EncodeToString(wif.PrivKey.Serialize())
	tx := gobcy.TempNewTX(from, destination, *big.NewInt(int64(amount)))

	//tx := gobcy.TX{//fields} // send multi: TODO support

	if len(preference) == 0 {
		tx.Preference = PreferenceMedium
	} else {
		tx.Preference = preference
	}
	skel, err := bs.chain.NewTX(tx, false) // gobcy.TX

	if err != nil {
		log.Println("bs.chain.NewTX err: ", err, tx)
		return "", err
	}
	log.Println("[SendTransactionWithPreference] fee", skel.Trans.Fees)
	prikHexs := []string{}
	for i := 0; i < len(skel.ToSign); i++ {
		prikHexs = append(prikHexs, pkHex)
	}

	err = skel.Sign(prikHexs)
	if err != nil {
		log.Println("skel.Sign error: ", err)
		return "", err
	}

	// add this one with segwit address:
	for i, _ := range skel.Signatures {
		skel.Signatures[i] = skel.Signatures[i] + "01"
	}

	skel, err = bs.chain.SendTX(skel)
	if err != nil {
		log.Println("bs.chain.SendTX err:", err)
		return "", err
	}
	return skel.Trans.Hash, nil
}

func (bs *BlockcypherService) GetBalance(address string) (*big.Int, int, error) {

	b := new(big.Int)

	confirm := 0

	btcAddrInfo, err := bs.BTCGetAddrInfo(address)
	if err != nil {
		return nil, confirm, errors.Wrap(err, "c.btc.BalanceAt")
	}
	// check confirmations number: 6
	if len(btcAddrInfo.TxRefs) > 0 {
		for _, tx := range btcAddrInfo.TxRefs {
			fmt.Println("btcAddrInfo.TxRefs[0].Confirmations ", tx.Confirmations)
			confirm = tx.Confirmations
			break
			// if tx.Confirmations < 6 {
			// 	return b, nil
			// } else {
			// 	return b, errors.Errorf("need 6 confirm, current confirm %d", tx.Confirmations)
			// }
		}

	}

	b.SetUint64(btcAddrInfo.Balance)
	return b, confirm, nil
}

// BTCGetAddrInfo :
func (bs *BlockcypherService) BTCGetAddrInfo(address string) (*AddrInfo, error) {
	url := fmt.Sprintf("%s/%s?limit=1&unspentOnly=true&includeScript=false&token=%s", bs.chainEndpoint, address, bs.bcyToken)
	fmt.Println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	var (
		result *AddrInfo
	)

	if err != nil {
		fmt.Println("BTC get UTXO failed", address, err.Error())
		return result, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Call BTC get UTXO failed", err.Error())
		return result, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return result, errors.New("GetUTXO Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Read body failed", err.Error())
		return result, errors.New("Read body failed")
	}

	return result, nil
}

func (bs *BlockcypherService) GetLastTxs(address string) ([]Txs, error) {

	var txs []Txs

	// get last tx:
	url := bs.chainEndpoint + "/" + address + "?token=" + bs.bcyToken
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

	if err != nil {
		fmt.Println("Read body failed", err.Error())
		return txs, errors.New("Read body failed")
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("Unmarshaly failed", err.Error())
		return txs, errors.New("Unmarshal failed")
	}

	return result.Txs, nil

}

// gen a segwit address:
func GenerateAddressSegwit() (privKey, pubKey, addressSegwit string, err error) {

	secret, err := btcec.NewPrivateKey()
	if err != nil {
		err = errors.Wrap(err, "c.GenerateAddressSegwit")
		return
	}

	wif, err := btcutil.NewWIF(secret, &chaincfg.MainNetParams, true)
	if err != nil {
		err = errors.Wrap(err, "c.GenerateAddressSegwit")
		return
	}

	privKey = wif.String()

	witnessProg := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)

	if err != nil {
		err = errors.Wrap(err, "btcutil.NewAddressWitnessPubKeyHash")
		return
	}

	addressSegwit = addressWitnessPubKeyHash.EncodeAddress()

	return
}
