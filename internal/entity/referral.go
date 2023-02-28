package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Referral struct {
	BaseEntity`bson:",inline"`
	ReferrerID string `bson:"referrer_id"`
	ReferreeID string `bson:"referree_id"`
	Referrer   *Users `bson:"referrer"`
	Referree   *Users `bson:"referree"`
	Percent    int32 	`bson:"percent"`
	//ReferreeVolumn ReferreeVolumn `bson:"referreeVolumn"`
}

type ReferreeVolumn struct {
	Amount string  `bson:"amount"`
	AmountType string  `bson:"amountType"`
	ProjectID string  `bson:"projectID"`
	Earn string  `bson:"earn"`
	GenEarn string  `bson:"genEarn"`
}

type FilterReferrals struct {
	BaseFilters
	ReferrerID *string
	ReferreeID *string
	ReferrerAddress *string
	ReferreeAddress *string
}

func (u Referral) TableName() string { 
	return utils.COLLECTION_REFERRALS
}

func (u Referral) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
