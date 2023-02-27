package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)


type  AggregateWalletAddressItem struct {
	ID AggregateItemID `bson:"_id" json:"id"`
	Amount float64 `bson:"amount" json:"amount"`
}

type  AggregateWalleRespItem struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
	Amount string `bson:"amount" json:"amount"`
}

type  AggregateAmount struct {
	Amount float64 `bson:"amount" json:"amount"`
}

type  AggregateItemID struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
}

//analytis
// type  AggregateWalletAddres struct {
// 	Items []AggregateWalletAddressItem  `json:"items"`
// 	TotalBTC float64 `json:"totalAmountBTC"`
// 	TotalETH float64 `json:"totalAmountETH"`
// }

type FilterVolume struct {
	projectIDs []string
	amountType *string
	UserID *string
	ProjectID *string
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
	AmountType *string  `bson:"amountType"`
	CreatorAddress *string `bson:"creatorAddress"`
	ProjectID *string  `bson:"projectID"`
	Amount *string  `bson:"amount"`
	Project  VolumeProjectInfo `bson:"project"`
	User  VolumnUserInfo `bson:"user"`
}

func (u UserVolumn) TableName() string {
	return utils.COLLECTION_USER_VOLUMN
}

func (u UserVolumn) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
