package entity

import (
	"rederinghub.io/utils"
	"time"
)

type FileFragmentStatus int

const (
	FileFragmentStatusCreated FileFragmentStatus = iota + 1
	FileFragmentStatusProcessing
	FileFragmentStatusDone
	FileFragmentStatusError
)

type TokenFileFragment struct {
	BaseEntity `bson:",inline" json:"base_entity"`
	TokenId    string             `json:"token_id" bson:"token_id"`
	FilePath   string             `json:"file_path" bson:"file_path"`
	Sequence   int                `json:"sequence" bson:"sequence"`
	Data       []byte             `json:"data" bson:"data"`
	Status     FileFragmentStatus `json:"status" bson:"status"`
	Note       string             `json:"note" bson:"note"`
	UploadTime *time.Time         `json:"upload_time" bson:"upload_time"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

func (m TokenFileFragment) TableName() string {
	return utils.TOKEN_FILE_FRAGMENT
}

type TokenFragmentJobStatus int

const (
	FragmentJobStatusPending TokenFragmentJobStatus = iota + 1
	FragmentJobStatusProcessing
	FragmentJobStatusDone
	FragmentJobStatusError
)

type TokenFragmentJob struct {
	BaseEntity `bson:",inline" json:"base_entity"`
	TokenId    string                 `json:"token_id" bson:"token_id"`
	FilePath   string                 `json:"file_path" bson:"file_path"`
	Status     TokenFragmentJobStatus `json:"status" bson:"status"`
	Note       string                 `json:"note" bson:"note"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
}

func (m TokenFragmentJob) TableName() string {
	return utils.TOKEN_FILE_FRAGMENT_JOB
}
