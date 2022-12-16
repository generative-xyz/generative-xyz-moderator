package model

type Users struct {
	BaseModel `bson:",inline"`
	WalletAddress string `bson:"wallet_address" json:"wallet_address"`
	Message string `bson:"message" json:"message"`
}

func (m Users) CollectionName() string {
	return "users"
}
