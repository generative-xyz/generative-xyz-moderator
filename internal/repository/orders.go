package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) InsertOrder(obj *entity.OrdersAddress) error {
	obj.SetCreatedAt()

	bData, err := obj.ToBson()
	if err != nil {
		return err
	}

	inserted, err := r.DB.Collection(obj.TableName()).InsertOne(context.TODO(), &bData)
	if err != nil {
		return err
	}

	objIDObject := inserted.InsertedID.(primitive.ObjectID)
	objIDStr := objIDObject.Hex()

	r.CreateCache(obj.TableName(), objIDStr, obj)
	err = obj.Decode(bData)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindOrderBy(orderID, payType string) (*entity.OrdersAddress, error) {
	resp := &entity.OrdersAddress{}
	usr, err := r.FilterOne(entity.OrdersAddress{}.TableName(), bson.D{
		{"order_id", orderID},
		{"address_type", payType},
	})

	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindOrderByIDs(orderIDs []string) ([]*entity.OrdersAddress, error) {
	resp := []*entity.OrdersAddress{}
	f := bson.D{{"order_id", bson.D{{"$in", orderIDs}}}}

	cursor, err := r.DB.Collection(utils.ORDERS).Find(context.TODO(), f, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindOrderByStatus(statuses []entity.OrderStatus) ([]*entity.OrdersAddress, error) {
	resp := []*entity.OrdersAddress{}
	f := bson.D{{"status", bson.D{{"$in", statuses}}}}

	cursor, err := r.DB.Collection(utils.ORDERS).Find(context.TODO(), f, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
