package entity

import "rederinghub.io/utils"

type ModularWorkshopEntity struct {
	BaseEntity `bson:",inline"`
	Name       string `bson:"name" json:"name"`
	OwnerAddr  string `bson:"owner_addr" json:"owner_addr" `
	MetaData   string `bson:"meta_data" json:"meta_data"`
}

func (u ModularWorkshopEntity) TableName() string {
	return utils.COLLECTION_MODULAR_WORKSHOP
}
