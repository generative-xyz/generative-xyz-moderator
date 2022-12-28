package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Files struct {
	BaseEntity`bson:",inline"`
	FileName string `bson:"file_name"`
	UploadedBy string `bson:"uploaded_by"`
	URL string `bson:"url"`
	MineType     string       `bson:"mime_type"`
	FileSize     int       `bson:"file_size"`
}

type FilterFiles struct {
	BaseFilters
	Name *string
	UploadedBy *string
}

func (u Files) TableName() string { 
	return utils.COLLECTION_FILES
}

func (u Files) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}