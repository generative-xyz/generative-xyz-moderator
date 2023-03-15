package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
	copierInternal "rederinghub.io/utils/copier"
	"rederinghub.io/utils/logger"
)

func (s *Usecase) ListDAOProject(ctx context.Context, userWallet string, request *request.ListDaoProjectRequest) (*entity.Pagination, error) {
	result := &entity.Pagination{
		PageSize: request.PageSize,
		Result:   make([]*response.DaoProject, 0),
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
		action.CanVote = user.IsVerified &&
			user.WalletAddress != project.CreatedBy &&
			!project.Expired()
		if action.CanVote {
			for _, voted := range project.DaoProjectVoted {
				if voted.CreatedBy == user.WalletAddress {
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
		Status:    dao_project.Voting,
	}
	seqId, err := s.Repo.NextId(ctx, daoProject.TableName())
	if err != nil {
		return "", err
	}
	daoProject.SeqId = seqId
	daoProject.SetID()
	daoProject.SetCreatedAt()
	id, err := s.Repo.Create(ctx, daoProject.TableName(), daoProject)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
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
	results := pag.Result.([]*response.DaoProject)
	if len(results) < 0 {
		return nil, nil
	}
	daoProject := results[0]
	userWallets := make([]string, 0, len(daoProject.DaoProjectVoted))
	for _, voted := range daoProject.DaoProjectVoted {
		userWallets = append(userWallets, voted.CreatedBy)
		if voted.Status == dao_project_voted.Voted {
			daoProject.TotalVote += 1
		}
		if voted.Status == dao_project_voted.Against {
			daoProject.TotalAgainst += 1
		}
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
	if !createdBy.IsVerified || strings.EqualFold(daoProject.CreatedBy, userWallet) {
		return errors.New("haven't permission")
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
	if req.Status != dao_project_voted.Voted {
		return nil
	}
	go func() {
		voted := []*entity.DaoProjectVoted{}
		err = s.Repo.Find(ctx,
			entity.DaoProjectVoted{}.TableName(),
			bson.M{
				"dao_project_id": daoProject.ID,
				"status":         dao_project_voted.Voted,
			}, &voted)
		if err != nil {
			return
		}
		count := s.Config.CountVoteDAO
		if count <= 0 {
			count = 2
		}
		if len(voted) < count || daoProject.Status == dao_project.Executed {
			return
		}
		project := &entity.Projects{}
		if err := s.Repo.FindOneBy(ctx, project.TableName(), bson.M{"_id": daoProject.ProjectId}, project); err != nil {
			logger.AtLog.Logger.Error("Get project failed", zap.Error(err))
			return
		}
		_, err = s.Repo.UpdateByID(ctx, project.TableName(), project.ID,
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "isHidden", Value: false}}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update project failed", zap.Error(err))
			return
		}
		_, err = s.Repo.UpdateByID(ctx, daoProject.TableName(), daoProject.ID,
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "status", Value: dao_project.Executed}}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update DAO project failed", zap.Error(err))
		}
	}()
	return nil
}
