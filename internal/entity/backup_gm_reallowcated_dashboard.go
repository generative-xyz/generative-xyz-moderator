package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type CachedGMReAllocatedDashBoard struct {
	BaseEntity     `bson:",inline"`
	Contributors   int     `bson:"contributors" json:"contributors"`
	UsdtValue      float64 `bson:"usdtValue" json:"usdtValue"`
	UsdtExtra      float64 `bson:"usdtExtra" json:"usdtExtra"`
	TotalGMReceive float64 `bson:"totalGMReceive" json:"totalGMReceive"`

	//only save uploaded link
	BackupURL      string `bson:"backup_url" json:"backup_url"`
	BackupFilePath string `bson:"backup_file_path" json:"backup_file_path"`
	BackupFileName string `bson:"backup_file_name" json:"backup_file_name"`
}

func (u CachedGMReAllocatedDashBoard) TableName() string {
	return utils.COLLECTION_CACHED_REALLOCATED_GM_DASHBOARD
}

func (u CachedGMReAllocatedDashBoard) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
