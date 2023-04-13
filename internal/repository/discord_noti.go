package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateDiscordNoti(noti entity.DiscordNoti) error {
	err := r.InsertOne(noti.TableName(), &noti)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetDiscordNotifications(req entity.GetDiscordNotiReq) (*entity.Pagination, error) {
	confs := make([]entity.DiscordNoti, 0)
	resp := &entity.Pagination{}

	s := []Sort{
		{SortBy: "created_at", Sort: entity.SORT_ASC},
	}

	f := bson.M{}

	if req.Status != nil {
		f["status"] = req.Status
	}

	p, err := r.Paginate(entity.DiscordNoti{}.TableName(), req.Page, req.Limit, f, bson.D{}, s, &confs)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = req.Limit

	return resp, nil
}

func (r Repository) UpdateDiscordNotiAddRetry(id string) error {
	filter := bson.M{
		"uuid": id,
	}
	update := bson.M{
		"$inc": bson.M{"num_retried": 1},
	}
	_, err := r.DB.Collection(entity.DiscordNoti{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r Repository) UpdateDiscordStatus(id string, status entity.DiscordNotiStatus, note string) error {
	filter := bson.M{
		"uuid": id,
	}
	update := bson.M{
		"$set": bson.M{
			"status": status,
			"note":   note,
		},
	}
	_, err := r.DB.Collection(entity.DiscordNoti{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
