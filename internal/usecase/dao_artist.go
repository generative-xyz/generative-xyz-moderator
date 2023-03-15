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
	"rederinghub.io/utils/constants/dao_artist"
	"rederinghub.io/utils/constants/dao_artist_voted"
	copierInternal "rederinghub.io/utils/copier"
	"rederinghub.io/utils/logger"
)

func (s *Usecase) ListDAOArtist(ctx context.Context, userWallet string, request *request.ListDaoArtistRequest) (*entity.Pagination, error) {
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
	artists, total, err := s.Repo.ListDAOArtist(ctx, request)
	if err != nil {
		return nil, err
	}
	artistsResp := []*response.DaoArtist{}
	if err := copierInternal.Copy(&artistsResp, artists); err != nil {
		return nil, err
	}
	for _, artist := range artistsResp {
		action := &response.ActionDaoArtist{}
		action.CanVote = user.ProfileSocial.TwitterVerified &&
			user.WalletAddress != artist.CreatedBy &&
			!artist.Expired()
		if action.CanVote {
			for _, voted := range artist.DaoArtistVoted {
				if voted.CreatedBy == user.WalletAddress {
					action.CanVote = false
					break
				}
			}
		}
		artist.SetFields(
			artist.WithAction(action),
		)
	}
	result.Result = artistsResp
	result.Total = total
	if len(artistsResp) > 0 {
		result.Cursor = artistsResp[len(artistsResp)-1].ID
	}
	return result, nil
}

func (s *Usecase) CreateDAOArtist(ctx context.Context, userWallet string, req *request.CreateDaoArtistRequest) (string, error) {
	user := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, user.TableName(), bson.M{"wallet_address": userWallet}, user); err != nil {
		return "", err
	}
	if user.ProfileSocial.TwitterVerified {
		return "", errors.New("Haven't permission")
	}
	if req.Twitter != "" && user.ProfileSocial.Twitter == "" {
		user.ProfileSocial.Twitter = req.Twitter
		_, err := s.Repo.UpdateByID(ctx, user.TableName(), user.ID,
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "profile_social", Value: user.ProfileSocial},
					{Key: "updated_at", Value: time.Now()},
				}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update twitter artist failed", zap.Error(err))
			return "", err
		}
	}
	daoArtist := &entity.DaoArtist{
		CreatedBy: user.WalletAddress,
		ExpiredAt: time.Now().Add(24 * 7 * time.Hour),
		Status:    dao_artist.Verifying,
	}
	seqId, err := s.Repo.NextId(ctx, daoArtist.TableName())
	if err != nil {
		return "", err
	}
	daoArtist.SeqId = seqId
	daoArtist.SetID()
	daoArtist.SetCreatedAt()
	id, err := s.Repo.Create(ctx, daoArtist.TableName(), daoArtist)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *Usecase) GetDAOArtist(ctx context.Context, id, userWallet string) (*response.DaoArtist, error) {
	request := &request.ListDaoArtistRequest{
		Pagination: &entity.Pagination{
			PageSize: 1,
		},
		Id: &id,
	}
	pag, err := s.ListDAOArtist(ctx, userWallet, request)
	if err != nil {
		return nil, err
	}
	results := pag.Result.([]*response.DaoArtist)
	if len(results) < 0 {
		return nil, nil
	}
	daoArtist := results[0]
	userWallets := make([]string, 0, len(daoArtist.DaoArtistVoted))
	for _, voted := range daoArtist.DaoArtistVoted {
		userWallets = append(userWallets, voted.CreatedBy)
		if voted.Status == dao_artist_voted.Report {
			daoArtist.TotalReport += 1
		}
		if voted.Status == dao_artist_voted.Verify {
			daoArtist.TotalVerify += 1
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
		for _, voted := range daoArtist.DaoArtistVoted {
			if val, ok := userMap[voted.CreatedBy]; ok {
				voted.DisplayName = val.DisplayName
			}
		}
	}
	return daoArtist, nil
}

func (s *Usecase) VoteDAOArtist(ctx context.Context, id, userWallet string, req *request.VoteDaoArtistRequest) error {
	createdBy := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, createdBy.TableName(), bson.M{"wallet_address": userWallet}, createdBy); err != nil {
		return err
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	daoArtist := &entity.DaoArtist{}
	if err := s.Repo.FindOneBy(ctx, daoArtist.TableName(), bson.M{"_id": objectId}, daoArtist); err != nil {
		return err
	}
	if daoArtist.Expired() {
		return errors.New("Proposal was expired")
	}
	if !createdBy.ProfileSocial.TwitterVerified || strings.EqualFold(daoArtist.CreatedBy, userWallet) {
		return errors.New("Haven't permission")
	}
	daoArtistVoted := &entity.DaoArtistVoted{
		CreatedBy:   userWallet,
		DaoArtistId: daoArtist.ID,
		Status:      req.Status,
	}
	daoArtistVoted.SetID()
	daoArtistVoted.SetCreatedAt()
	_, err = s.Repo.Create(ctx, daoArtistVoted.TableName(), daoArtistVoted)
	if err != nil {
		return err
	}

	if req.Status != dao_artist_voted.Verify {
		return nil
	}
	voted := []*entity.DaoArtistVoted{}
	err = s.Repo.Find(ctx, entity.DaoArtistVoted{}.TableName(), bson.M{"dao_artist_id": daoArtist.ID, "status": dao_artist_voted.Verify}, &voted)
	if err != nil {
		return nil
	}
	count := s.Config.CountVoteDAO
	if count <= 0 {
		count = 2
	}
	if len(voted) < count || daoArtist.Status == dao_artist.Verified {
		return nil
	}
	user := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, user.TableName(), bson.M{"wallet_address": daoArtist.CreatedBy}, user); err != nil {
		logger.AtLog.Logger.Error("Get artist failed", zap.Error(err))
		return nil
	}
	user.ProfileSocial.TwitterVerified = true
	_, err = s.Repo.UpdateByID(ctx, user.TableName(), user.ID,
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "profile_social", Value: user.ProfileSocial},
				{Key: "updated_at", Value: time.Now()},
			}},
		})
	if err != nil {
		logger.AtLog.Logger.Error("Update artist failed", zap.Error(err))
		return nil
	}
	_, err = s.Repo.UpdateByID(ctx, daoArtist.TableName(), daoArtist.ID,
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "status", Value: dao_artist.Verified},
				{Key: "updated_at", Value: time.Now()},
			}},
		})
	if err != nil {
		logger.AtLog.Logger.Error("Update DAO artist failed", zap.Error(err))
	}

	return nil
}
