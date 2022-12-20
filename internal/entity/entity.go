package entity

import (
	"time"

	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type BaseEntity struct {
	DeletedAt *time.Time `bson:"deleted_at"`
	CreatedAt *time.Time `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	ID            primitive.ObjectID `bson:"_id"`
	UUID            string `bson:"uuid"`
}


func (b *BaseEntity) SetID() {
	b.ID = primitive.NewObjectID()
	b.UUID = b.ID.Hex()

}

func (b *BaseEntity) GetID() string {
	return b.UUID
}

func (b *BaseEntity) SetCreatedAt() {
	now := time.Now().UTC()
	b.CreatedAt = &now

}

func (b *BaseEntity) SetUpdatedAt() {
	now := time.Now().UTC()
	b.UpdatedAt = &now

}

func (b *BaseEntity) SetDeletedAt() {
	now := time.Now().UTC()
	b.DeletedAt = &now
}

func (b *BaseEntity) Decode(from *primitive.D) error {
	err := helpers.Transform(from, b)
	if err != nil {
		return err
	}
	return nil
}
