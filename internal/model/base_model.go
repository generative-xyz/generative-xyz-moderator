package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type IEntity interface {
	CollectionName() string
	SetID()
	GetID() string
	SetCreatedAt()
	SetUpdatedAt()
	SetDeletedAt()
	WithTimeInfo()()
}

// BaseModel ...
type BaseModel struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	UUID            string `bson:"uuid"`
	CreatedAt *time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

func (b *BaseModel) WithTimeInfo() {
	if b == nil {
		return
	}
	now := time.Now().UTC()
	b.CreatedAt = &now
	b.UpdatedAt = &now
}


func (b *BaseModel) SetID() {
	b.ID = primitive.NewObjectID()
	b.UUID = b.ID.Hex()

}

func (b *BaseModel) GetID() string {
	return b.UUID
}

func (b *BaseModel) SetCreatedAt() {
	now := time.Now().UTC()
	b.CreatedAt = &now

}

func (b *BaseModel) SetUpdatedAt() {
	now := time.Now().UTC()
	b.UpdatedAt = &now

}

func (b *BaseModel) SetDeletedAt() {
	now := time.Now().UTC()
	b.DeletedAt = &now
}
