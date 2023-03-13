package usecase

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u Usecase) JobAggregateVolumns() {
	cached := []string{}
	str, err := u.Cache.GetData("btcIDs")
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(*str), &cached)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	projects, err := u.Repo.FindProjectsHaveMinted(cached)
	if err != nil {
		return
	}
	fileName := fmt.Sprintf("projects-volumn-%s.json", now)
	helpers.CreateFile(fileName, projects)
	

	payTypes := []string{
		string(entity.BIT),
		string(entity.ETH),
	}

	pLogs := []structure.VolumnLogs{}
	pLogsChannel := make(chan structure.VolumnLogs , len(projects) * 2)
	for _, project := range projects {
		for _, paytype := range payTypes {
			go func(project entity.ProjectsHaveMinted, paytype string, pLogsChannel chan structure.VolumnLogs) {
				logger.AtLog.Logger.Info("Calculating ...", zap.Any("project",project),  zap.Any("paytype",paytype))
				minted := 0
				amount := 0.0
				data, err := u.Repo.AggregateVolumn(project.TokenID, paytype)
				if err == nil && data != nil {
					if len(data) > 0 {
						minted = data[0].Minted
						amount = data[0].Amount
					}
				}

				oldMinted := 0
				oldAmount := 0.0
				oldData, err := u.AggregateOldData(project.TokenID, paytype)
				if err == nil && oldData != nil {
					oldMinted = oldData.Minted
					oldAmount = oldData.Amount
				}

				totalAmout := amount + oldAmount
				earning, gearning := helpers.CalculateVolumEarning(totalAmout, int32(utils.PERCENT_EARNING))
				earningF, _ := strconv.ParseFloat(earning, 10)

				wd := 0.0
				widthDraw, err := u.Repo.AggregateWithDrawByProject(project.TokenID, paytype)
				if err == nil && len(widthDraw) >0 {
					for _, wdItem := range widthDraw {
						wd += wdItem.Amount
					}
				}

				pLog := structure.VolumnLogs{
					ProjectID:     project.TokenID,
					Paytype:       paytype,
					OldMinted:     oldMinted,
					NewMinted:     minted,
					TotalMinted:   oldMinted + minted,
					OldAmount:     fmt.Sprintf("%d", int(oldAmount)),
					NewAmount:     fmt.Sprintf("%d", int(amount)),
					TotalAmount:   fmt.Sprintf("%d", int(totalAmout)),
					TotalEarnings: earning,
					ApprovedWithdraw: fmt.Sprintf("%d", int(wd)),
					Available: fmt.Sprintf("%d", int(earningF - wd)),
					GenEarnings:   gearning,
					SeparateRate:  fmt.Sprintf("%d", utils.PERCENT_EARNING),
					MintPrice: u.AggregateMintPrice(project, paytype),
				}
				
				
				pLogsChannel <- pLog
				

			}(project, paytype, pLogsChannel)
		}
	}

	for _, _ = range projects {
		for _, _ = range payTypes {
			pLog := <-pLogsChannel
			//u.CreateVolumn(pLog)
			a, _ :=  strconv.Atoi(pLog.Available)
			if pLog.OldMinted > 0 &&  a > 0 {
				pLogs = append(pLogs, pLog)
			}
			
		}
	}

	fileName = fmt.Sprintf("aggregated-volumn-%s.json", now)
	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName = strings.ReplaceAll(fileName, ":", "_")
	fileName = strings.ReplaceAll(fileName, "+", "_")
	fileName = strings.ToLower(fileName)

	helpers.CreateFile(fileName, pLogs)
	//bytes, err := json.Marshal(pLogs)
	if err == nil {

		//base64String := helpers.Base64Encode(bytes)
		//uploaded, err := u.GCS.UploadBaseToBucket(base64String, fileName)
		if err == nil {
			//spew.Dump(uploaded)
			//u.NotifyWithChannel(os.Getenv("SLACK_WITHDRAW_CHANNEL"), "[Volumns have been created]", "Please refer to the following URL", helpers.CreateURLLink(fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name), uploaded.Name))
		}
	}

	u.Cache.SetData("pLogs", pLogs)
	spew.Dump("done")
}

