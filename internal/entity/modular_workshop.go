package entity

import (
	"encoding/json"
	"rederinghub.io/utils"
	"strings"
)

type ModularWorkshopEntity struct {
	BaseNewEntity `bson:",inline"`
	Name          string `bson:"name" json:"name"`
	OwnerAddr     string `bson:"owner_addr" json:"owner_addr" `
	MetaData      string `bson:"meta_data" json:"meta_data"`
	Thumbnail     string `bson:"thumbnail" json:"thumbnail"`
	Public        bool   `bson:"public" json:"public"`   // save to DB
	IsGuestMode   *bool  `bson:"-" json:"is_guest_mode"` //from FE - Public, private keys are denied by FE code.
}

func (u ModularWorkshopEntity) TableName() string {
	return utils.COLLECTION_MODULAR_WORKSHOP
}

// defined by FE
type MetaDataBasicInfo struct {
	Type          string `json:"type"`
	GroupId       string `json:"groupId"`
	InscriptionId string `json:"inscriptionId"`
}

func (u ModularWorkshopEntity) GetListInscriptionIds() []string {
	var info []MetaDataBasicInfo
	json.Unmarshal([]byte(u.MetaData), &info)
	var inscriptionIds []string
	for _, data := range info {
		inscriptionIds = append(inscriptionIds, strings.TrimSpace(strings.ToLower(data.InscriptionId)))
	}
	return inscriptionIds
}

type ModularWorkshopStatistic struct {
	TotalOwner int   `json:"total_owner"`
	TotalModel int64 `json:"total_model"`
}
