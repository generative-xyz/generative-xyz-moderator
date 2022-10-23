package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel ...
type BaseModel struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