// func (u Usecase) AggregateVolumn(payType string) {
// 	data, err := u.Repo.AggregateVolumn(payType)
// 	if err != nil {
// 		return
// 	}

// 	now := time.Now().UTC()
// 	helpers.CreateFile(fmt.Sprintf("aggregateVolumn-%s-%s.json",payType, now), data)

// 	// for _, item := range data {
// 	// 	u.CreateVolumn(item)
// 	// }
// }

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
	f, err := os.Open("13-03-2023_paid_artists.csv")
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
	csvDataWithAmount := make(map[string]csvLine)
	projectIDs := []string{}
	//convert records to array of structs
	for i, line := range data {
		if i > 1 { // omit header line
			tmp := csvLine{
				ProjectID:   line[0],
				// Artist:      line[1],
				// Collection:  line[2],
				BTC:       	 line[5],
				ETH:          line[6],
			}
			projectIDs = append(projectIDs, line[0])
			csvData = append(csvData, tmp)

			csvDataWithAmount[tmp.ProjectID] = tmp
		}
	}

	type mintedData struct {
		ProjectID string
		Amount string
		Withdraw string
		PayType string
		Minted int
		Available string
	}

	helpers.CreateFile("projectIDs.json",projectIDs)
	dataResp := []mintedData{}
	btcWalletAddresses , err := u.Repo.FindWalletAddressesIn(projectIDs, entity.BTCWalletAddress{}.TableName())
	if err != nil {
		logger.AtLog.Error(zap.Error(err))
	}


	btcIDs := []string{}
	for _, btcWalletAddress := range btcWalletAddresses {
		amount :=  btcWalletAddress.Amount
		withdraw := csvDataWithAmount[btcWalletAddress.ProjectID].BTC
		withdrawF, _ := strconv.ParseFloat(withdraw, 10)  
		withdrawF = withdrawF * 1e8
		a := amount*0.9 - withdrawF

		dataResp  = append(dataResp, mintedData{
			ProjectID: btcWalletAddress.ProjectID,
			Amount: fmt.Sprintf("%d", int(amount)),
			Withdraw: fmt.Sprintf("%d", int(withdrawF)),
			Minted: btcWalletAddress.Minted,
			PayType: string(entity.BIT),
			Available: fmt.Sprintf("%d", int(a)) ,
		}) 

		btcIDs = append(btcIDs, btcWalletAddress.ProjectID)
	}
	
	helpers.CreateFile("dataResp-btc.json",dataResp)
	ethWalletAddresses , err := u.Repo.FindWalletAddressesIn(projectIDs, entity.ETHWalletAddress{}.TableName())
	if err != nil {
		logger.AtLog.Error(zap.Error(err))
	}

	//ethIDs := []string{}
	dataResp = []mintedData{}
	for _, ethWalletAddress := range ethWalletAddresses {
		amount :=  ethWalletAddress.Amount / 1e10
		withdraw := csvDataWithAmount[ethWalletAddress.ProjectID].ETH
		withdrawF, _ := strconv.ParseFloat(withdraw, 10)  
		withdrawF = withdrawF * 1e8
		a := amount*0.9 - withdrawF

		//spew.Dump(ethWalletAddress.ProjectID, ethWalletAddress.Minted, amount)
		dataResp  = append(dataResp, mintedData{
			ProjectID: ethWalletAddress.ProjectID,
			Amount: fmt.Sprintf("%d", int(amount)),
			Withdraw: fmt.Sprintf("%d", int(withdrawF)),
			Minted: ethWalletAddress.Minted,
			PayType: string(entity.ETH),
			Available: fmt.Sprintf("%d", int(a)) ,
		}) 

		btcIDs = append(btcIDs, ethWalletAddress.ProjectID)
	}

	u.Cache.SetData("btcIDs",btcIDs)
	helpers.CreateFile("dataResp-eth.json",dataResp)
	
	//spew.Dump(len(csvData))
	//processCsvData := []csvLine{}
	// for _, csv := range csvData {
	// 	// if strings.ToLower(csv.Status) == "scam" {
	// 	// 	continue
	// 	// }
	// 	// if csv.BTC == "0.00000" && csv.ETH == "0.00000" {
	// 	// 	continue
	// 	// }
	// 	processCsvData = append(processCsvData, csv)
	// }

	// tmp := csvLine{
	// 	ProjectID:  "1001311",
	// 	Artist:     "crashblossom",
	// 	Collection: "RECALL",
	// 	//Status:     line[3],
	// 	BTC:        "0.045",
	// 	ETH:        "43.40304",
	// }

	// //csvData = append(csvData, tmp)
	// processCsvData := []csvLine{}
	// processCsvData = append(processCsvData, tmp)
	// spew.Dump(len(processCsvData))
	// wdsETH := []*entity.Withdraw{}
	// for _, csv := range processCsvData {
	// 	wd, _, err := u.CreateWD(csv, string(entity.ETH))
	// 	if err != nil {
	// 		continue
	// 	}
	// 	wdsETH = append(wdsETH, wd)
	// 	// if isDuplicated {
	// 	// 	wd1, _, _ := u.CreateWD(csv, string(entity.ETH))
	// 	// 	wds = append(wds, wd1)
	// 	// }
	// 	u.Repo.CreateWithDraw(wd)

	// }
	// spew.Dump(len(wdsETH))

	// wdsBTC := []*entity.Withdraw{}
	// for _, csv := range processCsvData {
	// 	wd, _, err := u.CreateWD(csv, string(entity.BIT))
	// 	if err != nil {
	// 		continue
	// 	}

	// 	wdsBTC = append(wdsBTC, wd)
	// 	u.Repo.CreateWithDraw(wd)
	// 	// if isDuplicated {
	// 	// 	wd1, _, _ := u.CreateWD(csv, string(entity.ETH))
	// 	// 	wds = append(wds, wd1)
	// 	// }
	// }
	// spew.Dump(len(wdsBTC))

	// print the array
	//fmt.Printf("%+v\n", shoppingList)
}

