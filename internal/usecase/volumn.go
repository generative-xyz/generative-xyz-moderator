package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) AggregateVolumns() {
	payTypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	for _, payType := range payTypes {
		u.Logger.LogAny("AggregateVolumns",zap.Any("payType", payType))
		u.AggregateVolumn(payType)
	}
}

func (u Usecase) AggregateVolumn(payType string)   {
	data, err := u.Repo.AggregateVolumn(payType)
	if err != nil {
		u.Logger.ErrorAny("CreateVolume", zap.Any("err", err))
		return
	}
	
	u.Logger.LogAny("AggregateVolumn",zap.Any("payType", payType), zap.Any("data",data))
	processed := 0
	for _, item := range data {
		go func(item entity.AggregateProjectItemResp){
			u.Logger.LogAny("aggregateVolumn",zap.Any("item", item))
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
				amount :=  fmt.Sprintf("%d", int(item.Amount))
				earning, gearning := helpers.CalculateEarning(item.Amount, int32(utils.PERCENT_EARNING))
				if errors.Is(err, mongo.ErrNoDocuments) {
					//v := entity.FilterVolume
					ev := &entity.UserVolumn{
						CreatorAddress: &creatorID,
						PayType: &item.Paytype,
						ProjectID: &pID,
						Amount: &amount,
						Earning: &earning,
						GenEarning: &gearning,
						Minted: item.Minted,
						MintPrice: item.MintPrice,
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
				amount :=  fmt.Sprintf("%d", int(item.Amount))
				if amount != *ev.Amount  {
					earning, gearning := helpers.CalculateEarning(item.Amount, int32(utils.PERCENT_EARNING))
					_, err := u.Repo.UpdateVolumnAmount(ev.UUID, amount, earning, gearning)
					if err != nil {
						u.Logger.ErrorAny("UpdateVolumnAmount",zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
						return
					}
				}
				
				if item.Minted != ev.Minted  {
					_, err := u.Repo.UpdateVolumnMinted(ev.UUID, item.Minted)
					if err != nil {
						u.Logger.ErrorAny("UpdateVolumnAmount",zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
						return
					}
				}
				
				if int(item.MintPrice) != int(ev.MintPrice)  {
					_, err := u.Repo.UpdateVolumMintPrice(ev.UUID, item.MintPrice)
					if err != nil {
						u.Logger.ErrorAny("UpdateVolumnAmount",zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
						return
					}
				}
		}
		}(item)

		if processed % 10 == 0 {
			time.Sleep(2 * time.Second)
		}

		processed ++
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


func (u Usecase) GetEarningOfUser(walletAddress string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": 
		bson.M{"creatorAddress": "$creatorAddress", "payType": "$payType"}, 
		"amount": bson.M{"$sum": bson.M{"$toDouble": "$earning"}},
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
		"earning": bson.M{"$sum": bson.M{"$toDouble": "$earning"}},
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