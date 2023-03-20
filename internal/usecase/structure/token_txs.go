package structure

type Tx struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

type Input struct {
	Address string   `json:"address"`
	Witness []string `json:"witness"`
	Value   uint64   `json:"value"`
}

type Output struct {
	Address string        `json:"address"`
	Witness []string      `json:"witness"`
	Value   uint64        `json:"value"`
	Spender OutputSpender `json:"spender"`
}

type OutputSpender struct {
	TxId  string `json:"txid"`
	Input uint64 `json:"input"`
}