func (u Usecase) CreateWDs() {
	data, err := u.Cache.GetData("pLogs")
	if err != nil {
		return
	}

	pLogs := []structure.VolumnLogs{}
	err = json.Unmarshal([]byte(*data), &pLogs)
	if err != nil {
		return
	}

	wdLog := []entity.Withdraw{}
	for _, logs := range pLogs {
		wd, _, err := u.CreateWD(logs, logs.Paytype)
		if err != nil {
			return
		}

		err = u.Repo.CreateWithDraw(wd)
		if err != nil {
			return
		}
		wdLog = append(wdLog, *wd)
	}	
	helpers.CreateFile("withdraw-log.json",wdLog)

}

func (u Usecase) CreateWD(csv structure.VolumnLogs, paymentType string) (*entity.Withdraw, bool, error) {
	p, err := u.Repo.FindProjectByTokenID(csv.ProjectID)
	dateString := "2023-03-10T04:05:26.385+00:00"
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
		WithdrawFrom:   "fix_zero_balance",
		Amount:         "0",
		EarningReferal: "0",
		EarningVolume:  "0",
		TotalEarnings:  "0",
		WithdrawType:   entity.WithDrawProject,
		WithdrawItemID: p.TokenID,
	}

	amount := csv.Available
	arrge, err := u.Repo.FindVolumnByWalletAddress(p.CreatorProfile.WalletAddress, paymentType)
	if err == nil &&  arrge !=nil {
		wd.EarningVolume = amount
		wd.TotalEarnings = amount
		wdf, _ := strconv.ParseFloat(amount, 10)
		
		earningF, _  := strconv.ParseFloat(amount, 10)
		wd.AvailableBalance = fmt.Sprintf("%d", int(earningF - wdf))
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
	wd.PayType = csv.Paytype
	wd.CreatedAt = &date
	wd.Note = "Add the paid artist on Mar 2023"
	u.Logger.LogAny("CreateWD.wd", zap.String("paymentType", paymentType), zap.Any("wd", wd))
	wd.User = usr

	return wd, isDuplicated, nil
}

