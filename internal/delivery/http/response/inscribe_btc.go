package response

type InscribeBtcResp struct {
	BaseResponse
	UserAddress   string `json:"userAddress"` //user's wallet address from FE
	Amount        string `json:"amount"`
	MintFee       string `json:"mintFee"`
	SentTokenFee  string `json:"sentTokenFee"`
	OrdAddress    string `json:"ordAddress"` // address is generated from ORD service, which receive all amount
	SegwitAddress string `json:"segwitAddress"`
	FileURI       string `json:"fileURI"`       // FileURI will be mount if OrdAddress get all amount
	IsConfirm     bool   `json:"isConfirm"`     //default: false, if OrdAddress get all amount it will be set true
	InscriptionID string `json:"inscriptionID"` // tokenID in ETH
	Balance       string `json:"balance"`       // balance after check
	TimeoutAt     string `json:"timeoutAt"`

	ID string `json:"id"`
}

type InscribeInfoResp struct {
	BaseResponse
	Address            string `json:"address"`
	Index              string `json:"index"`
	OutputValue        string `json:"outputValue"`
	Sat                string `json:"sat"`
	Preview            string `json:"preview"`
	Content            string `json:"content"`
	ContentLength      string `json:"contentLength"`
	ContentType        string `json:"contentType"`
	Timestamp          string `json:"timestamp"`
	GenesisHeight      string `json:"genesisHeight"`
	GenesisFee         string `json:"genesisFee"`
	GenesisTransaction string `json:"genesisTransaction"`
	Location           string `json:"location"`
	Output             string `json:"output"`
	Offset             string `json:"offset"`
}
