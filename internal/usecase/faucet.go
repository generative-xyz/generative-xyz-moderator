package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiCreateFaucet(url string) (string, error) {

	// verify tw name:
	// //https://twitter.com/2712_at1999/status/1643190049981480961
	twNameRegex := regexp.MustCompile(`https://twitter.com/(\w+)/status/(\d+)(?:\?s=\d+)?`)
	// Find the first match in the tweet URL
	matchTwName := twNameRegex.FindStringSubmatch(url)

	twName := ""
	sharedID := ""

	if len(matchTwName) >= 3 {

		twName = matchTwName[1]
		shareID := matchTwName[2]
		fmt.Println("twName:", twName)   // Output: 2712_at1999
		fmt.Println("shareID:", shareID) // Output: 1643190049981480961

	} else {
		err := errors.New("The username or sharedID is not found in the tweet URL.")
		go u.sendSlack("", "ApiCreateFaucet.regexp", "check url"+url, err.Error())
		return "", err
	}
	// check valid vs twName first:
	err := u.CheckValidFaucet("no", twName)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.checkValidFaucet"), zap.Error(err))
		go u.sendSlack("", "ApiCreateFaucet.CheckValidFaucet.twName", twName, err.Error())
		return "", err
	}

	// check sharedID exist:
	sharedIDs, _ := u.Repo.FindFaucetBySharedID(sharedID)
	if len(sharedIDs) > 0 {
		err := errors.New("The shard ID already exists, please tweet a new one.")
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetBySharedID"), zap.Error(err))
		go u.sendSlack("", "ApiCreateFaucet.FindFaucetBySharedID.sharedID", sharedID, err.Error())
		return "", err
	}

	amountFaucet := big.NewInt(0.05 * 1e18) // todo: move to config

	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		logger.AtLog.Logger.Error("ApiCreateFaucet.ParseBool", zap.Error(err))
		return "", err
	}

	chromePath := "google-chrome"

	// if u.Config.ENV == "develop" {
	// 	chromePath = ""
	// }

	address, err := getFaucetPaymentInfo(url, chromePath, eCH)
	fmt.Println("address, err: ", address, err)

	if err != nil {
		go u.sendSlack("", "ApiCreateFaucet.getFaucetPaymentInfo.url", url, err.Error())
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.getFaucetPaymentInfo"), zap.Error(err))
		return "", err
	}

	err = u.CheckValidFaucet(address, twName)
	if err != nil {
		go u.sendSlack("", "ApiCreateFaucet.CheckValidFaucet.(address+twName)", address+","+twName, err.Error())
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.checkValidFaucet"), zap.Error(err))
		return "", err
	}

	faucetItem := &entity.Faucet{
		Address:     address,
		TwitterName: twName,
		Status:      0,
		Tx:          "",
		Amount:      amountFaucet.String(),
		TwShareID:   sharedID,
	}
	err = u.Repo.InsertFaucet(faucetItem)
	if err != nil {
		return "", err
	}

	go u.sendSlack("", "ApiCreateFaucet.NewFaucet", twName+"/"+address, "ok")

	return "", nil

	/*

		// transfer now:
		privateKeyDeCrypt, err := encrypt.DecryptToString(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET"), os.Getenv("SECRET_KEY"))
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.Decrypt.%s.Error", "can decrypt"), zap.Error(err))
			return "", err
		}

		tx, err := u.TcClient.Transfer(privateKeyDeCrypt, faucetItem.Address, amountFaucet)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.Transfer"), zap.Error(err))
			faucetItem.ErrLogs += "|" + err.Error()
			u.Repo.UpdateFaucet(faucetItem)
			return "", err
		}
		faucetItem.Tx = tx
		faucetItem.Status = 1
		_, err = u.Repo.UpdateFaucet(faucetItem)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.UpdateFaucetByUUid"), zap.Error(err))
			return "", err
		}
		// submit raw data:
		txBtc, err := u.SubmitTCToBtcChain(tx, 10)
		if err != nil {
			faucetItem.ErrLogs += "|" + err.Error()
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.SubmitTCToBtcChain"), zap.Error(err))
		} else {
			faucetItem.BtcTx = txBtc
			faucetItem.Status = 2
		}

		_, err = u.Repo.UpdateFaucet(faucetItem)

		_, err = u.Repo.UpdateFaucetByUUid(faucetItem.UUID, tx, txBtc, 1)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.UpdateFaucetByUUid btc tx"), zap.Error(err))
			return "", err
		}

		return "https://explorer.trustless.computer/" + tx, nil
	*/
}

