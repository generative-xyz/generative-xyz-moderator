package usecase

import (
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

const (
	DEFAULT_REFERRAL_PERCENT = 3
)

func (u Usecase) CreateReferral(rootSpan opentracing.Span, referrerID string, referreeID string) error {
	span, log := u.StartSpan("CreateReferral", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

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

func (u Usecase) GetReferrals(rootSpan opentracing.Span, req structure.FilterReferrals) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetReferrals", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	pe := &entity.FilterReferrals{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	referrals, err := u.Repo.GetReferrals(*pe)
	if err != nil {
		log.Error("u.Repo.FilterReferrals", err.Error(), err)
		return nil, err
	}

	log.SetData("referrals", referrals.Total)
	return referrals, nil
}

