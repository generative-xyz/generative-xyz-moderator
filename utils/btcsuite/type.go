package btcsuite

import (
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

const (
	// This will calculate and include appropriate fees for your transaction to be included in the next 1-2 blocks
	PreferenceHigh = "high"
	// This will calculate and include appropriate fees for your transaction to be included in the next 3-6 blocks
	PreferenceMedium = "medium"
	// This will calculate and include appropriate fees for your transaction to be included in the next 7 or more blocks
	PreferenceLow = "low"
	// No fee
	PreferenceZero = "zero"
)

type BtcSuiteService struct {
	chain             gobcy.API
	btcClient         *rpcclient.Client
	network           *chaincfg.Params
	broadcastEndPoint string
}

type UTXO struct {
	TxID          string `json:"txid"`
	Vout          uint32 `json:"vout"`
	Address       string `json:"address"`
	Account       string `json:"account"`
	ScriptPubKey  string `json:"scriptPubKey"`
	RedeemScript  string `json:"redeemScript,omitempty"`
	Amount        int64  `json:"amount"`
	Confirmations int64  `json:"confirmations"`
	Spendable     bool   `json:"spendable"`
}