func (u Usecase) CheckValidFaucet(address, twName string) error {
	// find address to faucet:
	faucetItems, err := u.Repo.FindFaucetByTwitterNameOrAddress(twName, address)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByTwitterName"), zap.Error(err))
		return err
	}

	totalFaucet := len(faucetItems)

	fmt.Println("totalFaucet: ", totalFaucet)

	if totalFaucet >= 5 {
		// check 5 times:
		err = errors.New("You have reached the maximum faucet.")
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
		return err

	}
	if totalFaucet > 0 {
		// last item:
		lastItem := faucetItems[0]

		t1 := lastItem.CreatedAt
		t2 := time.Now()

		diff := t2.Sub(*t1)

		maxHours := float64(24)

		fmt.Println("diff.Hours(): ", diff.Hours())

		if diff.Hours() < maxHours {
			err = errors.New(fmt.Sprintf("You can only request once within 24 hours. Please wait another %0.1f minutes.", maxHours-diff.Hours()))
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
			return err
		}

	}
	return nil
}

//////////
func getFaucetPaymentInfo(url, chromePath string, eCH bool) (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),  // uncomment on the server.
		chromedp.Flag("headless", eCH), // false => open chrome. true on the server
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-first-run", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	var res string
	cctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	err := chromedp.Run(cctx,
		chromedp.EmulateViewport(960, 960),
		chromedp.Navigate(url),
		chromedp.Text(ByTestId("tweetText"), &res, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	spew.Dump(res)

	if !strings.Contains(res, "DappsOnBitcoin") {
		return "", errors.New("tweet data invalid")
	}

	addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}") // payment address eth

	addressHex := addressRegex.FindString(res)
	if len(addressHex) == 0 {
		return "", errors.New("address not found")
	}

	fmt.Println("result: ", addressHex)

	return addressHex, nil

}
func ByTestId(s string) string {
	return "[data-testid='" + s + "']"
}

