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
	"rederinghub.io/utils/constants/dao_project_voted"
	copierInternal "rederinghub.io/utils/copier"
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
		action.CanVote = user.IsVerified && project.Status == dao_project.New && user.WalletAddress != project.CreatedBy
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

func (s *Usecase) GetDAOProject(ctx context.Context, id, userWallet string) (*response.DaoProject, error) {
	request := &request.ListDaoProjectRequest{
		Pagination: entity.Pagination{
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
	project := results[0]
	userWallets := make([]string, 0, len(project.DaoProjectVoted))
	for _, voted := range project.DaoProjectVoted {
		userWallets = append(userWallets, voted.CreatedBy)
		if voted.Status == dao_project_voted.Voted {
			project.TotalVote += 1
		}
		if voted.Status == dao_project_voted.Against {
			project.TotalAgainst += 1
		}
	}
	if len(userWallets) > 0 {
		users := []*entity.Users{}
		userMap := make(map[string]*entity.Users)
		if err := s.Repo.Find(ctx, entity.Users{}.TableName(), bson.M{"wallet_address": bson.M{"$in": userWallets}}, users); err != nil {
			return nil, err
		}
		for _, user := range users {
			userMap[user.WalletAddress] = user
		}
		for _, voted := range project.DaoProjectVoted {
			if val, ok := userMap[voted.CreatedBy]; ok {
				voted.DisplayName = val.DisplayName
			}
		}
	}
	return project, nil
}
