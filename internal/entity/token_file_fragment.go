package entity

import (
	"rederinghub.io/utils"
	"time"
)

type FileFragmentStatus int

const (
	FileFragmentStatusPending FileFragmentStatus = iota + 1
	FileFragmentStatusDone
	FileFragmentStatusError
)

type TokenFileFragment struct {
	BaseEntity  `bson:",inline" json:"base_entity"`
	TokenId     string             `json:"token_id" bson:"token_id"`
	FilePath    string             `json:"file_path" bson:"file_path"`
	Sequence    int                `json:"sequence" bson:"sequence"`
	Data        []byte             `json:"data" bson:"data"`
	Status      FileFragmentStatus `json:"status" bson:"status"`
	Note        string             `json:"note" bson:"note"`
	UploadedAt  *time.Time         `json:"uploaded_at" bson:"uploaded_at"`
	TxSendNft   string             `bson:"tx_send_nft" json:"tx_send_nft"`
	TxStoreNft  string             `bson:"tx_store_nft" json:"tx_store_nft"`
	GasPrice    string             `bson:"gas_price" json:"gas_price"`
	NewGasPrice string             `bson:"new_gas_price" json:"new_gas_price"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
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