// Job faucet now:
func (u Usecase) JobFaucet_SendTCNow() error {

	if len(os.Getenv("TC_MULTI_CONTRACT")) == 0 {
		err := errors.New("TC_MULTI_CONTRACT empty")
		go u.sendSlack("", "JobFaucet_SendTCNow.TC_MULTI_CONTRACT", "empty", err.Error())
		return err
	}
	if len(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET")) == 0 {
		err := errors.New("PRIVATE_KEY_FEE_TC_WALLET empty")
		go u.sendSlack("", "JobFaucet_SendTCNow.PRIVATE_KEY_FEE_TC_WALLET", "empty", err.Error())
		return err
	}
	if len(os.Getenv("SECRET_KEY")) == 0 {
		err := errors.New("SECRET_KEY empty")
		go u.sendSlack("", "JobFaucet_SendTCNow.SECRET_KEY", "empty", err.Error())
		return err
	}

	// check pending first:
	recordsPending, _ := u.Repo.FindFaucetByStatus(2)
	if len(recordsPending) > 0 {
		return u.JobFaucet_CheckTx(recordsPending)
	}

	faucets, _ := u.Repo.FindFaucetByStatus(0)
	fmt.Println("need faucet: ", len(faucets))

	if len(faucets) == 0 {
		return nil
	}

	// send TC:
	destinations := make(map[string]*big.Int)

	var uuids []string

	amountFaucet := big.NewInt(0.05 * 1e18) // todo: move to config

	feeRate := 6

	feeRateCurrent, err := u.getFeeRateFromChain()
	if err == nil {
		feeRate = feeRateCurrent.HourFee
	}

	// get list again:
	for _, item := range faucets {
		destinations[item.Address] = amountFaucet
		uuids = append(uuids, item.UUID)
	}

	uuidStr := strings.Join(uuids, ",")

	fmt.Println("destinations: ", destinations)

	privateKeyDeCrypt, err := encrypt.DecryptToString(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET"), os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("GenMintFreeTemAddress.Decrypt.%s.Error", "can decrypt"), zap.Error(err))
		return err
	}

	txID, err := u.TcClient.SendMulti(
		os.Getenv("TC_MULTI_CONTRACT"),
		privateKeyDeCrypt,
		destinations,
		nil,
		0,
	)
	fmt.Println("txID, err ", txID, err)

	if err != nil {
		go u.sendSlack(uuidStr, "ApiCreateFaucet.SendMulti", "call send "+txID, err.Error())
		return err
	}

	go u.sendSlack(uuidStr, "ApiCreateFaucet.SendMulti", "ok=> tx", txID)

	// submit raw data:
	txBtc, err := u.SubmitTCToBtcChain(txID, feeRate)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.SubmitTCToBtcChain"), zap.Error(err))
		go u.sendSlack(uuidStr, "ApiCreateFaucet.SubmitTCToBtcChain", "call send vs tcTx: "+txID, err.Error())
		return err
	}
	go u.sendSlack(uuidStr, "ApiCreateFaucet.SubmitTCToBtcChain", "ok=>tcTx/btcTx", txID+"/"+txBtc)
	// update tx by uuids:
	if len(uuids) > 0 {
		for _, item := range faucets {
			item.Status = 2
			item.Tx = txID
			item.BtcTx = txBtc

			_, err = u.Repo.UpdateFaucet(item)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.UpdateFaucet"), zap.Error(err))
				go u.sendSlack(uuidStr, "ApiCreateFaucet.UpdateFaucet", "UpdateFaucet", err.Error())
			}
			// return err
		}
	}

	fmt.Println("done------")

	return nil
}

func (u Usecase) sendSlack(ids, funcName, requestMsgStr, errorStr string) {
	preText := fmt.Sprintf("[App: %s][recordIDs %s] - %s", "Faucet", ids, requestMsgStr)
	if _, _, err := u.Slack.SendMessageToSlackWithChannel("C052K111MK6", preText, funcName, errorStr); err != nil {
		fmt.Println("s.Slack.SendMessageToSlack err", err)
	}
}

func (u Usecase) JobFaucet_CheckTx(recordsToCheck []*entity.Faucet) error {
	// recordsToCheck, _ := u.Repo.FindFaucetByStatus(2)
	// fmt.Println("need check tx: ", len(recordsToCheck))

	mapCheckTxPass := make(map[string]bool)
	mapCheckTxFalse := make(map[string]string)

	// check tc tx:
	for _, item := range recordsToCheck {
		context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		status, err := u.TcClient.GetTransaction(context, item.Tx)

		fmt.Println("GetTransaction status, err ", item.Tx, status, err)

		if err == nil {
			if status > 0 {
				// pass:
				mapCheckTxPass[item.Tx] = true
				_, err = u.Repo.UpdateFaucet(item)
				if err != nil {
					go u.sendSlack(item.UUID, "JobFaucet_CheckTx.UpdateFaucet", "UpdateFaucet", err.Error())
				}

			} else {
				mapCheckTxFalse[item.Tx] = "status != 1"
			}
		} else {
			// if error maybe tx is pending or rejected
			// TODO check timeout to detect tx is rejected or not.
			mapCheckTxFalse[item.Tx] = "err: " + err.Error()
		}
	}
	if len(mapCheckTxFalse) > 0 {
		var uuids []string
		var errs []string
		for k, v := range mapCheckTxFalse {
			uuids = append(uuids, k)
			errs = append(errs, v)
		}

		go u.sendSlack(strings.Join(uuids, ","), "JobFaucet_CheckTx.UpdateFaucet", "mapCheckTxFalse", strings.Join(errs, ","))
	}
	return nil

}
