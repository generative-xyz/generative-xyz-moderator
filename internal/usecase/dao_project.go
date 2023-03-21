package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
	copierInternal "rederinghub.io/utils/copier"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/rediskey"
)

func (s *Usecase) ListDAOProject(ctx context.Context, userWallet string, request *request.ListDaoProjectRequest) (*entity.Pagination, error) {
	result := &entity.Pagination{}
	redisKey := rediskey.Beauty(entity.DaoProject{}.TableName()).
		WithParams("list", userWallet).
		WithStructHash(request, &hashstructure.HashOptions{IgnoreZeroValue: true}).
		String()
	if err := s.RedisV9.Get(ctx, redisKey, result); err == nil {
		return result, nil
	}
	user := &entity.Users{}
	if userWallet != "" {
		if err := s.Repo.FindOneBy(ctx, user.TableName(), bson.M{"wallet_address": userWallet}, user); err != nil {
			return nil, err
		}
	}
	projects, total, err := s.Repo.ListDAOProject(ctx, request)
	if err != nil {
		return nil, err
	}
	projectsResp := []*response.DaoProject{}
	if err := copierInternal.Copy(&projectsResp, projects); err != nil {
		return nil, err
	}
	for _, project := range projectsResp {
		action := &response.ActionDaoProject{}
		action.CanVote = user.ProfileSocial.TwitterVerified &&
			user.WalletAddress != project.CreatedBy &&
			!project.Expired()
		if action.CanVote {
			for _, voted := range project.DaoProjectVoted {
				if strings.EqualFold(voted.CreatedBy, user.WalletAddress) {
					action.CanVote = false
					break
				}
			}
		}
		project.SetFields(
			project.WithAction(action),
		)
	}
	result.Result = projectsResp
	result.Total = total
	if len(projectsResp) > 0 {
		result.Cursor = projectsResp[len(projectsResp)-1].ID
	}
	_ = s.RedisV9.Set(ctx, redisKey, result, time.Minute*5)
	return result, nil
}

func (s *Usecase) CreateDAOProject(ctx context.Context, req *request.CreateDaoProjectRequest) ([]string, error) {
	createdBy := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, createdBy.TableName(), bson.M{"wallet_address": req.CreatedBy}, createdBy); err != nil {
		return nil, err
	}
	objectIds, err := utils.StringsToObjects(req.ProjectIds)
	if err != nil {
		return nil, err
	}
	expireDay := s.Config.VoteDAOExpireDay
	if expireDay <= 0 {
		expireDay = 7
	}
	daoProjects := make([]interface{}, 0, len(objectIds))
	for _, objectId := range objectIds {
		project := &entity.Projects{}
		if err := s.Repo.FindOneBy(ctx, project.TableName(), bson.M{"_id": objectId}, project); err != nil {
			return nil, err
		}
		if !strings.EqualFold(project.CreatorAddrr, req.CreatedBy) {
			return nil, errors.New("Haven't permission")
		}
		if s.Repo.CheckDAOProjectAvailableByProjectId(ctx, project.ID) {
			return nil, errors.New("Proposal is exists")
		}
		daoProject := &entity.DaoProject{
			CreatedBy: req.CreatedBy,
			ProjectId: project.ID,
			ExpiredAt: time.Now().Add(24 * time.Duration(expireDay) * time.Hour),
			Status:    dao_project.Voting,
		}
		seqId, err := s.Repo.NextId(ctx, daoProject.TableName())
		if err != nil {
			return nil, err
		}
		daoProject.SeqId = seqId
		daoProject.SetID()
		daoProject.SetCreatedAt()
		daoProjects = append(daoProjects, daoProject)
	}
	ids, err := s.Repo.CreateMany(ctx, entity.DaoProject{}.TableName(), daoProjects)
	if err != nil {
		return nil, err
	}

	_ = s.RedisV9.DelPrefix(ctx, rediskey.Beauty(entity.DaoProject{}.TableName()).WithParams("list").String())

	return utils.ObjectsToHex(ids), nil
}

