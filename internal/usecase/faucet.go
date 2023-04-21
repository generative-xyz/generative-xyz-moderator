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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth/contract/tcartifact"
	"rederinghub.io/utils/eth/contract/tcbns"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiListCheckFaucet(address string) ([]*entity.Faucet, error) {
	faucetItems, err := u.Repo.FindFaucetByTwitterNameOrAddress(address, address)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByTwitterName"), zap.Error(err))
		return nil, err
	}
	for _, item := range faucetItems {
		if len(item.Tx) > 0 {
			item.Tx = "https://explorer.trustless.computer/tx/" + item.Tx
		}
		item.StatusStr = "Pending"
		if item.Status == 2 {
			item.StatusStr = "Processing"
		} else if item.Status == 3 {
			item.StatusStr = "Success"
		}
	}
	return faucetItems, nil

}

func (u Usecase) ApiCreateFaucet(addressInput, url, txhash, faucetType string) (string, error) {

	// verify tw name:
	// //https://twitter.com/2712_at1999/status/1643190049981480961
	// //https://twitter.com/abc/status/1647374585451663361?s=46&t=B7w70LBsAJFhv8XbJlpvCA
	twNameRegex := regexp.MustCompile(`https?://(?:www\.)?twitter\.com/([^/]+)/status/(\d+)(?:\?.*)?$`)
	// Find the first match in the tweet URL
	matchTwName := twNameRegex.FindStringSubmatch(url)

	twName := ""
	sharedID := ""

	if len(matchTwName) >= 3 {
		twName = matchTwName[1]
		sharedID = matchTwName[2]
		fmt.Println("twName:", twName)    // Output: 2712_at1999
		fmt.Println("shareID:", sharedID) // Output: 1643190049981480961

	} else {
		err := errors.New("The username or sharedID is not found in the tweet URL.")
		go u.sendSlack("", "ApiCreateFaucet.regexp", "check url"+url, err.Error())
		return "", err
	}
	// check valid vs twName first:
	err := u.CheckValidFaucet(addressInput, twName, txhash, faucetType)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.checkValidFaucet"), zap.Error(err))
		go u.sendSlack("", "ApiCreateFaucet.CheckValidFaucet.twName", twName, err.Error())
		return "", err
	}

	// check sharedID exist:
	sharedIDs, _ := u.Repo.FindFaucetBySharedID(sharedID)
	if len(sharedIDs) > 0 {
		err := errors.New("The tweet has already been used to claim the faucet. Please send out a new tweet.")
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetBySharedID"), zap.Error(err))
		go u.sendSlack("", "ApiCreateFaucet.FindFaucetBySharedID.sharedID", sharedID, err.Error())
		return "", err
	}

	amountFaucet := big.NewInt(0.1 * 1e18) // todo: move to config

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

	err = u.CheckValidFaucet(address, twName, txhash, faucetType)
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
		SharedLink:  url,
		UserTx:      txhash,
		FaucetType:  faucetType,
	}
	err = u.Repo.InsertFaucet(faucetItem)
	if err != nil {
		return "", err
	}

	go u.sendSlack("", "ApiCreateFaucet.NewFaucet", twName+"/"+address, "ok")

	return "The request was submitted successfully. You will receive TC after 1-2 block confirmations (10~20 minutes).", nil

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

const (
	BNSAddress      = "0x8b46F89BBA2B1c1f9eE196F43939476E79579798"
	ArtifaceAddress = "0x16efdc6d3f977e39dac0eb0e123feffed4320bc0"
)

