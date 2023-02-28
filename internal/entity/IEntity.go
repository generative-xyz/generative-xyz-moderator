package entity

import (
	"net/http"
	"strconv"

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
	Result   interface{} `json:"result"`
	Page     int64       `json:"page"`
	PageSize int64       `json:"pageSize"`
	Total    int64       `json:"total"`
	Currsor  string      `json:"currsor"`
}

func NewDefaultPagination(opts ...int64) *Pagination {
	page := int64(1)
	limit := int64(10)
	if len(opts) > 0 && opts[0] > 0 {
		page = opts[0]
	}
	if len(opts) > 1 && opts[1] > 0 {
		limit = opts[1]
	}
	return &Pagination{
		PageSize: limit,
		Page:     page,
	}
}

func GetPagination(r *http.Request) *Pagination {
	pag := NewDefaultPagination()
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	// TODO add sort
	pag.Page = int64(page)
	pag.PageSize = int64(limit)
	return pag
}
