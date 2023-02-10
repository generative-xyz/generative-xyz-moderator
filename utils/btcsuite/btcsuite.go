package btcsuite

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/pkg/errors"
)

func NewBtcSuiteService(broadcastEndPoint string, username, password string, bcyToken string, env *chaincfg.Params) *BtcSuiteService {
	chain := gobcy.API{bcyToken, "btc", "main"}
	disableTLS := false

	if env.Name == chaincfg.TestNet3Params.Name {
		// For BlockCypher's internal testnet:
		chain = gobcy.API{bcyToken, "btc", "test3"}
		disableTLS = true
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         broadcastEndPoint,
		User:         username,
		Pass:         password,
		HTTPPostMode: true,       // Bitcoin core only supports HTTP POST mode
		DisableTLS:   disableTLS, // Bitcoin core does not provide TLS by default
	}

	btcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic("init rpcclient btc error")
	}

	return &BtcSuiteService{
		btcClient:         btcClient,
		broadcastEndPoint: broadcastEndPoint,
		chain:             chain,
		network:           env,
	}
}

func (bs *BtcSuiteService) NewMsgTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func (bs *BtcSuiteService) GetUTXO(from string, amount int64) ([]UTXO, error) {
	addr, err := btcutil.DecodeAddress(from, bs.network)
	if err != nil {
		return nil, err
	}
	fmt.Println("get utxo with amount = ", amount)

	listUnspentResult, err := bs.btcClient.ListUnspentMinMaxAddresses(6, 99999999, []btcutil.Address{addr})
	if err != nil {
		return nil, err
	}

	if len(listUnspentResult) <= 0 {
		return nil, err
	}

	sort.Slice(listUnspentResult, func(i, j int) bool {
		return listUnspentResult[i].Amount > listUnspentResult[j].Amount
	})

	var amountInt int64
	utxtoIndex := make(map[string]int)
	var UTXOResult []UTXO

	for i, item := range listUnspentResult {
		if item.Spendable {
			continue
		}

		utxoAmount := bs.convertToNanoBtc(item.Amount)
		amountInt += utxoAmount

		utxo := UTXO{
			TxID:          item.TxID,
			Vout:          item.Vout,
			Address:       item.Address,
			Account:       item.Account,
			ScriptPubKey:  item.ScriptPubKey,
			RedeemScript:  item.RedeemScript,
			Amount:        utxoAmount,
			Confirmations: item.Confirmations,
			Spendable:     item.Spendable,
		}

		UTXOResult = append(UTXOResult, utxo)
		utxtoIndex[item.TxID] = i
		if amountInt >= amount {
			break
		}
	}

	//add small value
	for _, item := range listUnspentResult {
		if item.Spendable {
			continue
		}

		utxoAmount := bs.convertToNanoBtc(item.Amount)
		_, ok := utxtoIndex[item.TxID]
		if utxoAmount <= 1000 && !ok {
			utxo := UTXO{
				TxID:          item.TxID,
				Vout:          item.Vout,
				Address:       item.Address,
				Account:       item.Account,
				ScriptPubKey:  item.ScriptPubKey,
				RedeemScript:  item.RedeemScript,
				Amount:        utxoAmount,
				Confirmations: item.Confirmations,
				Spendable:     item.Spendable,
			}

			UTXOResult = append(UTXOResult, utxo)
			break
		}
	}

	return UTXOResult, nil
}

func (bs *BtcSuiteService) ImportAddress(address string) error {
	err := bs.btcClient.ImportAddressRescan(address, "", true)
	if err != nil {
		return err

	}

	return nil
}

func (bs *BtcSuiteService) PushRawData(rawData string) (string, error) {
	result, err := bs.chain.PushTX(rawData)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result.Trans.Hash, nil
}

func (bs *BtcSuiteService) PushRawDataWithInternalNode(rawData string) (string, error) {
	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(rawData)
	if err != nil {
		return "", err
	}

	msgTx := &wire.MsgTx{}
	err = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return "", err
	}

	txId, err := bs.btcClient.SendRawTransaction(msgTx, true)
	if err != nil {
		return "", err
	}

	return txId.String(), nil
}

