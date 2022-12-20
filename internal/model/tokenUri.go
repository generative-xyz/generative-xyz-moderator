package model

type TokenUri struct {
	BaseModel `bson:",inline"`
	TokenID string `bson:"token_id" json:"token_id"`
	ContractAddress string `bson:"contract_address" json:"contract_address"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image string `bson:"image" json:"image"`
	AnimationURL string `bson:"animation_url" json:"animation_url"`
	Attributes string `bson:"attributes" json:"attributes"`
}

func (m TokenUri) CollectionName() string {
	return "token_uri"
}
