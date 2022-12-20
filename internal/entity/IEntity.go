package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEntity interface {
	TableName() string
	ToBson() (*bson.D, error)
	SetID()
	GetID() string
	SetCreatedAt()
	SetUpdatedAt()
	SetDeletedAt()
	Decode(from *primitive.D) error
}

type SortType int 
const (
	SORT_ASC SortType = 1
	SORT_DESC SortType = -1	
)

type BaseFilters struct {
	Page int64
	Limit int64
	SortBy string
	Sort SortType
}

type Pagination struct {
	Data interface{}
	Page int64
	TotalPage int64
	Next int64
	Prev int64
	Limit int64
	Total int64
	Currsor string
}