package usecase

import (
	"rederinghub.io/internal/entity"
)

const (
	DEFAULT_REFERRAL_PERCENT = 3
)

func (u Usecase) CreateReferral(referrerID string, referreeID string) error {
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
