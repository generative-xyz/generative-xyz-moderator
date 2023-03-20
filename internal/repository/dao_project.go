package repository

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
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
	filters["$or"] = bson.A{
		bson.M{"expired_at": bson.M{"$gt": time.Now()}},
		bson.M{"status": dao_project.Executed},
	}
	if request.Id != nil {
		id, err := primitive.ObjectIDFromHex(*request.Id)
		if err != nil {
			return nil, 0, err
		}
		filters["_id"] = id
	}
	if request.SeqId != nil {
		filters["seq_id"] = *request.SeqId
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
		search := bson.A{
			bson.M{"project_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
			bson.M{"user_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
		}
		if seqId, err := strconv.Atoi(*request.Keyword); err == nil {
			search = append(search, bson.M{"seq_id": seqId})
		} else {
			if id, err := primitive.ObjectIDFromHex(*request.Keyword); err == nil {
				search = append(search, bson.M{"_id": id})
			}
		}
		filterSearch["$or"] = search
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
func (s Repository) CheckDAOProjectAvailableByProjectId(ctx context.Context, projectId primitive.ObjectID) bool {
	daoProject := &entity.DaoProject{}
	if err := s.FindOneBy(ctx, daoProject.TableName(), bson.M{
		"project_id": projectId,
		"$or": bson.A{
			bson.M{"expired_at": bson.M{"$gt": time.Now()}},
			bson.M{"status": dao_project.Executed},
		},
	}, daoProject); err != nil {
		return false
	}
	return true
}
func (s Repository) CheckDAOProjectAvailableByUser(ctx context.Context, userWallet string, projectId primitive.ObjectID) (*entity.DaoProject, bool) {
	daoProject := &entity.DaoProject{}
	if err := s.FindOneBy(ctx, daoProject.TableName(), bson.M{
		"created_by": userWallet,
		"project_id": projectId,
		"$or": bson.A{
			bson.M{"expired_at": bson.M{"$gt": time.Now()}},
			bson.M{"status": dao_project.Executed},
		},
	}, daoProject); err != nil {
		return nil, false
	}
	return daoProject, true
}
func (s Repository) CountDAOProjectVoteByStatus(ctx context.Context, daoProjectId primitive.ObjectID, status dao_project_voted.Status) int {
	match := bson.M{"$match": bson.M{
		"dao_project_id": daoProjectId,
		"status":         status,
	}}
	group := bson.M{
		"$group": bson.M{
			"_id":   "$dao_project_id",
			"count": bson.M{"$sum": 1},
		},
	}
	cur, err := s.DB.Collection(entity.DaoProjectVoted{}.TableName()).Aggregate(ctx, bson.A{match, group})
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
