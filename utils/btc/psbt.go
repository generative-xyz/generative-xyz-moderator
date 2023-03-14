package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

var BTCChainConf = &chaincfg.MainNetParams

const DummyValue = 1000

func ParsePSBTFromBase64(data string) (*psbt.Packet, error) {
	psbtTx, err := psbt.NewFromRawBytes(
		strings.NewReader(data), true,
	)
	if err != nil {
		return nil, err
	}
	return psbtTx, nil
}

func ParseTx(data string) (*wire.MsgTx, error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var tx wire.MsgTx
	err = tx.Deserialize(strings.NewReader(string(dataBytes)))
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func GetAddressFromPKScript(script []byte) (string, error) {
	pk1, err := txscript.ParsePkScript(script)
	if err != nil {
		return "", err
	}
	address, err := pk1.Address(&chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// TODO: update size here
func EstimateTxFee(numIns, numOuts, feeRatePerByte uint) uint64 {
	return (uint64(numIns)*68 + uint64(numOuts)*43) * uint64(feeRatePerByte)
}

// NOTE: Only support creating tx from BTC segwit
// CreateZeroValueOutputs creates tx that use uxtos are inputs and create new zero-value outputs
func CreateZeroValueOutputs(privateKey string, receiverAddress string, numOuts uint, utxos []UTXOType, feeRatePerByte uint32) (string, string, uint64, error) {
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// decode receiver address to PubkeyScript
	decodedAddr, err := btcutil.DecodeAddress(receiverAddress, BTCChainConf)
	if err != nil {
		return "", "", 0, fmt.Errorf("Error when decoding receiver address: %v", err)
	}
	destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
	if err != nil {
		return "", "", 0, fmt.Errorf("Error when new Address Script: %v", err)
	}

	// add TxIns into raw tx
	// totalInputAmount in external unit
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	totalInputAmount := uint64(0)
	for _, in := range utxos {
		utxoHash, err := chainhash.NewHashFromStr(in.TxHash)
		if err != nil {
			return "", "", 0, err
		}
		outPoint := wire.NewOutPoint(utxoHash, uint32(in.TxOutIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = feeRatePerByte
		msgTx.AddTxIn(txIn)
		totalInputAmount += uint64(in.Value)
		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{})
	}

	for i := 0; i < int(numOuts); i++ {
		// TODO: add op_return

		// create OP_RETURN script
		script, err := txscript.NullDataScript(destinationAddrByte)
		if err != nil {
			return "", "", 0, err
		}
		redeemTxOut := wire.NewTxOut(0, script)

		msgTx.AddTxOut(redeemTxOut)
	}

	// add change output
	fee := EstimateTxFee(uint(len(utxos)), numOuts+1, uint(feeRatePerByte))
	changeAmt := totalInputAmount - fee
	redeemTxOut := wire.NewTxOut(int64(changeAmt), destinationAddrByte)
	msgTx.AddTxOut(redeemTxOut)

	// sign tx
	msgTx, err = signTx(msgTx, privateKey, utxos, receiverAddress, prevOuts)
	if err != nil {
		return "", "", 0, err
	}

	var signedTx bytes.Buffer
	err = msgTx.Serialize(&signedTx)
	if err != nil {
		return "", "", 0, err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return msgTx.TxHash().String(), hexSignedTx, fee, nil

}

type PaymentInfo struct {
	Address string
	Amount  uint64
}

// NOTE: Only support creating tx from BTC segwit
// CreateZeroValueOutputs creates tx that use uxtos are inputs and create new zero-value outputs
func CreateTx(privateKey string, senderAddress string, utxos []UTXOType, paymentInfos []PaymentInfo, feeRatePerByte uint32) (string, string, uint64, error) {
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// decode receiver address to PubkeyScript
	senderDecodedAddr, err := btcutil.DecodeAddress(senderAddress, BTCChainConf)
	if err != nil {
		return "", "", 0, fmt.Errorf("Error when decoding receiver address: %v", err)
	}
	senderAddrByte, err := txscript.PayToAddrScript(senderDecodedAddr)
	if err != nil {
		return "", "", 0, fmt.Errorf("Error when new Address Script: %v", err)
	}

	// add TxIns into raw tx
	// totalInputAmount in external unit
	totalInputAmount := uint64(0)
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	for _, in := range utxos {
		utxoHash, err := chainhash.NewHashFromStr(in.TxHash)
		if err != nil {
			return "", "", 0, err
		}
		outPoint := wire.NewOutPoint(utxoHash, uint32(in.TxOutIndex))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = feeRatePerByte
		msgTx.AddTxIn(txIn)
		totalInputAmount += uint64(in.Value)
		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{})
	}

	totalOutputAmount := uint64(0)
	for _, out := range paymentInfos {
		// decode receiver address to PubkeyScript
		decodedAddr, err := btcutil.DecodeAddress(out.Address, BTCChainConf)
		if err != nil {
			return "", "", 0, fmt.Errorf("Error when decoding receiver address: %v", err)
		}
		destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
		if err != nil {
			return "", "", 0, fmt.Errorf("Error when new Address Script: %v", err)
		}

		redeemTxOut := wire.NewTxOut(int64(out.Amount), destinationAddrByte)
		msgTx.AddTxOut(redeemTxOut)

		totalOutputAmount += out.Amount
	}

	// add change output
	fee := EstimateTxFee(uint(len(utxos)), uint(len(paymentInfos)+1), uint(feeRatePerByte))
	changeAmt := totalInputAmount - totalOutputAmount - fee
	redeemTxOut := wire.NewTxOut(int64(changeAmt), senderAddrByte)
	msgTx.AddTxOut(redeemTxOut)

	// sign tx
	msgTx, err = signTx(msgTx, privateKey, utxos, senderAddress, prevOuts)
	if err != nil {
		return "", "", 0, err
	}

	var signedTx bytes.Buffer
	err = msgTx.Serialize(&signedTx)
	if err != nil {
		return "", "", 0, err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return msgTx.TxHash().String(), hexSignedTx, fee, nil

}

// sign raw tx by Segwit address (P2WPKH)
func signTx(msgTx *wire.MsgTx, privKey string, inputs []UTXOType, sourceAddressStr string, prevOuts *txscript.MultiPrevOutFetcher) (*wire.MsgTx, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, err
	}
	btcPrivateKey := wif.PrivKey
	if len(inputs) != len(msgTx.TxIn) {
		return nil, errors.New("Invalid inputs")
	}

	sourceAddress, err := btcutil.DecodeAddress(sourceAddressStr, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	sourcePKScript, _ := txscript.PayToAddrScript(sourceAddress)

	sigs := []wire.TxWitness{}

	for i, _ := range msgTx.TxIn {
		// sourcePKScript, err := hex.DecodeString(inputs[i].ScriptPubKey)
		// if err != nil {
		// 	return nil, err
		// }

		// sig, err := txscript.RawTxInWitnessSignature(msgTx, txscript.NewTxSigHashes(msgTx, prevOuts), i, int64(inputs[i].Value), sourcePKScript, txscript.SigHashAll, btcPrivateKey)
		// sig, err := txscript.SignatureScript(msgTx, i, sourcePKScript, txscript.SigHashAll, btcPrivateKey, true)
		sig, err := txscript.WitnessSignature(msgTx, txscript.NewTxSigHashes(msgTx, prevOuts), i, int64(inputs[i].Value), sourcePKScript, txscript.SigHashAll, btcPrivateKey, true)
		if err != nil {
			return nil, fmt.Errorf("Error when signing on raw btc tx: %v", err)
		}

		// msgTx.TxIn[i].Witness = [][]byte{sig, sourcePKScript}
		sigs = append(sigs, sig)

		// msgTx.TxIn[i].SignatureScript = sig

	}

	for i, _ := range msgTx.TxIn {
		msgTx.TxIn[i].Witness = sigs[i]
	}

	return msgTx, nil
}

// sortUTXOs sorts list of utxos descending by value
func sortUTXOs(utxos []UTXOType) []UTXOType {
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	return utxos
}

func createDummyOutputs(
	privateKey string, senderAddress string, utxos []UTXOType, feeRatePerByte uint32,
) (txID string, txHex string, fee uint64, dummyUTXO UTXOType, newUTXOs []UTXOType, err error) {

	utxos = sortUTXOs(utxos)
	if len(utxos) == 0 {
		err = errors.New("List of utxos is empty")
		return
	}
	// the smallest UTXO is dummy value, no need to split UTXO
	if utxos[len(utxos)-1].Value <= DummyValue {
		dummyUTXO = utxos[len(utxos)-1]
		newUTXOs = utxos[0 : len(utxos)-1]
		return
	}

	// otherwise, split the smallest UTXO
	UTXOsToSplit := []UTXOType{utxos[len(utxos)-1]}
	newUTXOs = utxos[0 : len(utxos)-1]

	paymentInfos := []PaymentInfo{{
		Address: senderAddress,
		Amount:  DummyValue,
	}}
	txID, txHex, fee, err = CreateTx(privateKey, senderAddress, UTXOsToSplit, paymentInfos, feeRatePerByte)
	if err != nil {
		log.Println("[createDummyOutputs] create tx error: ", err)
		return
	}

	// parse tx hex
	tx, err := ParseTx(txHex)
	if err != nil {
		log.Println("[createDummyOutputs] parse tx error: ", err)
		return
	}

	dummyUTXO = UTXOType{
		Value:      int(tx.TxOut[0].Value),
		TxHash:     txID,
		TxOutIndex: 0,
	}

	if len(tx.TxOut) > 1 {
		changeUTXO := UTXOType{
			Value:      int(tx.TxOut[1].Value),
			TxHash:     txID,
			TxOutIndex: 1,
		}
		newUTXOs = append(newUTXOs, changeUTXO)
	}

	return
}

// NOTE: Only support creating tx from BTC segwit
// CreateZeroValueOutputs creates tx that use uxtos are inputs and create new zero-value outputs
// func CreateDummyOutputs(privateKey string, receiverAddress string, numOuts uint, utxos []QuickNodeUTXO, feeRatePerByte uint32) (string, string, uint64, error) {
// 	msgTx := wire.NewMsgTx(wire.TxVersion)

// 	// decode receiver address to PubkeyScript
// 	decodedAddr, err := btcutil.DecodeAddress(receiverAddress, BTCChainConf)
// 	if err != nil {
// 		return "", "", 0, fmt.Errorf("Error when decoding receiver address: %v", err)
// 	}
// 	destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
// 	if err != nil {
// 		return "", "", 0, fmt.Errorf("Error when new Address Script: %v", err)
// 	}

// 	// add TxIns into raw tx
// 	// totalInputAmount in external unit
// 	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
// 	totalInputAmount := uint64(0)
// 	for _, in := range utxos {
// 		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
// 		if err != nil {
// 			return "", "", 0, err
// 		}
// 		outPoint := wire.NewOutPoint(utxoHash, uint32(in.Index))
// 		txIn := wire.NewTxIn(outPoint, nil, nil)
// 		txIn.Sequence = feeRatePerByte
// 		msgTx.AddTxIn(txIn)
// 		totalInputAmount += uint64(in.Value)
// 		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{})
// 	}

// 	for i := 0; i < int(numOuts); i++ {
// 		// TODO: add op_return

// 		// create OP_RETURN script
// 		script, err := txscript.NullDataScript(destinationAddrByte)
// 		if err != nil {
// 			return "", "", 0, err
// 		}
// 		redeemTxOut := wire.NewTxOut(0, script)

// 		msgTx.AddTxOut(redeemTxOut)
// 	}

// 	// add change output
// 	fee := EstimateTxFee(uint(len(utxos)), numOuts+1, uint(feeRatePerByte))
// 	changeAmt := totalInputAmount - fee
// 	redeemTxOut := wire.NewTxOut(int64(changeAmt), destinationAddrByte)
// 	msgTx.AddTxOut(redeemTxOut)

// 	// sign tx
// 	msgTx, err = signTx(msgTx, privateKey, utxos, receiverAddress, prevOuts)
// 	if err != nil {
// 		return "", "", 0, err
// 	}

// 	var signedTx bytes.Buffer
// 	err = msgTx.Serialize(&signedTx)
// 	if err != nil {
// 		return "", "", 0, err
// 	}

// 	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

// 	return msgTx.TxHash().String(), hexSignedTx, fee, nil

// }

func CreatePSBTToBuyInscription(
	sellerSignedPsbtB64 string,
	privateKey string,
	address string,
	receiverInscriptionAddress string,
	price uint64,
	utxos []UTXOType,
	feeRatePerByte uint64,
) (string, uint64, error) {

	// parse Seller Tx and validate
	psbt, err := ParsePSBTFromBase64(sellerSignedPsbtB64)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] ParsePSBTFromBase64 err: ", err)
		return "", 0, fmt.Errorf("[CreatePSBTToBuyInscription] ParsePSBTFromBase64 err: ", err)
	}

	sellerInputs := psbt.Inputs
	sellerOutputs := psbt.Outputs

	if len(sellerInputs) == 0 {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return "", 0, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}

	if len(sellerInputs) != len(sellerOutputs) {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return "", 0, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}
	if sellerInputs[0].WitnessUtxo == nil {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return "", 0, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}

	valueInscription := sellerInputs[0].WitnessUtxo.Value

	if valueInscription == 0 {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT - value inscription is zero.")
		return "", 0, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT - value inscription is zero.")
	}

	// create dummy UTXO
	splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs, err := createDummyOutputs(privateKey, address, utxos, uint32(feeRatePerByte))
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] create dummy utxo err: ", err)
		return "", 0, fmt.Errorf("[CreatePSBTToBuyInscription] create dummy utxo err: ", err)
	}

	fmt.Println("splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs: ", splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs)

	// TODO: create psbt to buy

	return "", 0, nil

}
