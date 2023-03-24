package usecase

import (
	"fmt"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

const (
	DEFAULT_REFERRAL_PERCENT = 100
)

func (u Usecase) CreateReferral( referrerID string, referreeID string) error {
	// check if referree is referred
	count, err := u.Repo.CountReferralOfReferee(referreeID)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.CountReferralOfReferee", zap.Any("FindUserByID", err))
		return err
	}
	if count > 0 {
		return fmt.Errorf("user is referred")
	}
	logger.AtLog.Logger.Info("CreateReferral", zap.Any("referrerID", referrerID), zap.Any("referreeID", referreeID))
	referrer, err := u.Repo.FindUserByID(referrerID)
	if err != nil {
		logger.AtLog.Logger.Error("CreateReferral", zap.Any("FindUserByID", err))
		return err
	}
	referree, err := u.Repo.FindUserByID(referreeID)
	if err != nil {
		logger.AtLog.Logger.Error("CreateReferral", zap.Any("FindUserByID", err))
		return err
	}

	inserted := &entity.Referral{
		ReferrerID: referrerID,
		ReferreeID: referreeID,
		Percent: DEFAULT_REFERRAL_PERCENT,
		Referrer: referrer,
		Referree: referree,
	}

	logger.AtLog.Logger.Info("CreateReferral", zap.Any("InsertReferral", inserted))
	err = u.Repo.InsertReferral(inserted)
	if err != nil {
		logger.AtLog.Logger.Error("CreateReferral", zap.Any("InsertReferral", err))
		return err
	}

	return nil
}

func (u Usecase) GetReferrals( req structure.FilterReferrals) (*entity.Pagination, error) {
	logger.AtLog.Logger.Info("GetReferrals", zap.Any("req", req))
	pe := &entity.FilterReferrals{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Logger.Error("GetReferrals",zap.Any("copier.Copy", err))
		return nil, err
	}
	referrals, err := u.Repo.GetReferrals(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.FilterReferrals", err.Error(), err)
		return nil, err
	}

	//spew.Dump(referrals)
	data := referrals.Result.([]entity.Referral)

	resp := []structure.ReferalResp{}
	for _, item := range data {
		tmp := &structure.ReferalResp{}
		err = copier.Copy(tmp, item)
		if err != nil {
			u.Logger.Error("copier.Copy", err.Error(), err)
			return nil, err
		}

		wdType := string(entity.WithDrawReferal)
		latestWd, _ := u.Repo.GetLastWithdraw(entity.FilterWithdraw{
			WalletAddress: &item.Referrer.WalletAddress,
			WithdrawItemID: &item.ReferreeID, 
			PaymentType:    req.PayType,
			WithdrawType: &wdType,
		})

		status := entity.StatusWithdraw_Available
		if  latestWd != nil {
			status = latestWd.Status
			if status == entity.StatusWithdraw_Approve {
				status = entity.StatusWithdraw_Available
			}
			
			if status == entity.StatusWithdraw_Reject {
				status = entity.StatusWithdraw_Available
			}
		}

		tmp.ReferreeVolume = structure.ReferralVolumnResp{
			Earn: item.ReferreeVolumn[*req.PayType].Earn,
			Amount: item.ReferreeVolumn[*req.PayType].Amount,
			GenEarn: item.ReferreeVolumn[*req.PayType].GenEarn,
			AmountType: *req.PayType,
			Status: status,
		}
		
		resp = append(resp, *tmp)
	}
	
	referrals.Result = resp
	logger.AtLog.Logger.Info("GetReferrals", zap.Any("req", req), zap.Any("referrals",referrals), zap.Any("referrals", referrals.Total))
	return referrals, nil
}
