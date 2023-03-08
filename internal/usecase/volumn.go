package usecase

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) AggregateVolumns() {
	payTypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	// data, err := u.Repo.FindDataMissingRate()
	// if err != nil {
	// 	u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))

	// 	return
	// }

	// for _, item := range data {
	// 	item.BtcRate = 14.7
	// 	item.EthRate = 1

	// 	updated, err := u.Repo.UpdateMintNftBtc(&item)
	// 	if err != nil {
	// 		u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))

	// 		continue
	// 	}
	// 	u.Logger.LogAny("AggregationBTCWalletAddress", zap.Any("updated", updated))
	// }

	for _, payType := range payTypes {
		u.Logger.LogAny("AggregateVolumns", zap.Any("payType", payType))
		u.AggregateVolumn(payType)
	}
}

func (u Usecase) AggregateVolumn(payType string) {
	data, err := u.Repo.AggregateVolumn(payType)
	if err != nil {
		u.Logger.ErrorAny("CreateVolume", zap.Any("err", err))
		return
	}

	u.Logger.LogAny("AggregateVolumn", zap.Any("payType", payType), zap.Any("data", data))

	for _, item := range data {
		u.CreateVolumn(item)
	}
}

func (u Usecase) JobAggregateReferral() {

	referrals, err := u.Repo.GetAllReferrals(entity.FilterReferrals{})
	if err != nil {
		u.Logger.ErrorAny("JobAggregateReferral", zap.Any("err", err))
		return
	}

	paytypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	for _, referral := range referrals {
		vol := make(map[string]entity.ReferreeVolumn)
		for _, paytype := range paytypes {

			volume, err := u.GetVolumeOfUser(referral.Referree.WalletAddress, &paytype)
			if err != nil {
				vol[paytype] = entity.ReferreeVolumn{
					Amount:        "0",
					AmountType:    paytype,
					Earn:          "0",
					GenEarn:       "0",
					RemainingEarn: "0",
				}
			} else {
				refEarning, genEarning := helpers.CalculateRefEarning(volume.Amount, referral.Percent)
				remaining := referral.ReferreeVolumn[paytype].RemainingEarn
				if remaining == "" {
					remaining = "0"
				}

				vol[paytype] = entity.ReferreeVolumn{
					Amount:        fmt.Sprintf("%d", int(volume.Amount)),
					AmountType:    paytype,
					Earn:          refEarning,
					GenEarn:       genEarning,
					RemainingEarn: remaining,
				}
			}
		}
		referral.ReferreeVolumn = vol
		_, err = u.Repo.UpdateReferral(referral.UUID, &referral)
		if err != nil {
			u.Logger.ErrorAny("JobAggregateReferral", zap.Error(err))
			return
		}
	}
	_ = referrals
}

func (u Usecase) GetVolumeOfUser(walletAddress string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": bson.M{"creatorAddress": "$creatorAddress", "payType": "$payType"},
		"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
	}}

	amount, err := u.Repo.AggregateAmount(entity.FilterVolume{
		CreatorAddress: &walletAddress,
		AmountType:     amountType,
	}, group)
	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}

func (u Usecase) GetEarningOfUser(walletAddress string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": bson.M{"creatorAddress": "$creatorAddress", "payType": "$payType"},
		"amount": bson.M{"$sum": bson.M{"$toDouble": "$earning"}},
	}}

	amount, err := u.Repo.AggregateAmount(entity.FilterVolume{
		CreatorAddress: &walletAddress,
		AmountType:     amountType,
	}, group)
	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}

func (u Usecase) GetVolumeOfProject(projectID string, amountType *string) (*entity.AggregateAmount, error) {
	group := bson.M{"$group": bson.M{"_id": bson.M{"projectID": "$projectID", "payType": "$payType"},
		"amount":  bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
		"earning": bson.M{"$sum": bson.M{"$toDouble": "$earning"}},
	}}

	amount, err := u.Repo.AggregateAmount(entity.FilterVolume{
		ProjectID:  &projectID,
		AmountType: amountType,
	}, group)

	if err != nil {
		return nil, err
	}
	if len(amount) == 0 {
		return nil, errors.New("no document")
	}
	return &amount[0], nil
}

type csvLine struct {
	ProjectID  string
	Artist     string
	Collection string
	Status     string
	BTC        string
	ETH        string
}

