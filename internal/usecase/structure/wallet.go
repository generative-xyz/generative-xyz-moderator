package structure

import "time"

type InscriptionOrdInfoByOutput struct {
	Outpoint string      `json:"outpoint"`
	List     interface{} `json:"list"`
	Chain    string      `json:"chain"`
	Output   struct {
		Value        int    `json:"value"`
		ScriptPubkey string `json:"script_pubkey"`
	} `json:"output"`
	Inscriptions []string `json:"inscriptions"`
}

type InscriptionOrdInfoByID struct {
	Chain         string `json:"chain"`
	GenesisFee    int    `json:"genesis_fee"`
	GenesisHeight int    `json:"genesis_height"`
	ContentType   string `json:"content_type"`
	InscriptionID string `json:"inscription_id"`
	Next          string `json:"next"`
	Number        int    `json:"number"`
	Output        struct {
		Value        int    `json:"value"`
		ScriptPubkey string `json:"script_pubkey"`
	} `json:"output"`
	Previous  string      `json:"previous"`
	Sat       interface{} `json:"sat"`
	Satpoint  string      `json:"satpoint"`
	Timestamp string      `json:"timestamp"`
}

type BlockCypherWalletInfo struct {
	Address            string `json:"address"`
	TotalReceived      int    `json:"total_received"`
	TotalSent          int    `json:"total_sent"`
	Balance            int    `json:"balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	FinalBalance       int    `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
	Txrefs             []struct {
		TxHash        string    `json:"tx_hash"`
		BlockHeight   int       `json:"block_height"`
		TxInputN      int       `json:"tx_input_n"`
		TxOutputN     int       `json:"tx_output_n"`
		Value         int       `json:"value"`
		RefBalance    int       `json:"ref_balance"`
		Spent         bool      `json:"spent"`
		Confirmations int       `json:"confirmations"`
		Confirmed     time.Time `json:"confirmed"`
		DoubleSpend   bool      `json:"double_spend"`
	} `json:"txrefs"`
	TxURL string `json:"tx_url"`
	Error string `json:"error"`
}

type WalletInfo struct {
	BlockCypherWalletInfo
	Inscriptions          []WalletInscriptionInfo `json:"inscriptions"`
	InscriptionsByOutputs map[string][]string     `json:inscriptions_by_outputs"`
}

type WalletInscriptionInfo struct {
	InscriptionID string `json:"inscription_id"`
	Offset        int64  `json:"offset"`
	Number        int    `json:"number"`
	ContentType   string `json:"content_type"`
	ProjectID     string `json:"project_id"`
	ProjecName    string `json:"projec_name"`
	Thumbnail     string `json:"thumbnail"`
}