func (s *Usecase) GetDAOProject(ctx context.Context, id, userWallet string) (*response.DaoProject, error) {
	request := &request.ListDaoProjectRequest{
		Pagination: &entity.Pagination{
			PageSize: 1,
		},
		Id: &id,
	}
	pag, err := s.ListDAOProject(ctx, userWallet, request)
	if err != nil {
		return nil, err
	}
	if pag.Result == nil {
		return nil, nil
	}
	results, ok := pag.Result.([]*response.DaoProject)
	if !ok || len(results) <= 0 {
		return nil, nil
	}
	daoProject := results[0]
	userWallets := make([]string, 0, len(daoProject.DaoProjectVoted))
	for _, voted := range daoProject.DaoProjectVoted {
		userWallets = append(userWallets, voted.CreatedBy)
	}
	if len(userWallets) > 0 {
		users := []*entity.Users{}
		userMap := make(map[string]*entity.Users)
		if err := s.Repo.Find(ctx, entity.Users{}.TableName(), bson.M{"wallet_address": bson.M{"$in": userWallets}}, &users); err != nil {
			return nil, err
		}
		for _, user := range users {
			userMap[user.WalletAddress] = user
		}
		for _, voted := range daoProject.DaoProjectVoted {
			if val, ok := userMap[voted.CreatedBy]; ok {
				voted.DisplayName = val.DisplayName
			}
		}
	}
	return daoProject, nil
}

func (s *Usecase) GetLastDAOProjectByProjectId(ctx context.Context, projectId primitive.ObjectID) (*entity.DaoProject, error) {
	daoProject := &entity.DaoProject{}
	opts := &options.FindOneOptions{}
	opts.SetSort(bson.M{"created_at": -1})
	if err := s.Repo.FindOneBy(ctx, daoProject.TableName(), bson.M{"project_id": projectId}, daoProject, &options.FindOneOptions{}, opts); err != nil {
		return nil, err
	}
	return daoProject, nil
}

func (s *Usecase) ListDAOProjectsByProjectId(ctx context.Context, projectId string) ([]*entity.DaoProject, error) {
	objectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, err
	}
	daoProjects := []*entity.DaoProject{}
	if err := s.Repo.Find(ctx, entity.DaoProject{}.TableName(), bson.M{"project_id": objectId}, &daoProjects); err != nil {
		return nil, err
	}
	return daoProjects, nil
}

func (s *Usecase) IsProjectReviewing(ctx context.Context, projectId string) bool {
	daoProjects, err := s.ListDAOProjectsByProjectId(ctx, projectId)
	if err != nil {
		return false
	}
	if len(daoProjects) <= 0 {
		return false
	}
	isReviewing := true
	for _, daoProject := range daoProjects {
		if daoProject.Status == dao_project.Executed {
			isReviewing = false
			break
		}
	}
	return isReviewing
}

func (s *Usecase) CheckDAOProjectAvailableByUser(ctx context.Context, walletAddress string, projectId primitive.ObjectID) (*entity.DaoProject, bool) {
	return s.Repo.CheckDAOProjectAvailableByUser(ctx, walletAddress, projectId)
}

func (s *Usecase) SetExpireAvailableDAOProject(ctx context.Context, projectId primitive.ObjectID) error {
	_, err := s.Repo.UpdateMany(ctx, entity.DaoProject{}.TableName(),
		bson.M{
			"project_id": projectId,
			"$and": bson.A{
				bson.M{"expired_at": bson.M{"$gt": time.Now()}},
				bson.M{"status": dao_project.Voting},
			},
		},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "expired_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			}},
		})

	_ = s.RedisV9.DelPrefix(ctx, rediskey.Beauty(entity.DaoProject{}.TableName()).WithParams("list").String())

	return err
}

