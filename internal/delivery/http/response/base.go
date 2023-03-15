package response

import (
	"time"
)

type PaginationResponse struct {
	Result    interface{} `json:"result"`
	Page      int64       `json:"page"`
	PageSize  int64       `json:"pageSize"`
	TotalPage int64       `json:"totalPage"`
	// Next int64
	// Prev int64
	//Limit int64
	Total   int64  `json:"total"`
	Currsor string `json:"cursor"`
}

type BaseEntity struct {
	ID        string     `json:"id" bson:"_id"`
	UUID      string     `json:"uuid" bson:"uuid"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}
