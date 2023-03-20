package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) GetPendingBTCSubmitTx() ([]entity.BTCTransactionSubmit, error) {
	listings := []entity.BTCTransactionSubmit{}
	f := bson.M{
		"status": bson.M{"$in": []entity.BTCTransactionSubmitStatus{entity.StatusBTCTransactionSubmit_Waiting, entity.StatusBTCTransactionSubmit_Pending}},
	}
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_TX_SUBMIT).Find(context.TODO(), f, &options.FindOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r Repository) GetPendingBTCSubmitTxByInscriptionID(inscriptionID string) ([]entity.BTCTransactionSubmit, error) {
	listings := []entity.BTCTransactionSubmit{}
	f := bson.M{
		"related_inscriptions": bson.M{"$in": []string{inscriptionID}},
		"status":               bson.M{"$in": []entity.BTCTransactionSubmitStatus{entity.StatusBTCTransactionSubmit_Waiting, entity.StatusBTCTransactionSubmit_Pending}},
	}
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_TX_SUBMIT).Find(context.TODO(), f, &options.FindOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r Repository) UpdateBTCTxSubmitStatus(model *entity.BTCTransactionSubmit) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"status": model.Status,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}
