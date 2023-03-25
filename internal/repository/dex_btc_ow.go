package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateDexBTCOWSubmit(listing *entity.DexBTCOWSubmitTx) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) CreateDexBTCOWInscriptions(listing []entity.DexBTCOWInscription) error {
	inscriptionListing := make([]entity.IEntity, 0, len(listing))
	for _, tokenHolder := range listing {
		_tokenHolder := tokenHolder
		inscriptionListing = append(inscriptionListing, &_tokenHolder)
	}
	err := r.InsertMany(listing[0].TableName(), inscriptionListing)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ClearDexBTCOWCollectionListing(collectionSlug string) error {
	_, err := r.DeleteMany(context.Background(), entity.DexBTCOWInscription{}.TableName(), bson.M{"collection_slug": collectionSlug})
	if err != nil {
		return err
	}
	return nil
}
