package usecase

import (
	"fmt"
	"os"

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
		//totalEarning := (refAmount + refAmount) - widthDrawAmount
		// (refAmount + refAmount) is pushed into volumn by crontab
		//TODO - calculate refAmount
	
		f := &entity.Withdraw{}
		err := copier.Copy(f, wr)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw", zap.Any("input", data), zap.Error(err))
			return nil, err
		}

		//check amount from User volume (projects earnings)
		fv := entity.FilterVolume{
			//ProjectIDs: []string{wr.ProjectID},
			AmountType: &wr.PaymentType,
			CreatorAddress: &walletAddress,
		}

		u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv))
		volumes, err := u.Repo.AggregateAmount(fv)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw", zap.Any("input", data), zap.Error(err))
			return nil, err
		}

		u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv), zap.Any("volumes", volumes))
		//check widthdraw history
		stat :=  entity.StatusWithdraw_Done
		fWdtd := &entity.FilterWithdraw{
			WalletAddress: &walletAddress,
			Status: &stat,
			PaymentType: &wr.PaymentType,
		}
		u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv), zap.Any("volumes", volumes), zap.Any("fWdtd", fWdtd))
		wtd, err := u.Repo.AggregateWithDraw(fWdtd)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv), zap.Any("volumeAmount", volumeAmount), zap.Any("fWdtd", fWdtd), zap.Error(err))
			return nil, err
		}

		//check referal amount
		if len(volumes) == 0 {
			volumeAmount = 0
		}else{
			volumeAmount = volumes[0].Amount
		}
		
		if len(wtd) == 0 {
			widthDrawAmount = 0
		}else{
			widthDrawAmount = wtd[0].Amount
		}

		
		earning := volumeAmount - widthDrawAmount
		//TODO - validate these things
		// if earning <= 0 {
		// 	err = errors.New("Not enough balance")
		// 	u.Logger.ErrorAny("CreateWithdraw", zap.Float64("earning", earning) , zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv), zap.Any("volumeAmount", volumeAmount), zap.Any("fWdtd", fWdtd), zap.Error(err))
		// 	return nil, err
		// }


		f.WalletAddress = walletAddress
		f.PayType = wr.PaymentType
		f.Amount = fmt.Sprintf("%d", int(earning))

		u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("filterVolume", fv), zap.Any("widthdraw",f))
		err = u.Repo.CreateWithDraw(f)
		if err != nil {
			u.Logger.ErrorAny("CreateWithdraw", zap.Any("CreateWithDraw", f), zap.Error(err))
			return nil, err
		}

		u.NotifyWithChannel(os.Getenv("SLACK_USER_CHANNEL"), fmt.Sprintf("[Withdraw has been created][User %s]", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress)), "", fmt.Sprintf("User %s making withdraw with %s %s ", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress), f.Amount, wr.PaymentType))
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