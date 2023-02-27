package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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

	payTypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	for _, payType := range payTypes {
		for _, item := range referrals {
			tmp := &structure.ReferalResp{}
			err = copier.Copy(tmp, item)
			if err != nil {
				u.Logger.ErrorAny("AggregateReferal", zap.Any("err", err))
				return 
			}
	
			amount := "0"
			volume, err := u.GetVolumeOfUser(item.Referree.WalletAddress, &payType)
			if err != nil {
				tmp.ReferreeVolume = structure.ReferralVolumnResp{Amount: "0", AmountType: payType, Percent: int(item.Percent), ProjectID: "", Earn: "0", GenEarn: "0" }
			}	else{
				amount = fmt.Sprintf("%d", int64(volume.Amount))
				refEarning, genEarning :=  helpers.CalculateEarning(volume.Amount, item.Percent)
				tmp.ReferreeVolume = structure.ReferralVolumnResp{
					Amount: fmt.Sprintf("%d", int(volume.Amount)), 
					AmountType: volume.ID.Paytype,
					Percent: int(item.Percent),
					ProjectID: volume.ID.ProjectID, 
					Earn: refEarning,
					GenEarn: genEarning,
				}
			}

			_ = amount
			// _, err  := u.Repo.FindVolumnByWalletAddress(item.Referrer.WalletAddress, payType)
			// if err != nil {
			// 	if errors.Is(err, mongo.ErrNoDocuments) {
			// 		ev = &entity.UserVolumn{
			// 			CreatorAddress: &item.Referrer.WalletAddress ,
			// 			PayType: &payType,
			// 			ProjectID: nil,
			// 			Amount: &amount,
			// 			Project: entity.VolumeProjectInfo{},
			// 			User: entity.VolumnUserInfo{
			// 				WalletAddress: &item.Referrer.WalletAddress,
			// 				WalletAddressBTC: &item.Referrer.WalletAddressBTC,
			// 				DisplayName:  &item.Referrer.WalletAddressBTC,
			// 				Avatar: &item.Referrer.Avatar,
			// 			},
			// 		}

					
			// 	}
			// }else{
				
			// }
		}
	}
}

func (u Usecase) GetVolumeOfUser(walletAddress string, amountType *string) (*entity.AggregateAmount, error) {
	
	amount, err :=  u.Repo.AggregateAmount(entity.FilterVolume{
		CreatorAddress: &walletAddress,
		AmountType: amountType,
	})
	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}