package usecase

import (
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateBTCWalletAddress(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("CreateConfig", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("input", input)

	walletAddress := &entity.BTCWalletAddress{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateBTCWalletAddress.Copy", err.Error(), err)
		return nil, err
	}

	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateBTCWalletAddress.InsertBtcWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) UpdateBTCWalletAddress(rootSpan opentracing.Span, input structure.ConfigData) (*entity.BTCWalletAddress, error) {
	return nil, nil
}

func (u Usecase) GetBTCWalletAddress(rootSpan opentracing.Span, input string) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("GetBTCWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)

	config, err := u.Repo.FindBtcWalletAddress(input)
	if err != nil {
		log.Error(" u.Repo.FindConfig", err.Error(), err)
		return nil, err
	}

	return config, nil
}

func (u Usecase) GetBTCWalletAddresses(rootSpan opentracing.Span, input structure.FilterConfigs) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetBTCWalletAddresses", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	f := &entity.FilterBTCWalletAddress{}
	err := copier.Copy(f, input)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	confs,  err := u.Repo.ListBtcWalletAddress(*f)
	if err != nil {
		log.Error(" u.Repo.ListBtcWalletAddress", err.Error(), err)
		return nil, err
	}

	return confs, nil

}