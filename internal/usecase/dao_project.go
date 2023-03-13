package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project"
)

func (s *Usecase) ListDAOProject(ctx context.Context, userWallet string, req *entity.Pagination) (*entity.Pagination, error) {
	return nil, nil
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