func (u Usecase) CheckValidFaucet(address, twName, txhash, faucetType string) error {
	// find address to faucet:
	faucetItems, err := u.Repo.FindFaucetByTwitterNameOrAddress(twName, address)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByTwitterName"), zap.Error(err))
		return err
	}

	totalFaucet := len(faucetItems)
	filteredTotalFaucet := 0
	for _, item := range faucetItems {
		if item.FaucetType == faucetType {
			filteredTotalFaucet++
		}
	}

	fmt.Println("totalFaucet: ", totalFaucet)
	fmt.Println("filteredTotalFaucet: ", filteredTotalFaucet)
	limitFaucet := 1
	switch faucetType {
	case "dapps":
		//check valid mint tx
		if txhash != "" {
			// check tx:
			tx, isPending, err := u.TcClient.GetClient().TransactionByHash(context.Background(), common.HexToHash(txhash))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash"), zap.Error(err))
				return err
			}
			if isPending {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash.isPending"), zap.Error(err))
				return errors.New("tx is pending")
			}
			txReceipt, err := u.TcClient.GetClient().TransactionReceipt(context.Background(), common.HexToHash(txhash))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash"), zap.Error(err))
				return err
			}
			if txReceipt.Status == 0 {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash.Status"), zap.Error(err))
				return errors.New("tx status is 0")
			}

			from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.Sender"), zap.Error(err))
				return err
			}

			if !strings.EqualFold(from.String(), address) {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.Sender"), zap.Error(err))
				return errors.New("invalid sender")
			}

			if strings.EqualFold(tx.To().String(), BNSAddress) {
				bnsAddress := common.HexToAddress(BNSAddress)
				bnsContract, err := tcbns.NewBNS(bnsAddress, u.TcClient.GetClient())
				if err != nil {
					return err
				}
				haveEvent := false
				for _, v := range txReceipt.Logs {
					evt, err := bnsContract.ParseNameRegistered(*v)
					if err != nil {
						continue
					}
					if len(evt.Name) > 0 {
						haveEvent = true
						break
					}
				}
				if !haveEvent {
					return errors.New("invalid tx")
				}
			}
			if strings.EqualFold(tx.To().String(), ArtifaceAddress) {
				artifactAddress := common.HexToAddress(ArtifaceAddress)

				artifactContract, err := tcartifact.NewNFT721(artifactAddress, u.TcClient.GetClient())
				if err != nil {
					return err
				}
				haveEvent := false
				for _, v := range txReceipt.Logs {
					evt, err := artifactContract.ParseTransfer(*v)
					if err != nil {
						continue
					}
					if strings.EqualFold(evt.To.String(), address) {
						haveEvent = true
						break
					}
				}
				if !haveEvent {
					return errors.New("invalid tx")
				}
			}
		} else {
			return errors.New("invalid tx")
		}
	}

	if filteredTotalFaucet >= limitFaucet {
		// check times:
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
			err = errors.New(fmt.Sprintf("The faucet only allows one request per day. Please try again later in %0.1f hours.", maxHours-diff.Hours()))
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
			return err
		}

	}
	return nil
}

// ////////
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

	if !strings.Contains(res, "@generative_xyz") {
		return "", errors.New("Tweet not found. Please double-check and try again")
	}
	addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}") // payment address eth

	addressHex := addressRegex.FindString(res)
	if len(addressHex) == 0 {
		return "", errors.New("Address not found.")
	}

	fmt.Println("result: ", addressHex)

	return addressHex, nil
}

func getFaucetInfo(url, chromePath string, eCH bool) (string, string, error) {

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
		return "", "", err
	}

	spew.Dump(res)

	if !strings.Contains(res, "@generative_xyz") {
		return "", "", errors.New("Tweet not found. Please double-check and try again")
	}
	addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}") // payment address eth

	addressHex := addressRegex.FindString(res)

	txRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{64}") // tx eth
	txHex := txRegex.FindString(res)

	return addressHex, txHex, nil
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
		u.JobFaucet_CheckTx(recordsPending)

	}

	// check pending again:
	recordsPending, _ = u.Repo.FindFaucetByStatus(2)
	if len(recordsPending) > 0 {
		return nil
	}

	faucets, _ := u.Repo.FindFaucetByStatus(0)
	fmt.Println("need faucet: ", len(faucets))

	if len(faucets) == 0 {
		return nil
	}

	// send TC:
	destinations := make(map[string]*big.Int)

	var uuids []string

	amountFaucet := big.NewInt(0.1 * 1e18) // todo: move to config

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
				item.Status = 3
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
