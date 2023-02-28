package usecase

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateWithdraw(walletAddress string, data structure.WithDrawRequest) ([]entity.Withdraw, error) {
	resp := []entity.Withdraw{}
	u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("input", data))
	for _, wr := range data.Items {
		volumeAmount := 0.0
		widthDrawAmount := 0.0
		refAmount := 0.0
		//totalEarning := (refAmount + refAmount) - widthDrawAmount
		// (refAmount + refAmount) is pushed into volumn by crontab
		//TODO - calculate refAmount
		status := entity.StatusWithdraw_Done
		wdf := &entity.FilterWithdraw{
			WalletAddress: &walletAddress,
			PaymentType: &wr.PaymentType,
			Status: &status,
		}
		wd, err := u.Repo.AggregateWithDrawByUser(wdf)
		if len(wd) > 0 {
			widthDrawAmount = wd[0].Amount
		}

		f := &entity.Withdraw{}
		err = copier.Copy(f, wr)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw.Copy", zap.Any("input", data), zap.Error(err))
			return nil, err
		}

		u.Logger.LogAny("CreateWithdraw.FilterVolume", zap.String("walletAddress", walletAddress))
		volumes, _ := u.CreatorVolume(walletAddress, f.PayType)
		
		fr := entity.FilterReferrals{
			ReferrerAddress: &walletAddress,
		}
		
		refs, _ := u.Repo.GetReferral(fr)
		if len(refs) > 0 {
			for _, ref := range refs {
				tmp, err := strconv.ParseFloat(ref.ReferreeVolumn[wr.PaymentType].Amount, 10)
				if err == nil {
					refAmount +=   tmp
				}
				
			}
		}
		u.Logger.LogAny("CreateWithdraw.volumes", zap.String("walletAddress", walletAddress), zap.Any("volumes", volumes))
		//check referal amount
		if volumes == nil {
			volumeAmount = 0
		}else{
			f, err := strconv.ParseFloat(volumes.Amount, 10)
			if err != nil {
				volumeAmount = 0
			}else{
				volumeAmount = f
			}
		}
		
		totalEarnings := refAmount + volumeAmount
		earning := totalEarnings - widthDrawAmount
		if earning <= 0 {
			err = errors.New("Not enough balance")
			u.Logger.ErrorAny("CreateWithdraw", zap.Float64("earning", earning) , zap.String("walletAddress", walletAddress),  zap.Any("volumeAmount", volumeAmount), zap.Error(err))
			return nil, err
		}

		requestEarnings, err := strconv.ParseFloat(wr.Amount, 10)
		if err != nil {
			requestEarnings = 0
		}

		if requestEarnings  > earning {
			requestEarnings = earning
		}

		f.WalletAddress = walletAddress
		f.PayType = wr.PaymentType
		f.Amount = fmt.Sprintf("%d", int(requestEarnings))
		f.EarningReferal = fmt.Sprintf("%d", int(refAmount))
		f.EarningVolume = fmt.Sprintf("%d", int(volumeAmount))
		f.TotalEarnings = fmt.Sprintf("%d", int(totalEarnings))

		u.Logger.LogAny("CreateWithdraw.CreateWithDraw", zap.String("walletAddress", walletAddress),  zap.Any("widthdraw",f))
		err = u.Repo.CreateWithDraw(f)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw.CreateWithDraw", zap.Any("CreateWithDraw", f), zap.Error(err))
			return nil, err
		}

		if wr.PaymentType == string(entity.BIT){
			requestEarnings = requestEarnings / 1e8
		}else{
			requestEarnings = requestEarnings / 10e8
		}
		
		u.NotifyWithChannel(os.Getenv("SLACK_USER_CHANNEL"), fmt.Sprintf("[Withdraw has been created][User %s]", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress)), "", fmt.Sprintf("User %s making withdraw with %f %s ", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress), requestEarnings, wr.PaymentType))
		resp = append(resp, *f)
	}

	return resp, nil
}

func (u Usecase) FilterWidthdraw(data structure.FilterWithdraw) (*entity.Pagination, error) {
	f := &entity.FilterWithdraw{}

	err := copier.Copy(f, data)
	if err != nil {
		u.Logger.ErrorAny("FilterWidthdraw", zap.Any("data", data), zap.Error(err))
		return nil, err
	}

	p, err := u.Repo.FilterWithDraw(f)
	if err != nil {
		u.Logger.ErrorAny("FilterWidthdraw", zap.Any("FilterWithDraw", f), zap.Error(err))
		return nil, err
	}

	return p, nil
}

func (u Usecase) UpdateWithdraw(UUID string, status int) error {
	u.Logger.LogAny("UpdateWithdraw", zap.String("UUID",UUID), zap.Int("status",status) )
	err := u.Repo.UpdateWithDrawStatus(UUID, status)
	if err != nil {
		u.Logger.ErrorAny("UpdateWithdraw", zap.String("UUID",UUID), zap.Int("status",status) , zap.Error(err))
		return  err
	}

	return nil
}