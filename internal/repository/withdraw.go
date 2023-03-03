package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)


func (r Repository) CreateWithDraw(input *entity.Withdraw) error {
	err := r.InsertOne(input.TableName(), input)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FilterWithDraw(filter *entity.FilterWithdraw) (*entity.Pagination, error) {
	confs := []entity.Withdraw{}
	resp := &entity.Pagination{}
	f := bson.M{}

	var s []Sort
	s = append(s, Sort{SortBy: "created_at", Sort: entity.SORT_DESC})
	
	p, err := r.Paginate(utils.COLLECTION_WITHDRAW, filter.Page, filter.Limit, f, nil, s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil

}


func (r Repository) AggregateWithDrawByUser(filter *entity.FilterWithdraw) ( []entity.AggregateAmount, error) {
	confs := []entity.AggregateAmount{}
	f :=  bson.A{}

	if filter.PaymentType != nil && *filter.PaymentType != "" {
		f = append(f, bson.M{"payType": *filter.PaymentType})
	}
	
	if filter.WalletAddress != nil && *filter.WalletAddress != "" {
		f = append(f, bson.M{"walletAddress": *filter.WalletAddress})
	}
	
	if filter.WithdrawItemID != nil && *filter.WithdrawItemID != "" {
		f = append(f, bson.M{"withdrawItemID": *filter.WithdrawItemID})
	}
	
	if len(filter.WithdrawItemIDs) > 0 {
		f = append(f, bson.M{"$in": bson.M{"withdrawItemID": filter.WithdrawItemIDs } })
	}
	
	if filter.Status != nil {
		f = append(f, bson.M{"status": *filter.Status})
	}
	
	if  len(filter.Statuses) > 0 {
		f = append(f, bson.M{"status": bson.M{"$in": filter.Statuses}})
	}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": f}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": 
			bson.M{"walletAddress": "$walletAddress", "payType": "$payType"}, 
			"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
		}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}
	
	cursor, err := r.DB.Collection(entity.Withdraw{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {	
		tmp := &entity.AggregateAmount{}
		err = helpers.Transform(item, tmp)
		confs = append(confs, *tmp)
	}
	
	return confs, nil

}

func (r Repository) UpdateWithDrawStatus(UUID string, status int) error {
	filter := bson.D{
		{Key: utils.KEY_UUID, Value: UUID},
	}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_WITHDRAW).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) GetLastWithdraw(filter entity.FilterWithdraw) (*entity.Withdraw, error) {
	wd := []entity.Withdraw{}
	f := bson.M{}
	if filter.WalletAddress != nil && *filter.WalletAddress != "" {
		f = bson.M{"walletAddress": *filter.WalletAddress}
	}
	
	if filter.WithdrawItemID != nil && *filter.WithdrawItemID != "" {
		f =  bson.M{"withdrawItemID": *filter.WithdrawItemID}
	}
	
	if filter.PaymentType != nil && *filter.PaymentType != "" {
		f =  bson.M{"payType": *filter.PaymentType}
	}
	
	if len(filter.Statuses) > 0 {
		f =  bson.M{"status": bson.M{"$in": filter.Statuses}}
	}

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.DB.Collection(entity.Withdraw{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &wd); err != nil {
		return nil, err
	}

	if len(wd) <= 0 {
		return  nil, errors.New("document not found")
	}

	return &wd[0], nil
}


