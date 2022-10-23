package model

type Template struct {
	BaseModel
	TokenID string `bson:"token_id" json:"token_id"`
}

func (m Template) CollectionName() string {
	return "templates"
}
