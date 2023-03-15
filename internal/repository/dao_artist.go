package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_artist_voted"
)

func (s Repository) ListDAOArtist(ctx context.Context, request *request.ListDaoArtistRequest) ([]*entity.DaoArtist, int64, error) {
	limit := int64(100)
	filters := make(bson.M)
	filterIdOperation := "$lt"
	sorts := bson.M{
		"$sort": bson.D{{Key: "_id", Value: -1}},
	}
	matchFilters := bson.M{"$match": filters}
	lookupUser := bson.M{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "wallet_address",
			"as":           "user",
		},
	}
	unwindUser := bson.M{"$unwind": "$user"}
	addFieldUserName := bson.M{
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
			bson.M{"user_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
		}
	}
	lookupDaoArtistVoted := bson.M{
		"$lookup": bson.M{
			"from":         "dao_artist_voted",
			"localField":   "_id",
			"foreignField": "dao_artist_id",
			"as":           "dao_artist_voted",
		},
	}
	addFieldsCount := bson.M{
		"$addFields": bson.M{
			"verify": bson.M{
				"$filter": bson.M{
					"input": "$dao_artist_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 1},
					},
				},
			},
			"report": bson.M{
				"$filter": bson.M{
					"input": "$dao_artist_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 0},
					},
				},
			},
		},
	}
	projectAgg := bson.M{
		"$project": bson.M{
			"_id":              1,
			"uuid":             1,
			"created_at":       1,
			"seq_id":           1,
			"created_by":       1,
			"user":             1,
			"expired_at":       1,
			"status":           1,
			"dao_artist_voted": 1,
			"user_name":        1,
			"total_verify":     bson.M{"$size": "$verify"},
			"total_report":     bson.M{"$size": "$report"},
		},
	}
	projects := []*entity.DaoArtist{}
	total, err := s.Aggregation(ctx,
		entity.DaoArtist{}.TableName(),
		0,
		limit,
		&projects,
		matchFilters,
		lookupUser,
		unwindUser,
		addFieldUserName,
		matchSearch,
		lookupDaoArtistVoted,
		addFieldsCount,
		sorts,
		projectAgg)
	if err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}

func (s Repository) CountDAOArtistVoteByStatus(ctx context.Context, daoArtistId primitive.ObjectID, status dao_artist_voted.Status) int {
	match := bson.M{"$match": bson.M{
		"dao_artist_id": daoArtistId,
		"status":        status,
	}}
	group := bson.M{
		"$group": bson.M{
			"_id":   "$dao_artist_id",
			"count": bson.M{"$sum": 1},
		},
	}
	cur, err := s.DB.Collection(entity.DaoArtist{}.TableName()).Aggregate(ctx, bson.A{match, group})
	if err != nil {
		return 0
	}
	var results []*Count
	if err := cur.All(ctx, &results); err != nil {
		return 0
	}
	if len(results) > 0 {
		return results[0].Count
	}
	return 0
}
