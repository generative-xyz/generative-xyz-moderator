package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type  AggregateWalletAddressItem struct {
	ID AggregateItemID `bson:"_id" json:"id"`
	Amount float32 `bson:"amount" json:"amount"`
}

type  AggregateWalleRespItem struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
	Amount string `bson:"amount" json:"amount"`
}


type  AggregateProjectItemID struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
	Amount float32 `bson:"amount" json:"amount"`
	BtcRate float32 `bson:"btcRate" json:"btcRate"`
	EthRate float32 `bson:"ethRate" json:"ethRate"`
}

type  AggregateProjectItem struct {
	ID AggregateProjectItemID `bson:"_id" json:"id"`
	Amount float64 `bson:"amount" json:"amount"`
	Minted int `bson:"minted" json:"minted"`
}

type  AggregateProjectItemResp struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
	Amount float64 `bson:"amount" json:"amount"`
	Minted int `bson:"minted" json:"minted"`
	BtcRate float32 `bson:"btcRate" json:"btcRate"`
	EthRate float32 `bson:"ethRate" json:"ethRate"`
}

type  AggregateAmount struct {
	ID AggregateItemID `bson:"_id" json:"id"`
	Amount float64 `bson:"amount" json:"amount"`
}

type  AggregateItemID struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
	CreatorAddress string `bson:"creatorAddress" json:"creatorAddress"`
}

//analytis
// type  AggregateWalletAddres struct {
// 	Items []AggregateWalletAddressItem  `json:"items"`
// 	TotalBTC float64 `json:"totalAmountBTC"`
// 	TotalETH float64 `json:"totalAmountETH"`
// }

type FilterVolume struct {
	ProjectIDs []string
	AmountType *string
	UserID *string
	ProjectID *string
	CreatorAddress *string
}

type VolumeProjectInfo  struct{
	Name string `bson:"name"`
	TokenID string `bson:"tokenID"`
	Thumnail string `bson:"thumbnail"`
}

type VolumnUserInfo struct {
	WalletAddress *string `bson:"walletAddress"`
	WalletAddressBTC *string `bson:"walletAddressBTC"`
	DisplayName *string  `bson:"displayName"`
	Avatar   *string  `bson:"avatar"`
}

type UserVolumn struct {
	BaseEntity              `bson:",inline" json:"-"`
	PayType *string  `bson:"payType"`
	CreatorAddress *string `bson:"creatorAddress"`
	ProjectID *string  `bson:"projectID"`
	Amount *string  `bson:"amount"`
	Minted int  `bson:"minted"`
	Project  VolumeProjectInfo `bson:"project"`
	User  VolumnUserInfo `bson:"user"`
}

func (u UserVolumn) TableName() string {
	return utils.COLLECTION_USER_VOLUMN
}

func (u UserVolumn) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
