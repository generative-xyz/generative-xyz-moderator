package usecase

import (
	"fmt"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
)

const (
	DEFAULT_REFERRAL_PERCENT = 100
)

func (u Usecase) CreateReferral( referrerID string, referreeID string) error {
	// check if referree is referred
	count, err := u.Repo.CountReferralOfReferee(referreeID)
	if err != nil {
		u.Logger.ErrorAny("u.Repo.CountReferralOfReferee", zap.Any("FindUserByID", err))
		return err
	}
	if count > 0 {
		return fmt.Errorf("user is referred")
	}
	u.Logger.LogAny("CreateReferral", zap.Any("referrerID", referrerID), zap.Any("referreeID", referreeID))
	referrer, err := u.Repo.FindUserByID(referrerID)
	if err != nil {
		u.Logger.ErrorAny("CreateReferral", zap.Any("FindUserByID", err))
		return err
	}
	referree, err := u.Repo.FindUserByID(referreeID)
	if err != nil {
		u.Logger.ErrorAny("CreateReferral", zap.Any("FindUserByID", err))
		return err
	}

	inserted := &entity.Referral{
		ReferrerID: referrerID,
		ReferreeID: referreeID,
		Percent: DEFAULT_REFERRAL_PERCENT,
		Referrer: referrer,
		Referree: referree,
	}

	u.Logger.LogAny("CreateReferral", zap.Any("InsertReferral", inserted))
	err = u.Repo.InsertReferral(inserted)
	if err != nil {
		u.Logger.ErrorAny("CreateReferral", zap.Any("InsertReferral", err))
		return err
	}

	return nil
}

func (u Usecase) GetReferrals( req structure.FilterReferrals) (*entity.Pagination, error) {
	u.Logger.LogAny("GetReferrals", zap.Any("req", req))
	pe := &entity.FilterReferrals{}
	uuid := `63f896971aa5ce35134f391f`
	
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.ErrorAny("GetReferrals",zap.Any("copier.Copy", err))
		return nil, err
	}

	pe.ReferrerID = &uuid
	referrals, err := u.Repo.GetReferrals(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.FilterReferrals", err.Error(), err)
		return nil, err
	}

	data := referrals.Result.([]entity.Referral)
	resp := []structure.ReferalResp{}
	for _, item := range data {
		tmp := &structure.ReferalResp{}
		err = copier.Copy(tmp, item)
		if err != nil {
			u.Logger.Error("copier.Copy", err.Error(), err)
			return nil, err
		}

		volume, err := u.GetVolumeOfUser(item.Referree.WalletAddress, req.PayType)
		if err != nil   {
			tmp.ReferreeVolume = structure.ReferralVolumnResp{Amount: "0", AmountType: *req.PayType, Percent: int(item.Percent), ProjectID: "", Earn: "0", GenEarn: "0" }
		}else{
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
		resp = append(resp, *tmp)
	}
	
	referrals.Result = resp
	u.Logger.LogAny("GetReferrals", zap.Any("req", req), zap.Any("referrals",referrals), zap.Any("referrals", referrals.Total))
	return referrals, nil
}



