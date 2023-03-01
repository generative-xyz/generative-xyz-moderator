package usecase

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) AggregateVolumn() {

 	data, err := u.Repo.AggregateVolumn()
	if err != nil {
		u.Logger.ErrorAny("CreateVolume", zap.Any("err", err))
	}

	for _, item := range data {
		pID := strings.ToLower(item.ProjectID)
		p, err := u.Repo.FindProjectByTokenID(pID)
		if err != nil {
			u.Logger.ErrorAny("FindProjectByTokenID",zap.String("item.ProjectID", item.ProjectID), zap.Any("err", err))
			return
		}

		creatorID := strings.ToLower(p.CreatorAddrr)
		usr, err := u.Repo.FindUserByWalletAddress(creatorID)
		if err != nil {
			u.Logger.ErrorAny("FindUserByWalletAddress",zap.String("p.CreatorAddrr", creatorID), zap.Any("err", err))
			return
		}

		ev, err := u.Repo.FindVolumn(pID, item.Paytype) 
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				//v := entity.FilterVolume
				ev := &entity.UserVolumn{
					CreatorAddress: &creatorID,
					PayType: &item.Paytype,
					ProjectID: &pID,
					Amount: &item.Amount,
					Project: entity.VolumeProjectInfo{
						Name: p.Name,
						TokenID: p.TokenID,
						Thumnail: p.Thumbnail,
					},
					User: entity.VolumnUserInfo{
						WalletAddress: &p.CreatorAddrr,
						WalletAddressBTC: &usr.WalletAddressBTC,
						DisplayName: &usr.DisplayName,
						Avatar: &usr.Avatar,
					},
				}

				err = u.Repo.CreateVolumn(ev)
				if err != nil {
					u.Logger.ErrorAny("CreateVolumn",zap.Any("ev", ev), zap.Any("err", err))
					return
				}
			}
		}else{
			if item.Amount != *ev.Amount {
				_, err := u.Repo.UpdateVolumnAmount(ev.UUID, item.Amount)
				if err != nil {
					u.Logger.ErrorAny("UpdateVolumnAmount",zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
					return
				}
			}
		}
	}
}

func (u Usecase) AggregateReferal() {

	referrals, err := u.Repo.GetAllReferrals(entity.FilterReferrals{})
	if err != nil {
		u.Logger.ErrorAny("AggregateReferal", zap.Any("err", err))
		return
	}

	paytypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	for _, referral := range referrals  {
		vol := make(map[string]entity.ReferreeVolumn)
		for _, paytype := range paytypes{
			
			volume, err := u.GetVolumeOfUser(referral.Referree.WalletAddress, &paytype)
			if err != nil {
				vol[paytype] = entity.ReferreeVolumn{
					Amount: "0",
					AmountType: paytype,
					Earn: "0",
					GenEarn: "0",
					RemainingEarn: "0",
				}
			}else{
				refEarning, genEarning :=  helpers.CalculateEarning(volume.Amount, referral.Percent)
				remaining := referral.ReferreeVolumn[paytype].RemainingEarn
				if remaining == "" {
					remaining = "0"
				}

				vol[paytype] = entity.ReferreeVolumn{
					Amount: fmt.Sprintf("%d", int(volume.Amount)),
					AmountType: paytype,
					Earn: refEarning,
					GenEarn: genEarning,
					RemainingEarn: remaining,
				}
			}
		}
		referral.ReferreeVolumn = vol
		_, err = u.Repo.UpdateReferral(referral.UUID, &referral)
		if err != nil {
			u.Logger.ErrorAny("AggregateReferal",zap.Error(err))
			return 
		}
	}
	_ = referrals
}

func (u Usecase) GetVolumeOfUser(walletAddress string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": 
		bson.M{"creatorAddress": "$creatorAddress", "payType": "$payType"}, 
		"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
	}}

	amount, err :=  u.Repo.AggregateAmount(entity.FilterVolume{
		CreatorAddress: &walletAddress,
		AmountType: amountType,
	}, group)
	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}

func (u Usecase) GetVolumeOfProject(projectID string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": 
		bson.M{"projectID": "$projectID", "payType": "$payType"}, 
		"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
	}}

	amount, err :=  u.Repo.AggregateAmount(entity.FilterVolume{
		ProjectID: &projectID,
		AmountType: amountType,
	}, group)

	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}