func (s *Usecase) VoteDAOProject(ctx context.Context, id, userWallet string, req *request.VoteDaoProjectRequest) error {
	createdBy := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, createdBy.TableName(), bson.M{"wallet_address": userWallet}, createdBy); err != nil {
		return err
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	daoProject := &entity.DaoProject{}
	if err := s.Repo.FindOneBy(ctx, daoProject.TableName(), bson.M{"_id": objectId}, daoProject); err != nil {
		return err
	}
	if daoProject.Expired() {
		return errors.New("Proposal was expired")
	}
	if !createdBy.ProfileSocial.TwitterVerified || strings.EqualFold(daoProject.CreatedBy, userWallet) {
		return errors.New("Haven't permission")
	}
	voteDaoProject := &entity.DaoProjectVoted{
		CreatedBy:    userWallet,
		DaoProjectId: daoProject.ID,
		Status:       req.Status,
	}
	voteDaoProject.SetID()
	voteDaoProject.SetCreatedAt()
	_, err = s.Repo.Create(ctx, voteDaoProject.TableName(), voteDaoProject)
	if err != nil {
		return err
	}

	_ = s.RedisV9.DelPrefix(ctx, rediskey.Beauty(entity.DaoProject{}.TableName()).WithParams("list").String())

	if req.Status != dao_project_voted.Voted {
		return nil
	}

	_ = s.processEnableProject(ctx, daoProject)

	return nil
}
func (s *Usecase) processEnableProject(ctx context.Context, daoProject *entity.DaoProject) error {
	voted := s.Repo.CountDAOProjectVoteByStatus(ctx, daoProject.ID, dao_project_voted.Voted)
	count := s.Config.CountVoteDAO
	if count <= 0 {
		count = 2
	}
	if voted < count {
		return nil
	}
	project := &entity.Projects{}
	if err := s.Repo.FindOneBy(ctx, project.TableName(), bson.M{"_id": daoProject.ProjectId}, project); err != nil {
		logger.AtLog.Logger.Error("Get project failed", zap.Error(err))
		return err
	}

	if project.IsHidden {
		_, err := s.Repo.UpdateByID(ctx, project.TableName(), project.ID,
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "isHidden", Value: false},
					{Key: "updated_at", Value: time.Now()},
				}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update project failed", zap.Error(err))
			return err
		}
	}

	if daoProject.Status != dao_project.Executed {
		_, err := s.Repo.UpdateByID(ctx, daoProject.TableName(), daoProject.ID,
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "status", Value: dao_project.Executed},
					{Key: "updated_at", Value: time.Now()},
				}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update DAO project failed", zap.Error(err))
		}
		go s.NotifyCreateNewProjectToDiscord(project, &project.CreatorProfile, false, daoProject.UUID)
	}
	return nil
}

func (s *Usecase) ListYourProjectsIsHidden(ctx context.Context, userWallet string, req *request.ListProjectHiddenRequest) (*entity.Pagination, error) {
	result := &entity.Pagination{
		PageSize: req.PageSize,
		Result:   make([]*response.DaoProject, 0),
	}
	limit := int64(100)
	filters := make(bson.M)
	sorts := bson.M{
		"$sort": bson.D{{Key: "_id", Value: -1}},
	}
	matchFilters := bson.M{"$match": filters}
	lookupDaoProject := bson.M{
		"$lookup": bson.M{
			"from":         "dao_project",
			"localField":   "_id",
			"foreignField": "project_id",
			"as":           "dao_project",
		},
	}
	if req.PageSize > 0 && req.PageSize <= limit {
		limit = req.PageSize
	}
	if req.Keyword != nil {
		filters["$or"] = bson.A{
			bson.M{"name": primitive.Regex{
				Pattern: *req.Keyword,
				Options: "i",
			}},
		}
	}
	filters["creatorAddress"] = userWallet
	filters["isHidden"] = true
	filters["isSynced"] = true
	if req.Cursor != "" {
		if id, err := primitive.ObjectIDFromHex(req.Cursor); err == nil {
			filters["_id"] = bson.M{"$lt": id}
		}
	}
	addFieldsCount := bson.M{
		"$addFields": bson.M{
			"dao_project_not_expire": bson.M{
				"$filter": bson.M{
					"input": "$dao_project",
					"cond": bson.M{
						"$gt": []interface{}{"$$this.expired_at", time.Now()},
					},
				},
			},
			"dao_project_is_voted": bson.M{
				"$filter": bson.M{
					"input": "$dao_project",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", dao_project.Executed},
					},
				},
			},
		},
	}
	addFieldsTotal := bson.M{
		"$addFields": bson.M{
			"total_dao_project_not_expire": bson.M{"$size": "$dao_project_not_expire"},
			"total_dao_project_is_voted":   bson.M{"$size": "$dao_project_is_voted"},
		},
	}
	matchCount := bson.M{
		"$match": bson.M{
			"total_dao_project_not_expire": bson.M{"$lt": 1},
			"total_dao_project_is_voted":   bson.M{"$lt": 1},
		},
	}
	projects := []*entity.Projects{}
	total, err := s.Repo.Aggregation(ctx,
		entity.Projects{}.TableName(),
		0,
		limit,
		&projects,
		matchFilters,
		lookupDaoProject,
		addFieldsCount,
		addFieldsTotal,
		matchCount,
		sorts)
	if err != nil {
		return nil, err
	}
	projectsResp := []*response.ProjectForDaoProject{}
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
