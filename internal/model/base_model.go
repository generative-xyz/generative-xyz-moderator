package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel ...
type BaseModel struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (b *BaseModel) WithTimeInfo() {
	if b == nil {
		return
	}
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now
}
