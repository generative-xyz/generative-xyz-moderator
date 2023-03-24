package structure

type Tx struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
	Time    int64    `json:"time"`
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

type SearchInscriptionResult struct {
	Error  interface{} `json:"error"`
	Status bool        `json:"status"`
	Data   struct {
		Result []struct {
			ObjectType  string `json:"objectType"`
			Inscription struct {
				ObjectID       string      `json:"objectId"`
				InscriptionID  string      `json:"inscriptionId"`
				ContentType    string      `json:"contentType"`
				Number         int         `json:"number"`
				Sat            int64       `json:"sat"`
				Chain          string      `json:"chain"`
				Address        string      `json:"address"`
				GenesisFee     int         `json:"genesisFee"`
				GenesisHeight  int         `json:"genesisHeight"`
				Timestamp      string      `json:"timestamp"`
				ProjectName    string      `json:"projectName"`
				ProjectTokenID string      `json:"projectTokenId"`
				Buyable        bool        `json:"buyable"`
				PriceBtc       string      `json:"priceBtc"`
				Owner          interface{} `json:"owner"`
			} `json:"inscription"`
			Project  interface{} `json:"project"`
			Artist   interface{} `json:"artist"`
			TokenURI interface{} `json:"tokenUri"`
		} `json:"result"`
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		TotalPage int    `json:"totalPage"`
		Total     int    `json:"total"`
		Cursor    string `json:"cursor"`
	} `json:"data"`
}