func (bs *BtcSuiteService) GetAddressLegacyAndPK(privateKey string) (string, []byte, error) {
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", nil, err
	}

	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), bs.network)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("address", addrPubKey.EncodeAddress())

	destinationAddr, err := btcutil.DecodeAddress(addrPubKey.EncodeAddress(), bs.network)
	if err != nil {
		return "", nil, err
	}
	sourcePKScript, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("sourcePKScript", hex.EncodeToString(sourcePKScript))

	return addrPubKey.EncodeAddress(), sourcePKScript, nil
}

func (bs *BtcSuiteService) getFeeLevel(feeLevel string) (int64, error) {
	blocCh, err := bs.chain.GetChain()
	if err != nil {
		return 0, err
	}

	var fee int
	if feeLevel == PreferenceHigh {
		fee = blocCh.HighFee
	} else if feeLevel == PreferenceLow {
		fee = blocCh.LowFee
	} else {
		fee = blocCh.MediumFee
	}

	return int64(fee), nil
}

func (bs *BtcSuiteService) getVSizeBTCTx(txContent string) (int64, error) {
	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(txContent)
	if err != nil {

		return 0, err
	}
	txRawResult, err := bs.btcClient.DecodeRawTransaction(serializedTx)
	if err != nil {
		return 0, err
	}

	return int64(txRawResult.Vsize), nil
}

func (bs *BtcSuiteService) GetTxDetailByTxID(txId string) (*btcjson.GetTransactionResult, error) {
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return nil, errors.Wrap(err, "chainhash.NewHashFromStr")
	}

	txResponse, err := bs.btcClient.GetTransaction(txHash)

	if err != nil {
		return nil, errors.Wrap(err, "btcClient.GetTransaction")
	}

	return txResponse, nil
}

func (bs *BtcSuiteService) EstimateFee(privKey string, destination string, amount int64, feeLevel string) (int64, error) {
	addressSender, pkScript, err := bs.GetAddressLegacyAndPK(privKey)
	if err != nil {
		return 0, err
	}

	fmt.Println(fmt.Sprintf("Get utxo with amount = %v", amount))
	utxo, err := bs.GetUTXO(addressSender, amount)
	if err != nil {
		return 0, err
	}

	_, rawData, err := bs.BuildRawData(addressSender, destination, amount, 0, privKey, pkScript, utxo)

	if err != nil {
		return 0, err
	}

	feePerKb, err := bs.getFeeLevel(feeLevel)

	if err != nil {
		return 0, err
	}

	totalKb, err := bs.getVSizeBTCTx(rawData)
	if err != nil {
		return 0, err
	}

	rs := totalKb * feePerKb
	kb := rs / 1024
	return kb, nil
}

func (bs *BtcSuiteService) CreateLegacyTx(privKey string, destination string, amount int64, txFee int64) (string, string, string, error) {
	if txFee <= 0 {
		return "", "", "", errors.New("Fees should be not empty")
	}

	addressSender, pkScript, err := bs.GetAddressLegacyAndPK(privKey)
	if err != nil {
		return "", "", "", err
	}

	fmt.Println(fmt.Sprintf("Get utxo with amount = %v", amount))
	utxo, err := bs.GetUTXO(addressSender, amount+(txFee*2))
	if err != nil {
		return "", "", "", err
	}

	utxoByte, err := json.Marshal(utxo)
	if err != nil {
		return "", "", "", err
	}

	txHash, rawData, err := bs.BuildRawData(addressSender, destination, amount, txFee, privKey, pkScript, utxo)

	if err != nil {
		return "", "", "", err
	}

	return txHash, string(utxoByte), rawData, nil
}

func (bs *BtcSuiteService) CreateLegacyTxWithTxStuck(privKey string, destination string, amount int64, txFee int64, utxoInput string, txID string) (string, string, error) {
	if txFee <= 0 {
		return "", "", errors.New("Fees should be not empty")
	}

	addressSender, pkScript, err := bs.GetAddressLegacyAndPK(privKey)
	if err != nil {
		return "", "", errors.Wrap(err, "bs.GetAddressLegacyAndPK")
	}

	txDetail, err := bs.GetTxDetailByTxID(txID)
	if err != nil {
		return "", "", errors.Wrap(err, "bs.GetTxDetailByTxID")
	}

	if txDetail.Confirmations != 0 {
		//return "", "", errors.Wrap(err, "Tx confirmations")
	}

	var utxo []UTXO
	err = json.Unmarshal([]byte(utxoInput), &utxo)
	if err != nil {
		return "", "", errors.Wrap(err, "json.Unmarshal")
	}

	return bs.BuildRawData(addressSender, destination, amount, txFee, privKey, pkScript, utxo)
}

