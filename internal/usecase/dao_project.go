package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project"
	copierInternal "rederinghub.io/utils/copier"
)

func (s *Usecase) ListDAOProject(ctx context.Context, userWallet string, request *request.ListDaoProjectRequest) (*entity.Pagination, error) {
	result := &entity.Pagination{PageSize: request.PageSize}
	limit := int64(100)
	filters := make(bson.M)
	filterIdOperation := "$lt"
	sorts := bson.M{
		"$sort": bson.D{{Key: "_id", Value: -1}},
	}
	match := bson.M{"$match": filters}
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
		"$addFields": bson.M{"user_name": "$user.name"},
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
	if request.Status != nil {
		filters["status"] = request.Status
	}
	if request.Cursor != "" {
		if id, err := primitive.ObjectIDFromHex(request.Cursor); err == nil {
			filters["_id"] = bson.M{filterIdOperation: id}
		}
	}
	if request.Keyword != nil {
		keyword := *request.Keyword
		filters["$or"] = bson.A{
			bson.M{"project_name": primitive.Regex{
				Pattern: keyword,
				Options: "i",
			}},
			bson.M{"user_name": primitive.Regex{
				Pattern: keyword,
				Options: "i",
			}},
		}
	}
	projects := []*entity.DaoProject{}
	total, err := s.Repo.Aggregation(ctx,
		entity.DaoProject{}.TableName(),
		0,
		limit,
		&projects,
		match,
		lookupProject,
		unwindProject,
		lookupUser,
		unwindUser,
		lookupDaoProjectVoted,
		addProjectName,
		addUserName,
		sorts)
	if err != nil {
		return nil, err
	}
	projectsResp := []*response.DaoProject{}
	if err := copierInternal.Copy(&projectsResp, projects); err != nil {
		return nil, err
	}
	result.Result = projectsResp
	result.Total = total
	if len(projectsResp) > 0 {
		result.Cursor = projectsResp[len(projectsResp)-1].ID
	}
	return result, nil
}

func (s *Usecase) CreateDAOProject(ctx context.Context, req *request.CreateDaoProjectRequest) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(req.ProjectId)
	if err != nil {
		return "", err
	}
	project := &entity.Projects{}
	if err := s.Repo.FindOneBy(ctx, project.TableName(), bson.M{"_id": objectId}, project); err != nil {
		return "", err
	}
	if !strings.EqualFold(project.CreatorAddrr, req.CreatedBy) {
		return "", errors.New("haven't permission")
	}
	createdBy := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, createdBy.TableName(), bson.M{"wallet_address": req.CreatedBy}, createdBy); err != nil {
		return "", err
	}
	daoProject := &entity.DaoProject{
		CreatedBy: req.CreatedBy,
		ProjectId: project.ID,
		ExpiredAt: time.Now().Add(24 * 7 * time.Hour),
		Status:    dao_project.New,
	}
	daoProject.SetID()
	daoProject.SetCreatedAt()
	id, err := s.Repo.Create(ctx, daoProject.TableName(), daoProject)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}
