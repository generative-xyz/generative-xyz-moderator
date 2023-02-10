package response

type BctWalletResp struct {
	BaseResponse
	UserAddress string `json:"user_address"` //user's wallet address from FE
	Amount string `bson:"amount"`
	OrdAddress string `json:"ordAddress"` // address is generated from ORD service, which receive all amount
	FileURI string `json:"fileURI"` // FileURI will be mount if OrdAddress get all amount
	IsConfirm bool  `json:"isConfirm"` //default: false, if OrdAddress get all amount it will be set true
	InscriptionID string `json:"inscriptionID"` // tokenID in ETH
	Balance string `json:"balance"` // balance after check
}

type BctReceiveWalletResp struct {
	
	Address string `json:"address"` 
	Pricce string `bson:"price"`
	
}


type BtcWalletRespV2 struct {
	BaseResponse
	UserAddress string `json:"userAddress"` //user's wallet address from FE
	Amount string `bson:"amount"`
	MintFee string `bson:"mintFee"`
	SentTokenFee string `bson:"sentTokenFee"`
	OrdAddress string `json:"ordAddress"` // address is generated from ORD service, which receive all amount
	FileURI string `json:"fileURI"` // FileURI will be mount if OrdAddress get all amount
	IsConfirm bool  `json:"isConfirm"` //default: false, if OrdAddress get all amount it will be set true
	InscriptionID string `json:"inscriptionID"` // tokenID in ETH
	Balance string `json:"balance"` // balance after check
}
