package usecase

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateWithdraw(walletAddress string, wr structure.WithDrawItemRequest) (*entity.Withdraw, error) {
	logger.AtLog.Logger.Info("CreateWithdraw", zap.String("walletAddress", walletAddress), zap.Any("input", zap.Any("wr)", wr)))
	if wr.WithdrawType == string(entity.WithDrawProject) {
		return u.CreateWithdrawProject(walletAddress, wr)
	}

	if wr.WithdrawType == string(entity.WithDrawReferal) {
		return u.CreateWithdrawReferral(walletAddress, wr)
	}

	return nil, errors.New("Withdraw type is not allowed")
}

func (u Usecase) CreateWithdrawProject(walletAddress string, wr structure.WithDrawItemRequest) (*entity.Withdraw, error) {

	logger.AtLog.Logger.Info("CreateWithdrawProject", zap.String("walletAddress", walletAddress), zap.Any("input", zap.Any("wr)", wr)))
	volumeAmount := 0.0 //earning
	widthDrawAmount := 0.0
	refAmount := 0.0

	project, err := u.Repo.FindProjectByTokenIDOrGenNFTAddr(wr.ID)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(project.CreatorAddrr) != strings.ToLower(walletAddress) {
		err := errors.New(fmt.Sprintf("Yout don't have permission to make withdraw to this collection"))
		return nil, err
	}

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

	wdType := string(entity.WithDrawProject)
	wdf := &entity.FilterWithdraw{
		WalletAddress:  &walletAddress,
		WithdrawItemID: &wr.ID,
		PaymentType:    &wr.PaymentType,
		WithdrawType:   &wdType,
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
		logger.AtLog.Logger.Error("CreateWithdraw.Copy", zap.Any("input", wr), zap.Error(err))
		return nil, err
	}

	f.PayType = wr.PaymentType
	logger.AtLog.Logger.Info("CreateWithdraw.FilterVolume", zap.String("walletAddress", walletAddress))
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

	logger.AtLog.Logger.Info("CreateWithdraw.volumes", zap.String("walletAddress", walletAddress), zap.Any("volumes", zap.Any("volumes)", volumes)))
	//check referal amount
	if volumes == nil {
		volumeAmount = 0
	} else {
		volumeAmount = volumes.Earning
	}

	//totalEarnings := refAmount + volumeAmount
	totalEarnings := volumeAmount
	availableBalance := totalEarnings - widthDrawAmount
	if availableBalance < 0 {
		err = errors.New("Not enough balance")
		logger.AtLog.Logger.Error("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
		return nil, err
	}

	if requestEarnings > availableBalance {
		err = errors.New("RequestEarnings must be less than availableBalance")
		logger.AtLog.Logger.Error("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
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

	user, err := u.Repo.FindUserByWalletAddress(walletAddress)
	if err != nil {
		return nil, err
	}
	if user.WalletAddressPayment == "" && f.PayType == "eth" {
		return nil, fmt.Errorf("eth payment address not found")
	}
	f.User = entity.WithdrawUserInfo{
		WalletAddress:        &user.WalletAddress,
		WalletAddressPayment: &user.WalletAddressPayment,
		WalletAddressBTC:     &user.WalletAddressBTC,
		DisplayName:          &user.DisplayName,
		Avatar:               &user.Avatar,
	}

	logger.AtLog.Logger.Info("CreateWithdraw.CreateWithDraw", zap.String("walletAddress", walletAddress), zap.Any("widthdraw", zap.Any("f)", f)))
	err = u.Repo.CreateWithDraw(f)
	if err != nil {
		logger.AtLog.Logger.Error("CreateWithdraw.CreateWithDraw", zap.Any("CreateWithDraw", f), zap.Error(err))
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

func (u Usecase) CreateWithdrawReferral(walletAddress string, wr structure.WithDrawItemRequest) (*entity.Withdraw, error) {

	logger.AtLog.Logger.Info("CreateWithdrawReferral", zap.String("walletAddress", walletAddress), zap.Any("input", zap.Any("wr)", wr)))
	volumeAmount := 0.0 //earning
	widthDrawAmount := 0.0
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

	wdType := string(entity.WithDrawReferal)
	wdf := &entity.FilterWithdraw{
		WalletAddress:  &walletAddress,
		WithdrawItemID: &wr.ID,
		WithdrawType:   &wdType,
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
		logger.AtLog.Logger.Error("CreateWithdraw.Copy", zap.Any("input", wr), zap.Error(err))
		return nil, err
	}

	f.PayType = wr.PaymentType
	logger.AtLog.Logger.Info("CreateWithdraw.FilterVolume", zap.String("walletAddress", walletAddress))
	refEarning, _ := u.Repo.GetAReferral(walletAddress, wr.ID)

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

	//logger.AtLog.Logger.Info("CreateWithdraw.volumes", zap.String("walletAddress", walletAddress), zap.Any("volumes", zap.Any("volumes)", volumes)))
	//check referal amount
	if refEarning == nil {
		volumeAmount = 0
	} else {
		volumes := refEarning.ReferreeVolumn[f.PayType].Earn
		volumeF, err := strconv.ParseFloat(volumes, 10)
		if err == nil {
			//volumeAmount = volumes.A
			volumeAmount = volumeF
		}

	}

	//totalEarnings := refAmount + volumeAmount
	totalEarnings := volumeAmount
	availableBalance := totalEarnings - widthDrawAmount
	if availableBalance < 0 {
		err = errors.New("Not enough balance")
		logger.AtLog.Logger.Error("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
		return nil, err
	}

	if requestEarnings > availableBalance {
		err = errors.New("RequestEarnings must be less than availableBalance")
		logger.AtLog.Logger.Error("CreateWithdraw", zap.Float64("earning", availableBalance), zap.String("walletAddress", walletAddress), zap.Any("volumeAmount", volumeAmount), zap.Error(err))
		return nil, err
	}

	f.WalletAddress = walletAddress
	f.PayType = wr.PaymentType
	f.Amount = fmt.Sprintf("%d", int(requestEarnings))
	f.EarningReferal = fmt.Sprintf("%d", int(volumeAmount))
	f.EarningVolume = fmt.Sprintf("%d", 0)
	f.TotalEarnings = fmt.Sprintf("%d", int(totalEarnings))
	f.AvailableBalance = fmt.Sprintf("%d", int(availableBalance))
	f.WithdrawType = entity.Withdrawtype(wr.WithdrawType)
	f.WithdrawItemID = wr.ID
	f.Status = entity.StatusWithdraw_Pending

	user, err := u.Repo.FindUserByWalletAddress(walletAddress)
	if err != nil {
		return nil, err
	}
	if user.WalletAddressPayment == "" && f.PayType == "eth" {
		return nil, fmt.Errorf("eth payment address not found")
	}

	f.User = entity.WithdrawUserInfo{
		WalletAddress:        &user.WalletAddress,
		WalletAddressPayment: &user.WalletAddressPayment,
		WalletAddressBTC:     &user.WalletAddressBTC,
		DisplayName:          &user.DisplayName,
		Avatar:               &user.Avatar,
	}

	logger.AtLog.Logger.Info("CreateWithdraw.CreateWithDraw", zap.String("walletAddress", walletAddress), zap.Any("widthdraw", zap.Any("f)", f)))
	err = u.Repo.CreateWithDraw(f)
	if err != nil {
		logger.AtLog.Logger.Error("CreateWithdraw.CreateWithDraw", zap.Any("CreateWithDraw", f), zap.Error(err))
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
		logger.AtLog.Logger.Error("FilterWidthdraw", zap.Any("data", data), zap.Error(err))
		return nil, err
	}

	p, err := u.Repo.FilterWithDraw(f)
	if err != nil {
		logger.AtLog.Logger.Error("FilterWidthdraw", zap.Any("FilterWithDraw", f), zap.Error(err))
		return nil, err
	}

	return p, nil
}

func (u Usecase) UpdateWithdraw(UUID string, status int) error {
	logger.AtLog.Logger.Info("UpdateWithdraw", zap.String("UUID", UUID), zap.Int("status", status))
	err := u.Repo.UpdateWithDrawStatus(UUID, status)
	if err != nil {
		logger.AtLog.Logger.Error("UpdateWithdraw", zap.String("UUID", UUID), zap.Int("status", status), zap.Error(err))
		return err
	}

	return nil
}

func (u Usecase) UpdateRefObject(withdraw entity.Withdraw) {
	logger.AtLog.Logger.Info("UpdateRefObject", zap.Any("withdraw", zap.Any("withdraw)", withdraw)))
	switch withdraw.WithdrawType {
	case entity.WithDrawProject:
		p, err := u.Repo.FindProjectByTokenID(withdraw.WithdrawItemID)
		if err != nil {
			logger.AtLog.Logger.Error("UpdateRefObject.FindProjectByTokenID", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}

		logger.AtLog.Logger.Error("UpdateRefObject.FindProjectByTokenID", zap.Any("withdraw", withdraw), zap.Any("project", p))
		break
	case entity.WithDrawReferal:

		ref, err := u.Repo.GetAReferral(withdraw.WalletAddress, withdraw.WithdrawItemID)
		if err != nil {
			logger.AtLog.Logger.Error("UpdateRefObject.GetAReferral", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}

		logger.AtLog.Logger.Error("UpdateRefObject.GetAReferral", zap.Any("withdraw", withdraw), zap.Any("referral", ref))
		earning, _ := strconv.ParseFloat(ref.ReferreeVolumn[withdraw.PayType].Earn, 10)
		withDraw, _ := strconv.ParseFloat(withdraw.Amount, 10)
		remaining := earning - withDraw

		data := ref.ReferreeVolumn[withdraw.PayType]
		data.RemainingEarn = fmt.Sprintf("%d", int(remaining))

		ref.ReferreeVolumn[withdraw.PayType] = data
		updated, err := u.Repo.UpdateReferral(ref.UUID, ref)
		if err != nil {
			logger.AtLog.Logger.Error("UpdateRefObject.UpdateReferral", zap.Any("withdraw", withdraw), zap.Error(err))
			return
		}
		logger.AtLog.Logger.Info("UpdateRefObject.UpdateReferral", zap.Any("updated", zap.Any("updated)", updated)))
		break
	}
}
