package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertReferral(data *entity.Referral) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FilterReferrals(filter entity.FilterReferrals) bson.M {
	f := bson.M{}
	if filter.ReferreeID != nil && *filter.ReferreeID != "" {
		f["referree_id"] = primitive.Regex{Pattern:  *filter.ReferreeID, Options: "i"}
	}
	if filter.ReferrerID != nil && *filter.ReferrerID != "" {
		f["referrer_id"] = primitive.Regex{Pattern:  *filter.ReferrerID, Options: "i"}
	}
	
	if filter.ReferrerAddress != nil && *filter.ReferrerAddress != "" {
		f["referrer.wallet_address"] = primitive.Regex{Pattern:  *filter.ReferrerAddress, Options: "i"}
	}
	
	if filter.ReferreeAddress != nil && *filter.ReferreeAddress != "" {
		f["referree.wallet_address"] = primitive.Regex{Pattern:  *filter.ReferreeAddress, Options: "i"}
	}
	
	// if filter.PayType != nil {
	// 	f["referrer_id"] = primitive.Regex{Pattern:  *filter.ReferrerID, Options: "i"}
	// }
	return f
}

func (r Repository) GetReferrals(filter entity.FilterReferrals) (*entity.Pagination, error) {
	confs := []entity.Referral{}
	resp := &entity.Pagination{}
	f := r.FilterReferrals(filter)
	p, err := r.Paginate(utils.COLLECTION_REFERRALS, filter.Page, filter.Limit, f, bson.M{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetAllReferrals(filter entity.FilterReferrals) ([]entity.Referral, error) {
	refs := []entity.Referral{}
	f := r.FilterReferrals(filter)
	cursor, err := r.DB.Collection(entity.Referral{}.TableName()).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &refs); err != nil {
		return nil, err
	}
	return refs, err
}

func (r Repository) CountReferralOfReferee(referreeID string) (int64, error) {
	f := bson.M{
		"referree_id": referreeID,
	}

	count, err := r.DB.Collection(entity.Referral{}.TableName()).CountDocuments(context.TODO(), f)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r Repository) UpdateReferral(ID string, data *entity.Referral) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.Referral{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) GetReferral(filter entity.FilterReferrals) ([]entity.Referral, error) {
	ref := []entity.Referral{}
	f := r.FilterReferrals(filter)
	cursor, err := r.DB.Collection(entity.Referral{}.TableName()).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (r Repository) GetAReferral(filter entity.FilterReferrals) (*entity.Referral, error) {
	ref := &entity.Referral{}
	f := r.FilterReferrals(filter)
	cursor, err := r.DB.Collection(entity.Referral{}.TableName()).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &ref); err != nil {
		return nil, err
	}

	return ref, nil
}