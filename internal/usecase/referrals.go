package usecase

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

const (
	DEFAULT_REFERRAL_PERCENT = 300
)

func (u Usecase) CreateReferral( referrerID string, referreeID string) error {

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
	pe := &entity.FilterReferrals{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.ErrorAny("GetReferrals",zap.Any("copier.Copy", err))
		return nil, err
	}

	referrals, err := u.Repo.GetReferrals(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.FilterReferrals", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("referrals", referrals.Total)
	return referrals, nil
}

