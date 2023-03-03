package btc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

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

func TempNewTXMultiOut(inAddr string, outAddrs map[string]big.Int) (trans gobcy.TX) {
	trans.Inputs = make([]gobcy.TXInput, 1)
	trans.Inputs[0].Addresses = make([]string, 1)
	trans.Inputs[0].Addresses[0] = inAddr
	for addr, amount := range outAddrs {
		newOutput := gobcy.TXOutput{
			Value:     amount,
			Addresses: []string{addr},
		}
		trans.Outputs = append(trans.Outputs, newOutput)
	}
	return
}

func (bs *BlockcypherService) EstimateFeeTransactionWithPreferenceFromSegwitAddressMultiAddress(secret string, from string, destinations map[string]int,
	preference string) (*big.Int, error) {

	outAddrs := make(map[string]big.Int)

	for addr, amount := range destinations {
		outAddrs[addr] = *big.NewInt(int64(amount))
	}

	tx := TempNewTXMultiOut(from, outAddrs)

	if len(preference) == 0 {
		tx.Preference = PreferenceMedium
	} else {
		tx.Preference = preference
	}
	skel, err := bs.chain.NewTX(tx, false) // gobcy.TX

	if err != nil {
		log.Println("bs.chain.NewTX err: ", err, tx)
		return nil, err
	}
	log.Println("[SendTransactionWithPreference] fee", skel.Trans.Fees)
	return &skel.Trans.Fees, nil
}

func (bs *BlockcypherService) SendTransactionWithPreferenceFromSegwitAddressMultiAddress(secret string, from string, destinations map[string]int,
	preference string) (string, error) {
	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return "", err
	}

	pkHex := hex.EncodeToString(wif.PrivKey.Serialize())

	outAddrs := make(map[string]big.Int)

	for addr, amount := range destinations {
		if curA, ok := outAddrs[addr]; !ok {
			outAddrs[addr] = *big.NewInt(int64(amount))
		} else {
			newA := big.NewInt(0)
			outAddrs[addr] = *newA.Add(&curA, big.NewInt(int64(amount)))
		}
	}

	tx := TempNewTXMultiOut(from, outAddrs)

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

// SendTX from Segwit address by lib gobcy, with preference, without manually setting fees :
// send from segwit to legacy address |
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

func (bs *BlockcypherService) GetEnpointURL() string {
	return bs.chainEndpoint
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

func (bs *BlockcypherService) CheckTx(tx string) (gobcy.TX, error) {
	return bs.chain.GetTX(tx, nil)
}

func ValidateAddress(crypto, address string) (bool, error) {
	crypto = strings.ToLower(crypto)

	var cryptoRegexMap = map[string]string{
		"btc":   "^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,62}$",
		"btg":   "^([GA])[a-zA-HJ-NP-Z0-9]{24,34}$",
		"dash":  "^([X7])[a-zA-Z0-9]{33}$",
		"dgb":   "^(D)[a-zA-Z0-9]{24,33}$",
		"eth":   "^(0x)[a-zA-Z0-9]{40}$",
		"smart": "^(S)[a-zA-Z0-9]{33}$",
		"xrp":   "^(r)[a-zA-Z0-9]{33}$",
		"zcr":   "^(Z)[a-zA-Z0-9]{33}$",
		"zec":   "^(t)[a-zA-Z0-9]{34}$",
		"ltc":   "^L[a-km-zA-HJ-NP-Z1-9]{26,33}$",
		"ltc2":  "^(ltc1|[LM])[a-zA-HJ-NP-Z0-9]{26,40}$",
		"doge":  "^D{1}[5-9A-HJ-NP-U]{1}[1-9A-HJ-NP-Za-km-z]{32}$",
		"dot":   "^(1)[a-zA-Z0-9]{47}$",
		"near":  "^(([a-z\\d]+[\\-_])*[a-z\\d]+\\.)*([a-z\\d]+[\\-_])*[a-z\\d]+$",
	}

	regex, ok := cryptoRegexMap[crypto]
	if !ok {
		return false, fmt.Errorf("Cryptocurrency not available: %s ", crypto)
	}

	re := regexp.MustCompile(regex)

	return re.MatchString(address), nil

}

func GetBTCTxStatusExtensive(txhash string, bs *BlockcypherService) (string, error) {
	var status string
	txStatus, err := bs.CheckTx(txhash)
	if err != nil {
		txInfo, err := CheckTxFromBTC(txhash)
		if err != nil {
			fmt.Printf("checkTxFromBTC err: %v", err)
			status = "Failed"
		} else {
			if txInfo.Data.Confirmations > 0 {
				status = "Success"
			} else {
				status = "Pending"
			}
		}
	} else {
		if txStatus.Confirmations > 0 {
			status = "Success"
		} else {
			status = "Pending"
		}
	}
	return status, nil
}

type BTCTxInfo struct {
	Data struct {
		BlockHeight   int    `json:"block_height"`
		BlockHash     string `json:"block_hash"`
		BlockTime     int    `json:"block_time"`
		CreatedAt     int    `json:"created_at"`
		Confirmations int    `json:"confirmations"`
		Fee           int    `json:"fee"`
		Hash          string `json:"hash"`
		InputsCount   int    `json:"inputs_count"`
		InputsValue   int    `json:"inputs_value"`
		IsCoinbase    bool   `json:"is_coinbase"`
		IsDoubleSpend bool   `json:"is_double_spend"`
		IsSwTx        bool   `json:"is_sw_tx"`
		LockTime      int    `json:"lock_time"`
		OutputsCount  int    `json:"outputs_count"`
		OutputsValue  int64  `json:"outputs_value"`
		Sigops        int    `json:"sigops"`
		Size          int    `json:"size"`
		Version       int    `json:"version"`
		Vsize         int    `json:"vsize"`
		Weight        int    `json:"weight"`
		WitnessHash   string `json:"witness_hash"`
		Inputs        *[]struct {
			PrevAddresses []string `json:"prev_addresses"`
			PrevPosition  int      `json:"prev_position"`
			PrevTxHash    string   `json:"prev_tx_hash"`
			PrevType      string   `json:"prev_type"`
			PrevValue     int      `json:"prev_value"`
			Sequence      int64    `json:"sequence"`
		} `json:"inputs"`
		Outputs *[]struct {
			Addresses         []string `json:"addresses"`
			Value             int      `json:"value"`
			Type              string   `json:"type"`
			SpentByTx         string   `json:"spent_by_tx"`
			SpentByTxPosition int      `json:"spent_by_tx_position"`
		} `json:"outputs"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrNo   int    `json:"err_no"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

var btcRateLock sync.Mutex

// TODO: 2077 add more apis
func checkTxfromQuicknode(txhash string) error {
	return nil
}

func CheckTxFromBTC(txhash string) (*BTCTxInfo, error) {
	btcRateLock.Lock()
	defer func() {
		time.Sleep(300 * time.Millisecond)
		btcRateLock.Unlock()
	}()
	url := fmt.Sprintf("https://chain.api.btc.com/v3/tx/%s?verbose=2", txhash)
	fmt.Println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}
	var result BTCTxInfo

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(string(body))
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getWalletInfo Response status != 200 " + result.Message + " " + url)
	}

	if result.Data.Hash != txhash {
		return nil, errors.New("tx not found")
	}

	return &result, nil
}
