package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
const DustValue = 546

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

	fmt.Printf("HHH TX after signing \n")

	for i, in := range msgTx.TxIn {
		fmt.Printf("HHH Tx in %v - %+v\n", i, in)
		// fmt.Printf("HHH Tx input %v - %+v\n", i, msgTx.TxIn[i])
	}

	for i, in := range msgTx.TxOut {
		fmt.Printf("HHH Tx out %v - %+v\n", i, in)
		// fmt.Printf("HHH Tx output %v - %+v\n", i, finalPsbt.Outputs[i])
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
		Value:      uint64(tx.TxOut[0].Value),
		TxHash:     txID,
		TxOutIndex: 0,
	}

	if len(tx.TxOut) > 1 {
		changeUTXO := UTXOType{
			Value:      uint64(tx.TxOut[1].Value),
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

func SelectUTXOs(utxos []UTXOType, price, fee uint64) ([]UTXOType, uint64, error) {
	totalPaymentAmt := price + fee
	utxos = sortUTXOs(utxos)
	totalInputAmount := uint64(0)
	resultUTXOs := []UTXOType{}

	if totalPaymentAmt > 0 {
		if len(utxos) == 0 {
			return nil, 0, errors.New("BTC balance is insufficient to create tx.")
		}
		if utxos[len(utxos)-1].Value >= totalPaymentAmt {
			// select the smallest utxo
			resultUTXOs = append(resultUTXOs, utxos[len(utxos)-1])
			totalInputAmount = utxos[len(utxos)-1].Value

		} else if utxos[0].Value < totalPaymentAmt {
			// select multiple UTXOs
			for i := 0; i < len(utxos); i++ {
				utxo := utxos[i]
				resultUTXOs = append(resultUTXOs, utxo)
				totalInputAmount += utxo.Value
				if totalInputAmount >= totalPaymentAmt {
					break
				}
			}
			if totalInputAmount < totalPaymentAmt {
				return nil, 0, errors.New("BTC balance is insufficient to create tx.")
			}
		} else {
			// select the nearest UTXO
			selectedUTXO := utxos[0]
			for i := 1; i < len(utxos); i++ {
				if utxos[i].Value < totalPaymentAmt {
					resultUTXOs = append(resultUTXOs, selectedUTXO)
					totalInputAmount = selectedUTXO.Value
					break
				}

				selectedUTXO = utxos[i]
			}
		}
	}
	return resultUTXOs, totalInputAmount, nil
}

type CreateTxBuyResp struct {
	SplitTxID  string
	SplitTxHex string
	SplitTxFee uint64
	SplitUTXOs []UTXO
	TxID       string
	TxHex      string
	BuyTxFee   uint64
}

func CreatePSBTToBuyInscription(
	sellerSignedPsbtB64 string,
	privateKey string,
	address string,
	receiverInscriptionAddress string,
	price uint64,
	utxos []UTXOType,
	feeRatePerByte uint64,
	maxFee uint64,
) (*CreateTxBuyResp, error) {

	// parse Seller Tx and validate
	sellerPsbt, err := ParsePSBTFromBase64(sellerSignedPsbtB64)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] ParsePSBTFromBase64 err: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] ParsePSBTFromBase64 err: %v", err)
	}

	sellerInputs := sellerPsbt.Inputs
	sellerOutputs := sellerPsbt.Outputs
	if len(sellerInputs) == 0 {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return nil, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}
	if len(sellerInputs) != len(sellerOutputs) {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return nil, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}
	if sellerInputs[0].WitnessUtxo == nil {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
		return nil, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT.")
	}

	// extract value in inscription
	valueInscription := sellerInputs[0].WitnessUtxo.Value
	if valueInscription <= 0 {
		log.Println("[CreatePSBTToBuyInscription] Invalid seller's PSBT - value inscription is zero.")
		return nil, errors.New("[CreatePSBTToBuyInscription] Invalid seller's PSBT - value inscription is zero.")
	}

	// filter pending UTXOs
	_, spendableUTXOs, err := FilterPendingUTXOs(utxos, address)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error filter pending utxos ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error filter pending utxos %v", err)
	}

	// create dummy UTXO
	splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs, err := createDummyOutputs(privateKey, address, spendableUTXOs, uint32(feeRatePerByte))
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] create dummy utxo err: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] create dummy utxo err: %v", err)
	}
	if splitFee > 0 {
		if maxFee <= splitFee {
			log.Println("[CreatePSBTToBuyInscription] invalid max fee: ", maxFee)
			return nil, fmt.Errorf("[CreatePSBTToBuyInscription] invalid max fee: %v", maxFee)
		}
		maxFee -= splitFee
	}

	fmt.Println("splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs: ", splitTxID, splitTxHex, splitFee, dummyUTXO, newUTXOs)

	// ====== create psbt to buy  ======

	// decode receiver address to PubkeyScript
	receiverAddr, err := btcutil.DecodeAddress(receiverInscriptionAddress, BTCChainConf)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when decoding receiver address: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when decoding receiver address: %v", err)
	}
	receiverPKScript, err := txscript.PayToAddrScript(receiverAddr)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when new receiver PKScript: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when new receiver PKScript: %v", err)
	}

	// decode sender address to PubkeyScript
	senderAddr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when new sender PKScript: ", err)
		return nil, err
	}
	senderPKScript, _ := txscript.PayToAddrScript(senderAddr)

	// get private key from string
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, err
	}
	btcPrivateKey := wif.PrivKey

	// select UTXOs to create tx
	selectedUTXOs, totalInputAmount, err := SelectUTXOs(newUTXOs, price, maxFee)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when selecting utxos to create tx buy: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when selecting utxos to create tx buy: %v", err)
	}

	inputs := []*wire.OutPoint{}
	outputs := []*wire.TxOut{}
	witnessInputs := []*wire.TxOut{}
	nSequences := []uint32{}
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)

	// Frist in - first out: dummy utxo & receiver inscription
	// the first output coin has value equal to the sum of dummy value and value inscription
	// this makes sure the first output coin is inscription outcoin
	// receiver inscription must be user's address
	dummyUtxoHash, err := chainhash.NewHashFromStr(dummyUTXO.TxHash)
	if err != nil {
		return nil, err
	}
	in0 := wire.NewOutPoint(dummyUtxoHash, uint32(dummyUTXO.TxOutIndex))
	out0 := wire.NewTxOut(int64(dummyUTXO.Value)+valueInscription, receiverPKScript)
	preOutIn := &wire.TxOut{
		Value:    int64(dummyUTXO.Value),
		PkScript: senderPKScript,
	}

	inputs = append(inputs, in0)
	witnessInputs = append(witnessInputs, preOutIn)
	prevOuts.AddPrevOut(*in0, preOutIn)
	nSequences = append(nSequences, uint32(feeRatePerByte))
	outputs = append(outputs, out0)

	// Add seller signed inputs and outputs
	for i := 0; i < len(sellerInputs); i++ {
		in := sellerPsbt.UnsignedTx.TxIn[i].PreviousOutPoint
		inputs = append(inputs, &in)
		witnessInputs = append(witnessInputs, sellerInputs[i].WitnessUtxo)
		prevOuts.AddPrevOut(in, sellerInputs[i].WitnessUtxo)
		nSequences = append(nSequences, sellerPsbt.UnsignedTx.TxIn[i].Sequence)

		outputs = append(outputs, sellerPsbt.UnsignedTx.TxOut[i])
	}

	// Add payment utxo inputs
	for _, utxo := range selectedUTXOs {
		utxoHash, err := chainhash.NewHashFromStr(utxo.TxHash)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(utxoHash, uint32(utxo.TxOutIndex))
		preOutIn := &wire.TxOut{
			Value:    int64(utxo.Value),
			PkScript: senderPKScript,
		}
		inputs = append(inputs, outPoint)
		witnessInputs = append(witnessInputs, preOutIn)
		prevOuts.AddPrevOut(*outPoint, preOutIn)
		nSequences = append(nSequences, uint32(feeRatePerByte))
	}

	// calcalate network fee
	fee := EstimateTxFee(uint(len(inputs)), uint(len(outputs)), uint(feeRatePerByte))
	// max fee can paid is defined from users in advance
	if fee > maxFee {
		fee = maxFee
	}

	// create change output
	changeAmount := totalInputAmount - price - fee
	if changeAmount > 0 {
		if changeAmount >= DustValue {
			out := wire.NewTxOut(int64(changeAmount), senderPKScript)
			outputs = append(outputs, out)
		} else {
			fee += changeAmount
		}
	}

	// init psbt from inputs and outputs
	finalPsbt, err := psbt.New(inputs, outputs, wire.TxVersion, 0, nSequences)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when new Psbt: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when new Psbt: %v", err)
	}

	updater, err := psbt.NewUpdater(finalPsbt)
	if err != nil {
		log.Println("[CreatePSBTToBuyInscription] Error when new updater PSBT: ", err)
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when new updater PSBT: %v", err)
	}
	// add witnesss for input coins
	for i := range finalPsbt.Inputs {
		err := updater.AddInWitnessUtxo(witnessInputs[i], i)
		if err != nil {
			log.Println("[CreatePSBTToBuyInscription] Error when add witness input index: ", i, err)
			return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when add witness input index %v: %v", i, err)
		}
		// if i > 0 && i <= len(sellerInputs) {
		// 	updater.AddInSighashType(txscript.SigHashSingle|txscript.SigHashAnyOneCanPay, i)
		// }
	}

	fmt.Printf("finalPsbt : %+v\n", finalPsbt)

	// todo: remove
	for i, in := range finalPsbt.UnsignedTx.TxIn {
		fmt.Printf("HHH Tx in %v - %+v\n", i, in)
		fmt.Printf("HHH Tx input %v - %+v\n", i, finalPsbt.Inputs[i])
	}
	for i, in := range finalPsbt.UnsignedTx.TxOut {
		fmt.Printf("HHH Tx out %v - %+v\n", i, in)
		fmt.Printf("HHH Tx output %v - %+v\n", i, finalPsbt.Outputs[i])
	}

	// sign tx
	sigs := []wire.TxWitness{}
	for i, input := range finalPsbt.Inputs {
		if i == 0 || i > len(sellerInputs) {
			sig, err := txscript.WitnessSignature(finalPsbt.UnsignedTx, txscript.NewTxSigHashes(finalPsbt.UnsignedTx, prevOuts), i, int64(input.WitnessUtxo.Value), senderPKScript, txscript.SigHashAll, btcPrivateKey, true)
			if err != nil {
				log.Println("[CreatePSBTToBuyInscription] Error when signing on raw btc tx: ", err)
				return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when signing on raw btc tx: %v", err)
			}
			fmt.Printf("sig: %+v\n", sig)

			// msgTx.TxIn[i].Witness = [][]byte{sig, sourcePKScript}
			sigs = append(sigs, sig)
		}
	}

	indexSig := 0
	indexSellerSig := 0
	for i := range finalPsbt.Inputs {
		if i == 0 || i > len(sellerInputs) {
			signSuccess, err := updater.Sign(i, sigs[indexSig][0], sigs[indexSig][1], nil, nil)
			if err != nil {
				log.Println("[CreatePSBTToBuyInscription] Sign index error: ", i, err)
				return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Sign index %v error: %v", i, err)
			}
			if signSuccess != psbt.SignSuccesful {
				log.Println("[CreatePSBTToBuyInscription] Sign index not success: ", i, err)
				return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Sign index %v not success: %v", i, err)
			}
			indexSig++
		} else {
			// signSuccess, err = updater.Sign(i, sellerInputs[indexSellerSig].FinalScriptWitness, nil, nil, nil)
			finalPsbt.Inputs[i].FinalScriptWitness = sellerInputs[indexSellerSig].FinalScriptWitness
			indexSellerSig++
		}

	}

	// todo: remove
	fmt.Printf("HHH PSBT after signing \n")
	for i, in := range finalPsbt.UnsignedTx.TxIn {
		fmt.Printf("HHH Tx in %v - %+v\n", i, in)
		fmt.Printf("HHH Tx input %v - %+v\n", i, finalPsbt.Inputs[i])
	}
	for i, in := range finalPsbt.UnsignedTx.TxOut {
		fmt.Printf("HHH Tx out %v - %+v\n", i, in)
		fmt.Printf("HHH Tx output %v - %+v\n", i, finalPsbt.Outputs[i])
	}

	// finalize tx
	for i := range finalPsbt.Inputs {
		if i == 0 || i > len(sellerInputs) {
			err := psbt.Finalize(finalPsbt, i)
			if err != nil {
				log.Println("[CreatePSBTToBuyInscription] Finalize Psbt index error : ", i, err)
				return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Finalize Psbt index %v error : %v", i, err)
			}
		}
	}

	fmt.Printf("finalPsbt.IsComplete() : %v\n", finalPsbt.IsComplete())

	msgTx, err := psbt.Extract(finalPsbt)
	if err != nil {
		return nil, fmt.Errorf("[CreatePSBTToBuyInscription] Error when extract final psbt: %v", err)
	}

	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	txHex := hex.EncodeToString(buf.Bytes())
	txID := msgTx.TxHash()

	resp := &CreateTxBuyResp{
		SplitTxID:  splitTxID,
		SplitTxHex: splitTxHex,
		SplitTxFee: splitFee,
		TxID:       txID.String(),
		TxHex:      txHex,
		BuyTxFee:   fee,
	}

	return resp, nil
}

