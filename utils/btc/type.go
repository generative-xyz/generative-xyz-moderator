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
	Value uint64 `json:"value" binding:"required"`
}

type Txo struct {
	Address string `json:"address"`
	Txs     []Txs  `json:"txrefs"`
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

type QuickNodeUTXO struct {
	Version  int    `json:"version"`
	Height   int    `json:"height"`
	Value    int    `json:"value"`
	Script   string `json:"script"`
	Address  string `json:"address"`
	Coinbase bool   `json:"coinbase"`
	Hash     string `json:"hash"`
	Index    int    `json:"index"`
}

type QuickNodeTx struct {
	Result struct {
		Txid     string `json:"txid"`
		Hash     string `json:"hash"`
		Version  int    `json:"version"`
		Size     int    `json:"size"`
		Vsize    int    `json:"vsize"`
		Weight   int    `json:"weight"`
		Locktime int    `json:"locktime"`
		Vin      []struct {
			Txid      string `json:"txid"`
			Vout      int    `json:"vout"`
			ScriptSig struct {
				Asm string `json:"asm"`
				Hex string `json:"hex"`
			} `json:"scriptSig"`
			Txinwitness []string `json:"txinwitness"`
			Sequence    int64    `json:"sequence"`
		} `json:"vin"`
		Vout []struct {
			Value        float64 `json:"value"`
			N            int     `json:"n"`
			ScriptPubKey struct {
				Asm     string `json:"asm"`
				Desc    string `json:"desc"`
				Hex     string `json:"hex"`
				Address string `json:"address"`
				Type    string `json:"type"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
		Hex           string `json:"hex"`
		Blockhash     string `json:"blockhash"`
		Confirmations int    `json:"confirmations"`
		Time          int    `json:"time"`
		Blocktime     int    `json:"blocktime"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    interface{} `json:"id"`
}

type UTXOType struct {
	Value      uint64 `json:"value"`
	TxHash     string `json:"tx_hash"`
	TxOutIndex int    `json:"tx_output_n"`
}