func (bs *BtcSuiteService) convertToNanoBtc(input float64) int64 {
	rs := input * 1e8
	return int64(rs)
}

func (bs *BtcSuiteService) BuildRawData(from, to string, amount int64, totalFeeAmount int64, privKey string, pkScript []byte, txInput []UTXO) (string, string, error) {

	fmt.Println("totalFeeAmount", totalFeeAmount)

	// creating a new bitcoin transaction, different sections of the tx, including
	// input list (contain UTXOs) and outputlist (contain destination address and usually our address)
	// in next steps, sections will be field and pass to sign
	redeemTx, err := bs.NewMsgTx()
	if err != nil {
		return "", "", err
	}

	var totalInputAmount int64
	for _, input := range txInput {
		totalInputAmount += input.Amount

		utxoHash, err := chainhash.NewHashFromStr(input.TxID)
		if err != nil {
			return "", "", err
		}

		// the second argument is vout or Tx-index, which is the index
		// of spending UTXO in the transaction that Txid referred to
		// in this case is 1, but can vary different numbers
		outPoint := wire.NewOutPoint(utxoHash, uint32(input.Vout))

		// making the input, and adding it to transaction
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = 0xfffffff0
		redeemTx.AddTxIn(txIn)
	}

	if totalInputAmount <= 0 {
		return "", "", errors.New("Balance should be not empty")
	}

	if amount <= 0 {
		return "", "", errors.New("Amount should be not empty")
	}

	totalReturnAmount := totalInputAmount - amount - totalFeeAmount

	fmt.Println("totalInputAmount", totalInputAmount)
	fmt.Println("totalSendAmount", amount)
	fmt.Println("totalReturnAmount", totalReturnAmount)

	if totalReturnAmount <= 0 {
		return "", "", errors.New(fmt.Sprintf("overload balance, total = %v, send = %v, fee = %v, retrun = %v", totalInputAmount, amount, totalFeeAmount, totalReturnAmount))
	}

	//todo: [important] receiver
	destinationAddr, err := btcutil.DecodeAddress(to, bs.network)
	if err != nil {
		return "", "", err
	}
	destinationAddrByte, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return "", "", err
	}
	redeemTxOutOfReceiver := wire.NewTxOut(amount, destinationAddrByte)
	redeemTx.AddTxOut(redeemTxOutOfReceiver)

	//todo: [important] return
	senderAddr, err := btcutil.DecodeAddress(from, bs.network)
	if err != nil {
		return "", "", err
	}

	senderAddrByte, err := txscript.PayToAddrScript(senderAddr)
	if err != nil {
		return "", "", err
	}
	redeemTxOutOfSender := wire.NewTxOut(totalReturnAmount, senderAddrByte)
	redeemTx.AddTxOut(redeemTxOutOfSender)

	// now sign the transaction
	txHash, finalRawTx, err := bs.SignLegacyTx(privKey, pkScript, len(txInput), redeemTx)
	if err != nil {
		return "", "", err
	}
	return txHash, finalRawTx, nil
}

func (bs *BtcSuiteService) SignLegacyTx(privKey string, sourcePKScript []byte, totalInput int, redeemTx *wire.MsgTx) (string, string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", "", err
	}

	for i := 0; i < totalInput; i++ {
		signature, err := txscript.SignatureScript(
			redeemTx,
			i,
			sourcePKScript,
			txscript.SigHashAll,
			wif.PrivKey,
			true)

		if err != nil {
			return "", "", err
		}

		// since there is only one input, and want to add
		// signature to it use 0 as index
		redeemTx.TxIn[i].SignatureScript = signature
	}

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return redeemTx.TxHash().String(), hexSignedTx, nil
}
