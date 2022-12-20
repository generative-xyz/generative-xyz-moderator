package paginator

import (
	"math"

	"gorm.io/gorm"
)

func NewPaginator() *GormPaginator {
	p := new(GormPaginator)
	return p
}

type GormPaginator struct {
	Items     interface{} `json:"items"`
	Limit     int         `json:"limit"`
	Page      int         `json:"page"`
	TotalItem int         `json:"total_item"`
	TotalPage int         `json:"total_page"`
	Offset    int         `json:"offset"`
	PrevPage  int         `json:"prev_page"`
	NextPage  int         `json:"next_page"`
}

func (pagination *GormPaginator) Paging(db *gorm.DB, result interface{}) {
	var countChannel = make(chan int64)
	var offset int

	go countRecords(db, result, countChannel)

	count := <-countChannel

	// calculate offset from limit & page
	if pagination.Page == 1 {
		offset = 0
	} else {
		offset = (pagination.Page - 1) * pagination.Limit
	}
	pagination.Offset = offset
	// set limit and offset for db
	db = db.Limit(pagination.Limit).Offset(pagination.Offset)

	pagination.TotalItem = int(count)
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(pagination.Limit)))
	if pagination.Page > 1 {
		pagination.PrevPage = pagination.Page - 1
	} else {
		pagination.PrevPage = pagination.Page
	}

	if pagination.Page == pagination.TotalPage {
		pagination.NextPage = pagination.Page
	} else {
		pagination.NextPage = pagination.Page + 1
	}

}

func countRecords(db *gorm.DB, anyType interface{}, countChannel chan int64) {
	var count int64
	db.Model(anyType).Count(&count)
	countChannel <- count
}