func (u Usecase) CreateVolumn(item structure.VolumnLogs) {
	logger.AtLog.Logger.Info("CreateVolumn...", zap.Any("item",item))
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

	ev, err := u.Repo.FindVolumn(pID, item.Paytype)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//v := entity.FilterVolume
			ev := &entity.UserVolumn{
				CreatorAddress: &creatorID,
				PayType:        &item.Paytype,
				ProjectID:      &pID,
				Amount:         &item.TotalAmount,
				Earning:        &item.TotalEarnings,
				GenEarning:     &item.GenEarnings,
				Minted:         item.TotalMinted,
				MintPrice:      int64(item.MintPrice),
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

		if item.TotalAmount != *ev.Amount {
			_, err := u.Repo.UpdateVolumnAmount(*ev.ProjectID, *ev.PayType, item.TotalAmount, item.TotalEarnings, item.GenEarnings)
			if err != nil {
				u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
				return
			}
		}
		
		_, err := u.Repo.UpdateVolumnMinted(*ev.ProjectID, *ev.PayType, item.TotalMinted)
		if err != nil {
			u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
			return
		}

		if item.MintPrice != int(ev.MintPrice) {
			_, err := u.Repo.UpdateVolumMintPrice(*ev.ProjectID, *ev.PayType, int64(item.MintPrice))
			if err != nil {
				u.Logger.ErrorAny("UpdateVolumnAmount", zap.String("p.CreatorAddrr", p.CreatorAddrr), zap.Any("err", err))
				return
			}
		}
	}
}

func (u Usecase) AggregateOldBtcAddress(projectID string) (*entity.AggregateProjectItemResp, error) {

	data, err := u.Repo.AggregationBTCWalletAddress(projectID)
	if err != nil {
		u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))
	}

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
	dataETH, err := u.Repo.AggregationETHWalletAddress(projectID)
	if err != nil {
		u.Logger.ErrorAny("AggregationBTCWalletAddress", zap.Error(err))
	}

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

func (u Usecase) AggregateOldData(projectID string, payType string) (*entity.AggregateProjectItemResp, error) {
	if payType == string(entity.ETH) {
		return u.AggregateOldETHAddress(projectID)
	}
	return u.AggregateOldBtcAddress(projectID)
}

func (u Usecase) AggregateMintPrice(project entity.ProjectsHaveMinted, payType string) int {
	mintPrice := 0.0
	ar, err := u.Repo.AggregateProjectMintPrice(project.TokenID, payType)
	if err == nil && len(ar) > 0 {
		mintPrice = ar[0].Amount
	} else {
		if payType == string(entity.ETH) {
			if project.MintPriceEth == "" {
				pFl, _ := strconv.ParseFloat(project.MintPrice, 10)
				mintPrice = pFl * 14.7 //hard code for the old projects
			}

			pFl, _ := strconv.ParseFloat(project.MintPriceEth, 10)
			mintPrice = pFl / 1e10
		}else{
			pFl, _ := strconv.ParseFloat(project.MintPrice, 10)
			mintPrice = pFl
		}
	}
	return int(mintPrice)
}


func (u Usecase) WriteCsv(fileName string, records []csvLine) {
	csvFile, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		csvFile.Close()
	
		csvwriter := csv.NewWriter(csvFile)
		for _, record := range records {
			row := []string{record.ProjectID}
			_ = csvwriter.Write(row)
		}
		csvwriter.Flush()
	
	
	
		return 
}