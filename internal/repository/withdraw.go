package repository

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
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
	
	if filter.ProjectID != nil && *filter.ProjectID != "" {
		f = append(f, bson.M{"projectID": *filter.ProjectID})
	}
	
	if len(filter.ProjectIDs) > 0 {
		f = append(f, bson.M{"$in": bson.M{"projectID": filter.ProjectIDs } })
	}
	
	if filter.Status != nil {
		f = append(f, bson.M{"status": *filter.Status})
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
	
	spew.Dump(f)
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