func (u Usecase) MigrateFromCSV() {
	f, err := os.Open("artist_balance_1.csv")
	if err != nil {
		return
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return
	}

	csvData := []csvLine{}
	// convert records to array of structs
	for i, line := range data {
		if i > 1 { // omit header line
			tmp := csvLine{
				ProjectID:  line[0],
				Artist:     line[1],
				Collection: line[2],
				Status:     line[3],
				BTC:        line[4],
				ETH:        line[5],
			}

			csvData = append(csvData, tmp)
		}
	}
	spew.Dump(len(csvData))
	processCsvData := []csvLine{}
	for _, csv := range csvData {
		// if strings.ToLower(csv.Status) == "scam" {
		// 	continue
		// }
		// if csv.BTC == "0.00000" && csv.ETH == "0.00000" {
		// 	continue
		// }
		processCsvData = append(processCsvData, csv)
	}

	spew.Dump(len(processCsvData))
	wdsETH := []*entity.Withdraw{}
	for _, csv := range processCsvData {
		wd, _, err := u.CreateWD(csv, string(entity.ETH))
		if err != nil {
			continue
		}
		wdsETH = append(wdsETH, wd)
		// if isDuplicated {
		// 	wd1, _, _ := u.CreateWD(csv, string(entity.ETH))
		// 	wds = append(wds, wd1)
		// }
		u.Repo.CreateWithDraw(wd)

	}
	spew.Dump(len(wdsETH))

	wdsBTC := []*entity.Withdraw{}
	for _, csv := range processCsvData {
		wd, _, err := u.CreateWD(csv, string(entity.BIT))
		if err != nil {
			continue
		}

		wdsBTC = append(wdsBTC, wd)
		u.Repo.CreateWithDraw(wd)
		// if isDuplicated {
		// 	wd1, _, _ := u.CreateWD(csv, string(entity.ETH))
		// 	wds = append(wds, wd1)
		// }
	}
	spew.Dump(len(wdsBTC))

	// print the array
	//fmt.Printf("%+v\n", shoppingList)
}

func (u Usecase) CreateWD(csv csvLine, paymentType string) (*entity.Withdraw, bool, error) {
	p, err := u.Repo.FindProjectByTokenID(csv.ProjectID)
	dateString := "2023-02-28T04:05:26.385+00:00"
	date, _ := time.Parse("2023-02-28T00:00:00.000+00:00", dateString)
	if err != nil {
		u.Logger.ErrorAny("CreateWD.FindProjectByTokenID", zap.Error(err), zap.String("csv", csv.ProjectID), zap.String("paymentType", paymentType))
		return nil, false, err
	}
	isDuplicated := false

	wd := &entity.Withdraw{
		PayType:        paymentType,
		Status:         entity.StatusWithdraw_Approve,
		WalletAddress:  p.CreatorProfile.WalletAddress,
		WithdrawFrom:   "migrate_csv",
		Amount:         "0",
		EarningReferal: "0",
		EarningVolume:  "0",
		TotalEarnings:  "0",
		WithdrawType:   entity.WithDrawProject,
		WithdrawItemID: p.TokenID,
	}

	amount := ""
	if paymentType == string(entity.ETH) {
		eth := csv.ETH
		ethFloat, err := strconv.ParseFloat(eth, 10)
		if err != nil {
			u.Logger.ErrorAny("CreateWD.ParseFloat", zap.Error(err), zap.String("csv", csv.ProjectID), zap.String("paymentType", paymentType), zap.String("eth", eth))
			return nil, false, err
		}
		if ethFloat > 0 {
			return nil, false, errors.New("User was not paid")
		}
		if ethFloat == 0 {
			return nil, false, errors.New("Witdraw with zero")
		}
		ethFloat = ethFloat * -1 * 1e10
		amount = fmt.Sprintf("%d", int(ethFloat))

	} else {
		btc := csv.BTC
		btcFloat, err := strconv.ParseFloat(btc, 10)
		if err != nil {
			u.Logger.ErrorAny("CreateWD.ParseFloat", zap.Error(err), zap.String("csv", csv.ProjectID), zap.String("paymentType", paymentType), zap.String("btc", btc))
			return nil, false, err
		}
		if btcFloat > 0 {
			return nil, false, errors.New("User was not paid")
		}
		if btcFloat == 0 {
			return nil, false, errors.New("Witdraw with zero")
		}
		btcFloat = btcFloat * -1 * 1e8
		amount = fmt.Sprintf("%d", int(btcFloat))
	}

	usr := entity.WithdrawUserInfo{}
	user, err := u.Repo.FindUserByWalletAddress(p.CreatorAddrr)
	if err == nil {
		usr.WalletAddress = &user.WalletAddress
		usr.WalletAddressPayment = &user.WalletAddressPayment
		usr.WalletAddressBTC = &user.WalletAddressBTC
		usr.DisplayName = &user.DisplayName
		usr.Avatar = &user.Avatar
	}
	wd.Amount = amount
	wd.EarningVolume = amount
	wd.TotalEarnings = amount
	wd.CreatedAt = &date
	wd.Note = "Add the paid artist on Feb 2023"
	u.Logger.LogAny("CreateWD.wd", zap.String("paymentType", paymentType), zap.Any("wd", wd))
	wd.User = usr

	return wd, isDuplicated, nil
}

