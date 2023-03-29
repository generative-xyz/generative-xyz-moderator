package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ProjectZipStatus int

const (
	UzipStatusFail    ProjectZipStatus = 0
	UzipStatusSuccess ProjectZipStatus = 1
)

type ProjectZipLinks struct {
	BaseEntity             `bson:",inline" json:"-"`
	ProjectID              string              `bson:"projectID" json:"projectID"`
	ZipLink                string              `bson:"zipLink" json:"zipLink"`
	Status                 ProjectZipStatus    `bson:"status" json:"status"`
	Message                string              `bson:"message" json:"message"`
	Logs                   []ProjectZipLinkLog `bson:"logs" json:"logs"`
	ReTries                int                 `bson:"retries" json:"retries"`
}

type ProjectZipLinkLog struct {
	Message     string           `bson:"message" json:"message"`
	Status      ProjectZipStatus `bson:"status" json:"status"`
	CreatedTime *time.Time       `bson:"time" json:"time"`
}

func (u ProjectZipLinks) TableName() string {
	return utils.COLLECTION_PROJECT_ZIPLINKS
}

func (u ProjectZipLinks) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
