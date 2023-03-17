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
	"rederinghub.io/utils/rediskey"
)

func (s *Usecase) ListDAOArtist(ctx context.Context, userWallet string, request *request.ListDaoArtistRequest) (*entity.Pagination, error) {
	result := &entity.Pagination{}
	redisKey := rediskey.Beauty(entity.DaoArtist{}.TableName()).WithParams("list", userWallet).WithStructHash(request, nil).String()
	if err := s.RedisV9.Get(ctx, redisKey, result); err == nil {
		return result, nil
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
				if strings.EqualFold(voted.CreatedBy, user.WalletAddress) {
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
	_ = s.RedisV9.Set(ctx, redisKey, result, time.Minute*5)
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
	_, exists := s.Repo.CheckDAOArtistAvailableByUser(ctx, userWallet)
	if exists {
		return "", errors.New("Proposal is exists")
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
	expireDay := s.Config.VoteDAOExpireDay
	if expireDay <= 0 {
		expireDay = 7
	}
	daoArtist := &entity.DaoArtist{
		CreatedBy: user.WalletAddress,
		ExpiredAt: time.Now().Add(24 * time.Duration(expireDay) * time.Hour),
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

	_ = s.RedisV9.DelPrefix(ctx, rediskey.Beauty(entity.DaoArtist{}.TableName()).WithParams("list").String())

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
	if pag.Result == nil {
		return nil, nil
	}
	results := pag.Result.([]*response.DaoArtist)
	if len(results) < 0 {
		return nil, nil
	}
	daoArtist := results[0]
	userWallets := make([]string, 0, len(daoArtist.DaoArtistVoted))
	for _, voted := range daoArtist.DaoArtistVoted {
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
		for _, voted := range daoArtist.DaoArtistVoted {
			if val, ok := userMap[voted.CreatedBy]; ok {
				voted.DisplayName = val.DisplayName
			}
		}
	}
	return daoArtist, nil
}

func (s *Usecase) GetDAOArtistByWallet(ctx context.Context, walletAddress string) (*entity.DaoArtist, error) {
	daoArtist := &entity.DaoArtist{}
	if err := s.Repo.FindOneBy(ctx, daoArtist.TableName(), bson.M{"created_by": walletAddress}, daoArtist); err != nil {
		return nil, err
	}
	return daoArtist, nil
}

func (s *Usecase) CanCreateNewProposalArtist(ctx context.Context, walletAddress string) (string, bool) {
	daoArtist, exists := s.Repo.CheckDAOArtistAvailableByUser(ctx, walletAddress)
	if exists && daoArtist != nil {
		return daoArtist.ID.Hex(), false
	}
	return "", true
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

	_ = s.RedisV9.DelPrefix(ctx, rediskey.Beauty(entity.DaoArtist{}.TableName()).WithParams("list").String())

	if req.Status != dao_artist_voted.Verify {
		return nil
	}

	_ = s.processVerifyArtist(ctx, daoArtist)

	return nil
}
func (s *Usecase) processVerifyArtist(ctx context.Context, daoArtist *entity.DaoArtist) error {
	voted := s.Repo.CountDAOArtistVoteByStatus(ctx, daoArtist.ID, dao_artist_voted.Verify)
	count := s.Config.CountVoteDAO
	if count <= 0 {
		count = 2
	}
	if voted < count {
		return nil
	}
	user := &entity.Users{}
	if err := s.Repo.FindOneBy(ctx, user.TableName(), bson.M{"wallet_address": daoArtist.CreatedBy}, user); err != nil {
		logger.AtLog.Logger.Error("Get artist failed", zap.Error(err))
		return err
	}
	if !user.ProfileSocial.TwitterVerified {
		user.ProfileSocial.TwitterVerified = true
		_, err := s.Repo.UpdateByID(ctx, user.TableName(), user.ID,
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "profile_social", Value: user.ProfileSocial},
					{Key: "updated_at", Value: time.Now()},
				}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update artist failed", zap.Error(err))
			return err
		}
	}

	if daoArtist.Status != dao_artist.Verified {
		_, err := s.Repo.UpdateByID(ctx, daoArtist.TableName(), daoArtist.ID,
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "status", Value: dao_artist.Verified},
					{Key: "updated_at", Value: time.Now()},
				}},
			})
		if err != nil {
			logger.AtLog.Logger.Error("Update DAO artist failed", zap.Error(err))
		}
	}

	return nil
}

func (s *Usecase) SetExpireYourProposalArtist(ctx context.Context, userWallet string) error {
	return s.Repo.SetExpireYourProposalArtist(ctx, userWallet)
}
