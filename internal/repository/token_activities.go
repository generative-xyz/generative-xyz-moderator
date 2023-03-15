package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertTokenActivity(data *entity.TokenActivity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) buildFilterTokenActivities(filter entity.FilterTokenActivities) bson.M {
	f := bson.M{}
	if filter.InscriptionID != nil {
		f["inscription_id"] = *filter.InscriptionID
	}
	if filter.ProjectID != nil {
		f["project_id"] = *filter.ProjectID
	}
	return f
}

func (r Repository) GetTokenActivities(filter entity.FilterTokenActivities) (*entity.Pagination, error) {
	activities := []entity.TokenActivity{}
	resp := &entity.Pagination{}
	f := r.buildFilterTokenActivities(filter)
	if filter.SortBy == "" {
		filter.SortBy = "time"
		filter.Sort = entity.SORT_DESC
	}
	s := []Sort{
		{SortBy: filter.SortBy, Sort: filter.Sort},
		{SortBy: "uuid", Sort: entity.SORT_DESC},
	}
	t, err := r.Paginate(entity.TokenUri{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, s, &activities)
	if err != nil {
		return nil, err
	}

	resp.Result = activities
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	//resp.PageSize = filter.Limit
	return resp, nil
}
