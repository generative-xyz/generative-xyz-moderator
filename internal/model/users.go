package model

type Users struct {
	BaseModel `bson:",inline"`
	WalletAddress string `bson:"wallet_address" json:"wallet_address"`
	Message string `bson:"message" json:"message"`
	DisplayName string `bson:"display_name" json:"display_name"`
	Bio string `bson:"bio" json:"bio"`
	AvatarURL string `bson:"avatar_url" json:"avatar_url"`
}

func (m Users) CollectionName() string {
	return "users"
}
