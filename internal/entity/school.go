package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type AISchoolJob struct {
	BaseEntity  `bson:",inline" json:"base_entity"`
	Params      string `bson:"params" json:"params"`
	DatasetUUID string `bson:"dataset_uuid" json:"dataset_uuid"`
	OutputUUID  string `bson:"output_uuid" json:"output_uuid"`
	OutputLink  string `bson:"output_link" json:"output_link"`
	JobID       string `bson:"job_id" json:"job_id"`
	Status      string `bson:"status" json:"status"`
	Progress    int    `bson:"progress" json:"progress"`
	ExecutedAt  int64  `bson:"executed_at" json:"executed_at"`
	CompletedAt int64  `bson:"completed_at" json:"completed_at"`
	ClearedAt   int64  `bson:"cleared_at" json:"cleared_at"`
	CreatedBy   string `bson:"created_by" json:"created_by"`
}

func (job AISchoolJob) TableName() string {
	return utils.AI_SCHOOL_JOB
}

func (job AISchoolJob) ToBson() (*bson.D, error) {
	return helpers.ToDoc(job)
}
