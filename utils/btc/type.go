package btc

import (
	"github.com/blockcypher/gobcy/v2"
	"github.com/btcsuite/btcd/chaincfg"
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

type Transaction struct {
	TxID               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

type UTXO struct {
	Address       string  `json:"address"`
	TxID          string  `json:"txid"`
	Vout          uint    `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      uint64  `json:"satoshis"`
	Height        uint64  `json:"height"`
	Confirmations int     `json:"confirmations"`
}

type Txref struct {
	TxHash        string  `json:"tx_hash"`
	BlockHeight   string  `json:"block_height"`
	TxInputN      uint    `json:"tx_input_n"`
	TxOutputN     string  `json:"tx_output_n"`
	Value         float64 `json:"value"`
	RefBalance    uint64  `json:"ref_balance"`
	Spent         uint64  `json:"spent"`
	Confirmations int     `json:"confirmations"`
	Confirmed     int     `json:"confirmed"`
	DoubleSpend   int     `json:"double_spend"`
}

type AddrInfo struct {
	Address            string  `json:"address"`
	TotalReceived      uint64  `json:"total_received"`
	TotalSent          uint64  `json:"total_sent"`
	Balance            uint64  `json:"balance"`
	UnconfirmedBalance uint64  `json:"unconfirmed_balance"`
	FinalBalance       uint64  `json:"final_balance"`
	NTx                uint    `json:"n_tx"`
	UnconfirmedNTx     uint    `json:"unconfirmed_n_tx"`
	FinalNTX           uint    `json:"final_n_tx"`
	TxRefs             []Txref `json:"txrefs"`
}

type BlockcypherService struct {
	chainEndpoint    string
	explorerEndPoint string
	bcyToken         string
	network          *chaincfg.Params
	chain            gobcy.API
}

type TxInfo struct {
	BlockHash     string `json:"block_hash"`
	BlockHeight   int    `json:"block_height"`
	DoubleSpend   bool   `json:"double_spend"`
	Confirmations int    `json:"confirmations"`
}

type Txs struct {
	Tx    string `json:"tx_hash"`
	Value string `json:"value" binding:"required"`
}

type Txo struct {
	Address string `json:"address"`
	Txs     []Txs  `json:"txrefs"`
}
