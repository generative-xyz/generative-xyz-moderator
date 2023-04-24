package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type AISchoolJob struct {
	BaseEntity    `bson:",inline" json:"base_entity"`
	Params        string   `bson:"params" json:"params"`
	DatasetUUID   string   `bson:"dataset_uuid" json:"dataset_uuid"` //TODO: remove this field
	OutputUUID    string   `bson:"output_uuid" json:"output_uuid"`
	OutputLink    string   `bson:"output_link" json:"output_link"`
	JobID         string   `bson:"job_id" json:"job_id"`
	Status        string   `bson:"status" json:"status"`
	Progress      int      `bson:"progress" json:"progress"`
	ExecutedAt    int64    `bson:"executed_at" json:"executed_at"`
	CompletedAt   int64    `bson:"completed_at" json:"completed_at"`
	ClearedAt     int64    `bson:"cleared_at" json:"cleared_at"`
	CreatedBy     string   `bson:"created_by" json:"created_by"`
	Errors        string   `bson:"errors" json:"errors"`
	UsePFPDataset bool     `bson:"use_pfp_dataset" json:"use_pfp_dataset"` //TODO: remove this field
	Logs          string   `bson:"logs" json:"logs"`
	ErrLogs       string   `bson:"err_logs" json:"err_logs"`
	Datasets      []string `bson:"datasets" json:"datasets"`
}

func (job AISchoolJob) TableName() string {
	return utils.AI_SCHOOL_JOB
}

func (job AISchoolJob) ToBson() (*bson.D, error) {
	return helpers.ToDoc(job)
}

type AISchoolPresetDataset struct {
	BaseEntity  `bson:",inline" json:"base_entity"`
	Name        string `bson:"name" json:"name"`
	Thumbnail   string `bson:"thumbnail" json:"thumbnail"`
	DatasetURI  string `bson:"dataset_uri" json:"dataset_uri"`
	FileUUID    string `bson:"file_uuid" json:"file_uuid"`
	Creator     string `bson:"creator" json:"creator"`
	IsPrivate   bool   `bson:"is_private" json:"is_private"`
	Size        int    `bson:"size" json:"size"`
	NumOfAssets int    `bson:"num_of_assets" json:"num_of_assets"`
}

func (job AISchoolPresetDataset) TableName() string {
	return utils.AI_SCHOOL_DATASET
}

func (job AISchoolPresetDataset) ToBson() (*bson.D, error) {
	return helpers.ToDoc(job)
}
