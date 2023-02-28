package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEntityNoID interface {
	TableName() string
	ToBson() (*bson.D, error)
	SetCreatedAt()
	SetUpdatedAt()
	SetDeletedAt()
	Decode(from *primitive.D) error
}

type IEntity interface {
	SetID()
	GetID() string
	IEntityNoID
}

type SortType int

const (
	SORT_ASC  SortType = 1
	SORT_DESC SortType = -1
)

type BaseFilters struct {
	Page   int64
	Limit  int64
	SortBy string
	Sort   SortType
}

type Pagination struct {
	Result   interface{}
	Page     int64
	PageSize int64
	Total    int64
	Currsor  string
}
