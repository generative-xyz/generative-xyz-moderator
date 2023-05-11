package mempool_space

type AddressTxItemResponse struct {
	TxID     string                      `json:"txid"`
	Version  int                         `json:"version"`
	LockTime int64                       `json:"locktime"`
	Vin      []AddressTxItemResponseVin  `json:"vin"`
	Vout     []AddressTxItemResponseVout `json:"vout"`
	Size     int64                       `json:"size"`
	Weight   int64                       `json:"weight"`
	Fee      int64                       `json:"fee"`
	Status   AddressTxItemResponseStatus `json:"status"`
}

type AddressTxItemResponseStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTine   int64  `json:"block_time"`
}

type AddressTxItemResponseVout struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
	ScriptpubkeyType    string `json:"scriptpubkey_type"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

type AddressTxItemResponseVin struct {
	Prevout prevout `json:"prevout"`
}

type prevout struct {
	Scriptpubkey_address string `json:"scriptpubkey_address"`
}