type CreatePSBTToBuyInscriptionRequest struct {
	Psbt     string     `json:"sellerSignedPsbtB64"`
	Receiver string     `json:"receiverInscriptionAddress"`
	Price    uint64     `json:"price"`
	FeeRate  uint64     `json:"feeRatePerByte"`
	UTXOs    []UTXOType `json:"utxos"`
}
type CreatePSBTToBuyInscriptionMultiRequest struct {
	BuyReqInfos []BuyReqInfo `json:"buyReqInfo"`
	// Psbt     []string   `json:"sellerSignedPsbtB64"`
	// Receiver string     `json:"receiverInscriptionAddress"`
	// Price    uint64     `json:"price"`
	FeeRate uint64     `json:"feeRatePerByte"`
	UTXOs   []UTXOType `json:"utxos"`
}

type CreatePSBTToBuyInscriptionRespond struct {
	TxID          string     `json:"txID"`
	TxHex         string     `json:"txHex"`
	Fee           int        `json:"fee"`
	SelectedUTXOs []UTXOType `json:"selectedUTXOs"`
	SplitTxID     string     `json:"splitTxID"`
	SplitUTXOs    []UTXOType `json:"splitUTXOs"`
	SplitTxRaw    string     `json:"splitTxRaw"`
}

func CreatePSBTToBuyInscriptionViaAPI(
	endpoint string,
	address string,
	sellerSignedPsbtB64 string,
	receiverInscriptionAddress string,
	price uint64,
	utxos []UTXOType,
	feeRatePerByte uint64,
	maxFee uint64,
) (*CreatePSBTToBuyInscriptionRespond, error) {
	_, spendableUTXOs, err := FilterPendingUTXOs(utxos, address)
	if err != nil {
		return nil, err
	}

	data := CreatePSBTToBuyInscriptionRequest{
		UTXOs:    spendableUTXOs,
		Psbt:     sellerSignedPsbtB64,
		Receiver: receiverInscriptionAddress,
		Price:    price,
		FeeRate:  feeRatePerByte,
	}

	json_data, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpoint+"/api/createtxbuy", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		return nil, err
	}

	var res CreatePSBTToBuyInscriptionRespond
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if res.TxHex == "" {
		var resErr struct {
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&resErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(resErr.Message)
	}
	return &res, nil
}

