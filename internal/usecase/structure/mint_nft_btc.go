package structure

import "time"

type MintNftBtcData struct {
	WalletAddress     string `json:"walletAddress"`
	RefundUserAddress string `json:"refundUserAddress"`
	ProjectID         string `json:"projectID"`
	PayType           string `json:"payType"`
}

type MintingInscription struct {
	ID            string     `json:"id"`
	CreatedAt     *time.Time `json:"createdAt"`
	ExpiredAt     *time.Time `json:"expiredAt"`
	StatusIndex   int        `json:"statusIndex"`
	Status        string     `json:"status"`
	FileURI       string     `json:"fileURI"`
	ProjectImage  string     `json:"projectImage"`
	ProjectID     string     `json:"projectID"`
	ProjectName   string     `json:"projectName"`
	InscriptionID string     `json:"inscriptionID"`

	ReceiveAddress string `json:"receiveAddress"`
	TxMint         string `json:"txMint"`
	TxSendNft      string `json:"txSendNft"`

	Amount  string `json:"amount"`
	PayType string `json:"payType"`

	IsCancel bool `json:"isCancel"`

	ProgressStatus interface{} `json:"progressStatus"`
}
