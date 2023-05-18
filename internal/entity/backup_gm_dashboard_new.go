package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type CachedGMDashBoardNew struct {
	BaseEntity   `bson:",inline"`
	Contributors int     `bson:"contributors" json:"contributors"`
	UsdtValue    float64 `bson:"usdtValue" json:"usdtValue"`

	//only save uploaded link
	BackupURL      string `bson:"backup_url" json:"backup_url"`
	BackupFilePath string `bson:"backup_file_path" json:"backup_file_path"`
	BackupFileName string `bson:"backup_file_name" json:"backup_file_name"`

	//old data
	OldBackupURL      string `bson:"old_backup_url" json:"old_backup_url"`
	OldBackupFilePath string `bson:"old_backup_file_path" json:"old_backup_file_path"`
	OldBackupFileName string `bson:"old_backup_file_name" json:"old_backup_file_name"`
}

func (u CachedGMDashBoardNew) TableName() string {
	return utils.COLLECTION_CACHED_GM_DASHBOARD_NEW
}

func (u CachedGMDashBoardNew) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