func SendTxBlockStream(txraw string) error {
	resp, err := http.Post("https://blockstream.info/api/tx", "application/json",
		bytes.NewBuffer([]byte(txraw)))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	if strings.Contains(bodyStr, "RPC error") {
		return errors.New(bodyStr)
	}
	return nil
}

type BuyReqInfo struct {
	SellerSignedPsbtB64        string `json:"sellerSignedPsbtB64"`
	ReceiverInscriptionAddress string `json:"receiverInscriptionAddress"`
	Price                      int    `json:"price"`
}

func CreatePSBTToBuyInscriptionMultiViaAPI(
	endpoint string,
	address string,
	// sellerSignedPsbtB64 []string,
	// receiverInscriptionAddress string,
	// price uint64,
	buyReqInfos []BuyReqInfo,
	utxos []UTXOType,
	feeRatePerByte uint64,
	// maxFee uint64,
) (*CreatePSBTToBuyInscriptionRespond, error) {
	_, spendableUTXOs, err := FilterPendingUTXOs(utxos, address)
	if err != nil {
		return nil, err
	}

	data := CreatePSBTToBuyInscriptionMultiRequest{
		BuyReqInfos: buyReqInfos,
		UTXOs:       spendableUTXOs,
		// Psbt:     sellerSignedPsbtB64,
		// Receiver: receiverInscriptionAddress,
		// Price:    price,
		FeeRate: feeRatePerByte,
	}

	json_data, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}
	resp, err := http.Post(endpoint+"/api/createtxbuymulti", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		return nil, err
	}

	log.Println("CreatePSBTToBuyInscriptionMultiViaAPI payload", string(json_data))
	var res CreatePSBTToBuyInscriptionRespond
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if res.TxHex == "" {
		var resErr struct {
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&resErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(resErr.Message)
	}
	return &res, nil
}
