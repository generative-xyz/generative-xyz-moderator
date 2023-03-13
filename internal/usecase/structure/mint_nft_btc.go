package structure

import "time"

type MintNftBtcData struct {
	UserAddress       string `json:"userAddress"`
	WalletAddress     string `json:"walletAddress"`
	RefundUserAddress string `json:"refundUserAddress"`
	ProjectID         string `json:"projectID"`
	PayType           string `json:"payType"`

	UserID string `json:"userID"`

	Quantity int `json:"quantity"`

	FeeRate int32 `json:"feeRate"`
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
	ArtistName    string     `json:"artist_name"`
	ArtistID      string     `json:"artist_id"`
	InscriptionID string     `json:"inscriptionID"`

	ReceiveAddress string `json:"receiveAddress"`
	TxMint         string `json:"txMint"`
	TxSendNft      string `json:"txSendNft"`

	Amount  string `json:"amount"`
	PayType string `json:"payType"`

	OriginUserAddress string `json:"userWallet"`

	IsCancel bool `json:"isCancel"`

	ProgressStatus interface{} `json:"progressStatus"`

	UserID string `json:"userID"`

	Quantity int `json:"quantity"`
}
