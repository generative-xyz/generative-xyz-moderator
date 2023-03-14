package entity

import (
	"time"

	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseEntity struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UUID           string             `json:"uuid" bson:"uuid"`
	BaseEntityNoID `bson:",inline"`
}

type BaseEntityNoID struct {
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

func (b *BaseEntity) SetID() {
	b.ID = primitive.NewObjectID()
	b.UUID = b.ID.Hex()

}

func (b *BaseEntity) GetID() string {
	return b.UUID
}

func (b *BaseEntityNoID) SetCreatedAt() {
	now := time.Now().UTC()
	b.CreatedAt = &now

}

func (b *BaseEntityNoID) SetUpdatedAt() {
	now := time.Now().UTC()
	b.UpdatedAt = &now

}

func (b *BaseEntityNoID) SetDeletedAt() {
	now := time.Now().UTC()
	b.DeletedAt = &now
}

func (b *BaseEntityNoID) Decode(from *primitive.D) error {
	err := helpers.Transform(from, b)
	if err != nil {
		return err
	}
	return nil
}

type FilterString struct {
	Keyword           string
	ListCollectionIDs string
	ListPrices        string
	ListIDs           string
}
