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

func (u Usecase) CreateWithdraw(walletAddress string, wr structure.WithDrawItemRequest) (*entity.Withdraw, error) {

	u.Logger.LogAny("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("input", wr))
	volumeAmount := 0.0 //earning
	widthDrawAmount := 0.0
	refAmount := 0.0

	requestEarnings, err := strconv.ParseFloat(wr.Amount, 10)
	if err != nil {
		requestEarnings = 0
	}

	if requestEarnings < 0 {
		err = errors.New("Withdraw must be greater than Zero")
		return nil, err
	}

	//totalEarning := (refAmount + refAmount) - widthDrawAmount
	// (refAmount + refAmount) is pushed into volumn by crontab
	//TODO - calculate refAmount

	wdf := &entity.FilterWithdraw{
		WalletAddress:  &walletAddress,
		WithdrawItemID: &wr.ID,
		PaymentType:    &wr.PaymentType,
		Statuses: []int{
			entity.StatusWithdraw_Pending,
			entity.StatusWithdraw_Approve,
		},
	}

	wd, err := u.Repo.AggregateWithDrawByUser(wdf)
	if len(wd) > 0 {
		widthDrawAmount = wd[0].Amount
	}

	f := &entity.Withdraw{}
	err = copier.Copy(f, wr)
	if err != nil {
		u.Logger.ErrorAny("CreateWithdraw.Copy", zap.Any("input", wr), zap.Error(err))
		return nil, err
	}

	f.PayType = wr.PaymentType
	u.Logger.LogAny("CreateWithdraw.FilterVolume", zap.String("walletAddress", walletAddress))
	volumes, _ := u.GetVolumeOfProject(wr.ID, &f.PayType)

	// fr := entity.FilterReferrals{
	// 	ReferrerAddress: &walletAddress,
	// }

	// refs, _ := u.Repo.GetReferral(fr)
	// if len(refs) > 0 {
	// 	for _, ref := range refs {
	// 		tmp, err := strconv.ParseFloat(ref.ReferreeVolumn[wr.PaymentType].Amount, 10)
	// 		if err == nil {
	// 			refAmount +=   tmp
	// 		}
	// 	}
	// }

	u.Logger.LogAny("CreateWithdraw.volumes", zap.String("walletAddress", walletAddress), zap.Any("volumes", volumes))
	//check referal amount
	if volumes == nil {
		volumeAmount = 0
	} else {
		volumeAmount = volumes.Earning
	}

	//totalEarnings := refAmount + volumeAmount
	totalEarnings := volumeAmount
	availableBalance := totalEarnings - widthDrawAmount
	if availableBalance <= 0 {
		err = errors.New("Not enough balance")
		u.Logger.ErrorAny("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
		return nil, err
	}

	if requestEarnings > availableBalance {
		err = errors.New("RequestEarnings must be less than availableBalance")
		u.Logger.ErrorAny("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
		return nil, err
	}

	f.WalletAddress = walletAddress
	f.PayType = wr.PaymentType
	f.Amount = fmt.Sprintf("%d", int(requestEarnings))
	f.EarningReferal = fmt.Sprintf("%d", int(refAmount))
	f.EarningVolume = fmt.Sprintf("%d", int(volumeAmount))
	f.TotalEarnings = fmt.Sprintf("%d", int(totalEarnings))
	f.AvailableBalance = fmt.Sprintf("%d", int(availableBalance))
	f.WithdrawType = entity.Withdrawtype(wr.WithdrawType)
	f.WithdrawItemID = wr.ID
	f.Status = entity.StatusWithdraw_Pending

	// todo: 0x2525
	user, err := u.Repo.FindUserByBtcAddress(walletAddress)
	f.User = entity.WithdrawUserInfo{
		WalletAddress:        &user.WalletAddress,
		WalletAddressPayment: &user.WalletAddressPayment,
		WalletAddressBTC:     &user.WalletAddressBTC,
		DisplayName:          &user.DisplayName,
		Avatar:               &user.Avatar,
	}

	u.Logger.LogAny("CreateWithdraw.CreateWithDraw", zap.String("walletAddress", walletAddress), zap.Any("widthdraw", f))
	err = u.Repo.CreateWithDraw(f)
	if err != nil {
		u.Logger.ErrorAny("CreateWithdraw.CreateWithDraw", zap.Any("CreateWithDraw", f), zap.Error(err))
		return nil, err
	}

	if wr.PaymentType == string(entity.BIT) {
		requestEarnings = requestEarnings / 1e8
	} else {
		requestEarnings = requestEarnings / 1e8
	}

	u.UpdateRefObject(*f)

	u.NotifyWithChannel(os.Getenv("SLACK_WITHDRAW_CHANNEL"), fmt.Sprintf("[Pending withdraw has been created][User %s][ProjectID %s]", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress), helpers.CreateProjectLink(f.WithdrawItemID, f.WithdrawItemID)), "", fmt.Sprintf("User %s made withdraw with %f %s ", helpers.CreateProfileLink(f.WalletAddress, f.WalletAddress), requestEarnings, wr.PaymentType))
	return f, nil
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
	u.Logger.LogAny("UpdateWithdraw", zap.String("UUID", UUID), zap.Int("status", status))
	err := u.Repo.UpdateWithDrawStatus(UUID, status)
	if err != nil {
		u.Logger.ErrorAny("UpdateWithdraw", zap.String("UUID", UUID), zap.Int("status", status), zap.Error(err))
		return err
	}

	return nil
}

func (u Usecase) UpdateRefObject(withdraw entity.Withdraw) {
	u.Logger.LogAny("UpdateRefObject", zap.Any("withdraw", withdraw))
	switch withdraw.WithdrawType {
	case entity.WithDrawProject:
		p, err := u.Repo.FindProjectByTokenID(withdraw.WithdrawItemID)
		if err != nil {
			u.Logger.ErrorAny("UpdateRefObject.FindProjectByTokenID", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}

		u.Logger.ErrorAny("UpdateRefObject.FindProjectByTokenID", zap.Any("withdraw", withdraw), zap.Any("project", p))
		break
	case entity.WithDrawReferal:

		ref, err := u.Repo.GetAReferral(withdraw.WalletAddress, withdraw.WithdrawItemID)
		if err != nil {
			u.Logger.ErrorAny("UpdateRefObject.GetAReferral", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}

		u.Logger.ErrorAny("UpdateRefObject.GetAReferral", zap.Any("withdraw", withdraw), zap.Any("referral", ref))
		earning, _ := strconv.ParseFloat(ref.ReferreeVolumn[withdraw.PayType].Earn, 10)
		withDraw, _ := strconv.ParseFloat(withdraw.Amount, 10)
		refEarnings := earning - withDraw

		data := ref.ReferreeVolumn[withdraw.PayType]
		data.RemainingEarn = fmt.Sprintf("%d", int(refEarnings))

		ref.ReferreeVolumn[withdraw.PayType] = data
		updated, err := u.Repo.UpdateReferral(ref.UUID, ref)
		if err != nil {
			u.Logger.ErrorAny("UpdateRefObject.UpdateReferral", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}
		u.Logger.LogAny("UpdateRefObject.UpdateReferral", zap.Any("updated", updated))
		break
	}
}
