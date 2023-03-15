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
	unwindProject := bson.M{"$unwind": "$project"}
	unwindUser := bson.M{"$unwind": "$user"}
	addFields := bson.M{
		"$addFields": bson.M{
			"project_name": "$project.name",
			"user_name":    "$user.display_name",
		},
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
	if request.PageSize > 0 && request.PageSize <= limit {
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
		filters["status"] = *request.Status
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
	lookupDaoProjectVoted := bson.M{
		"$lookup": bson.M{
			"from":         "dao_project_voted",
			"localField":   "_id",
			"foreignField": "dao_project_id",
			"as":           "dao_project_voted",
		},
	}
	addFieldsCount := bson.M{
		"$addFields": bson.M{
			"voted": bson.M{
				"$filter": bson.M{
					"input": "$dao_project_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 1},
					},
				},
			},
			"against": bson.M{
				"$filter": bson.M{
					"input": "$dao_project_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 0},
					},
				},
			},
		},
	}
	projectAgg := bson.M{
		"$project": bson.M{
			"_id":               1,
			"uuid":              1,
			"created_at":        1,
			"seq_id":            1,
			"created_by":        1,
			"user":              1,
			"project_id":        1,
			"project":           1,
			"expired_at":        1,
			"status":            1,
			"dao_project_voted": 1,
			"user_name":         1,
			"project_name":      1,
			"total_vote":        bson.M{"$size": "$voted"},
			"total_against":     bson.M{"$size": "$against"},
		},
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
		addFields,
		matchSearch,
		lookupDaoProjectVoted,
		addFieldsCount,
		sorts,
		projectAgg)
	if err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}
