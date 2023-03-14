package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
)

func (s Repository) ListDAOProject(ctx context.Context, request *request.ListDaoProjectRequest) ([]*entity.DaoProject, int64, error) {
	limit := int64(100)
	filters := make(bson.M)
	filterIdOperation := "$lt"
	sorts := bson.M{
		"$sort": bson.D{{Key: "_id", Value: -1}},
	}
	matchFilters := bson.M{"$match": filters}
	lookupProject := bson.M{
		"$lookup": bson.M{
			"from":         "projects",
			"localField":   "project_id",
			"foreignField": "_id",
			"as":           "project",
		},
	}
	lookupUser := bson.M{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "wallet_address",
			"as":           "user",
		},
	}
	lookupDaoProjectVoted := bson.M{
		"$lookup": bson.M{
			"from":         "dao_project_voted",
			"localField":   "_id",
			"foreignField": "dao_project_id",
			"as":           "dao_project_voted",
		},
	}
	unwindProject := bson.M{"$unwind": "$project"}
	unwindUser := bson.M{"$unwind": "$user"}
	addProjectName := bson.M{
		"$addFields": bson.M{"project_name": "$project.name"},
	}
	addUserName := bson.M{
		"$addFields": bson.M{"user_name": "$user.display_name"},
	}
	if len(request.Sorts) > 0 {
		sort := bson.D{}
		for _, srt := range request.Sorts {
			sort = append(sort, bson.E{Key: srt.Field, Value: srt.Type})
			if srt.Field == "_id" && srt.Type == entity.SORT_ASC {
				filterIdOperation = "$gt"
			}
		}
		sorts = bson.M{
			"$sort": sort,
		}
	}
	if request.PageSize <= limit {
		limit = request.PageSize
	}
	if request.Id != nil {
		id, err := primitive.ObjectIDFromHex(*request.Id)
		if err != nil {
			return nil, 0, err
		}
		filters["_id"] = id
	}
	if request.Status != nil {
		filters["status"] = request.Status
	}
	if request.Cursor != "" {
		if id, err := primitive.ObjectIDFromHex(request.Cursor); err == nil {
			filters["_id"] = bson.M{filterIdOperation: id}
		}
	}
	filterSearch := make(bson.M)
	matchSearch := bson.M{"$match": filterSearch}
	if request.Keyword != nil {
		filterSearch["$or"] = bson.A{
			bson.M{"project_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
			bson.M{"user_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
		}
	}
	projects := []*entity.DaoProject{}
	total, err := s.Aggregation(ctx,
		entity.DaoProject{}.TableName(),
		0,
		limit,
		&projects,
		matchFilters,
		lookupProject,
		unwindProject,
		lookupUser,
		unwindUser,
		addProjectName,
		addUserName,
		matchSearch,
		lookupDaoProjectVoted,
		sorts)
	if err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}
