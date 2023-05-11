package mempool_space

type AddressTxItemResponse struct {
	TxID     string                      `json:"txid"`
	Version  int                         `json:"version"`
	LockTime string                      `json:"locktime"`
	Vin      interface{}                 `json:"vin"`
	Vout     interface{}                 `json:"vout"`
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
