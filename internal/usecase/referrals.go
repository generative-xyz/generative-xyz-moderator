package usecase

import (
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

const (
	DEFAULT_REFERRAL_PERCENT = 3
)

func (u Usecase) CreateReferral( referrerID string, referreeID string) error {


	referrer, err := u.Repo.FindUserByID(referrerID)
	if err != nil {
		return err
	}
	referree, err := u.Repo.FindUserByID(referreeID)
	if err != nil {
		return err
	}

	err = u.Repo.InsertReferral(&entity.Referral{
		ReferrerID: referrerID,
		ReferreeID: referreeID,
		Percent: DEFAULT_REFERRAL_PERCENT,
		Referrer: referrer,
		Referree: referree,
	})

	if err != nil {
		return err
	}

	return nil
}

func (u Usecase) GetReferrals( req structure.FilterReferrals) (*entity.Pagination, error) {


	pe := &entity.FilterReferrals{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
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