func (u Usecase) CreateVolumn(item entity.AggregateProjectItemResp) {

	u.Logger.LogAny("aggregateVolumn", zap.Any("item", item))
	pID := strings.ToLower(item.ProjectID)
	p, err := u.Repo.FindProjectByTokenID(pID)
	if err != nil {
		u.Logger.ErrorAny("FindProjectByTokenID", zap.String("item.ProjectID", item.ProjectID), zap.Any("err", err))
		return
	}

	creatorID := strings.ToLower(p.CreatorAddrr)
	usr, err := u.Repo.FindUserByWalletAddress(creatorID)
	if err != nil {
		u.Logger.ErrorAny("FindUserByWalletAddress", zap.String("p.CreatorAddrr", creatorID), zap.Any("err", err))
		return
	}

	mintPrice := 0.0
	if item.Paytype == string(entity.BIT) {
		ar, err := u.Repo.AggregateProjectMintPrice(item.ProjectID, item.Paytype)
		if err == nil && len(ar) > 0 {
			mintPrice = ar[0].Amount
		} else {
			pFl, _ := strconv.ParseFloat(p.MintPrice, 10)
			mintPrice = pFl
		}

		oldData, err := u.AggregateOldBtcAddress(item.ProjectID)
		if err == nil {
			spew.Dump(item, oldData)
			item.Amount += oldData.Amount
			item.Minted += oldData.Minted
		}
	} else {
		ar, err := u.Repo.AggregateProjectMintPrice(item.ProjectID, item.Paytype)
		if err == nil && len(ar) > 0 {
			mintPrice = ar[0].Amount
		} else {
			pFl, _ := strconv.ParseFloat(p.MintPrice, 10)
			mintPrice = pFl
		}

		oldData, err := u.AggregateOldETHAddress(item.ProjectID)
		if err == nil {
			spew.Dump(item, oldData)
			item.Amount += oldData.Amount
			item.Minted += oldData.Minted
		}
	}

	ev, err := u.Repo.FindVolumn(pID, item.Paytype)
	if err != nil {
		amount := fmt.Sprintf("%d", int(item.Amount))
		earning, gearning := helpers.CalculateVolumEarning(item.Amount, int32(utils.PERCENT_EARNING))
		if errors.Is(err, mongo.ErrNoDocuments) {
			//v := entity.FilterVolume
			ev := &entity.UserVolumn{
				CreatorAddress: &creatorID,
				PayType:        &item.Paytype,
				ProjectID:      &pID,
				Amount:         &amount,
				Earning:        &earning,
				GenEarning:     &gearning,
				Minted:         item.Minted,
				MintPrice:      int64(mintPrice),
				Project: entity.VolumeProjectInfo{
					Name:     p.Name,
					TokenID:  p.TokenID,
					Thumnail: p.Thumbnail,
				},
				User: entity.VolumnUserInfo{
					WalletAddress:    &p.CreatorAddrr,
					WalletAddressBTC: &usr.WalletAddressBTC,
					DisplayName:      &usr.DisplayName,
					Avatar:           &usr.Avatar,
				},
			}

			err = u.Repo.CreateVolumn(ev)
			if err != nil {
				u.Logger.ErrorAny("CreateVolumn", zap.Any("ev", ev), zap.Any("err", err))
				return
			}
		}
	} else {
		amount := fmt.Sprintf("%d", int(item.Amount))
		if amount != *ev.Amount {
			earning, gearning := helpers.CalculateVolumEarning(item.Amount, int32(utils.PERCENT_EARNING))
			_, err := u.Repo.UpdateVolumnAmount(*ev.ProjectID, *ev.PayType, amount, earning, gearning)
			if err != nil {
				u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
				return
			}
		}

		_, err := u.Repo.UpdateVolumnMinted(*ev.ProjectID, *ev.PayType, item.Minted)
		if err != nil {
			u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
			return
		}

		if int(mintPrice) != int(ev.MintPrice) {
			_, err := u.Repo.UpdateVolumMintPrice(*ev.ProjectID, *ev.PayType, item.MintPrice)
			if err != nil {
				u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
				return
			}
		}
	}
}

func (u Usecase) AggregateOldBtcAddress(projectID string) (*entity.AggregateProjectItemResp, error) {
	u.Logger.LogAny("AggregationBTCWalletAddress", zap.Any("projectID", projectID))
	data, err := u.Repo.AggregationBTCWalletAddress(projectID)
	if err != nil {
		u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))
	}

	u.Logger.LogAny("AggregationBTCWalletAddress", zap.Any("data", data))
	if len(data) > 0 {
		item := data[0]
		item.Paytype = string(entity.BIT)
		item.BtcRate = 14.7
		item.EthRate = 1
		return &item, nil
	}
	return nil, errors.New("no olf data")

}

func (u Usecase) AggregateOldETHAddress(projectID string) (*entity.AggregateProjectItemResp, error) {
	u.Logger.LogAny("AggregateOldETHAddress", zap.Any("projectID", projectID))
	dataETH, err := u.Repo.AggregationETHWalletAddress(projectID)
	if err != nil {
		u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))
	}

	u.Logger.LogAny("AggregateOldETHAddress", zap.Any("dataETH", dataETH))
	if len(dataETH) > 0 {
		item := dataETH[0]
		item.MintPrice = item.MintPrice / 1e10
		item.Amount = item.Amount / 1e10
		item.Paytype = string(entity.ETH)
		item.BtcRate = 14.7
		item.EthRate = 1
		return &item, nil
	}

	return nil, errors.New("no old data")
}
