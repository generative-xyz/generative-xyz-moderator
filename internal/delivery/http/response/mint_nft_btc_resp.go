package response

type MintNftBtcResp struct {
	BaseResponse
	UserAddress    string `json:"user_address"` //user's wallet address from FE
	Amount         string `json:"amount"`
	ReceiveAddress string `json:"receiveAddress"` // address is generated from ORD service, which receive all amount
	FileURI        string `json:"fileURI"`        // FileURI will be mount if OrdAddress get all amount
	IsConfirm      bool   `json:"isConfirm"`      //default: false, if OrdAddress get all amount it will be set true
	InscriptionID  string `json:"inscriptionID"`  // nft id
	Balance        string `json:"balance"`        // balance after check

}

type MintNftBtcReceiveWalletResp struct {
	Address string `json:"address"`
	Price   string `json:"price"`
	PayType string `json:"payType"`
}
