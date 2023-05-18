package usecase

import (
	"bytes"
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
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	faucetconst "rederinghub.io/utils/constants/faucet"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/eth/contract/tcartifact"
	"rederinghub.io/utils/eth/contract/tcbns"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiFaucetGetNonce(address string) ([]uint64, error) {
	var nonces []uint64

	nonce1, err := u.TcClient.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		return nonces, err
	}
	nonces = append(nonces, nonce1)
	tcClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.TCPublicEndpoint)
	if err != nil {
		return nonces, err
	}
	tcPulicClient := eth.NewClient(tcClientWrap)
	nonce2, err := tcPulicClient.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		return nonces, err
	}
	nonces = append(nonces, nonce2)
	return nonces, nil

}

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

func (u Usecase) ApiCreateFaucet(addressInput, url, txhash, faucetType, source string) (string, error) {

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
	// specFaucetType, err := u.CheckValidFaucet(addressInput, twName, txhash, faucetType)
	// if err != nil {
	// 	logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.checkValidFaucet"), zap.Error(err))
	// 	go u.sendSlack("", "ApiCreateFaucet.CheckValidFaucet.twName", twName, err.Error())
	// 	return "", err
	// }

	// check sharedID exist:
	sharedIDs, _ := u.Repo.FindFaucetBySharedID(sharedID)
	if len(sharedIDs) > 0 {
		err := errors.New("The tweet has already been used to claim the faucet. Please send out a new tweet.")
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetBySharedID"), zap.Error(err))
		go u.sendSlack("", "ApiCreateFaucet.FindFaucetBySharedID.sharedID", sharedID, err.Error())
		return "", err
	}

	amountFaucet := big.NewInt(0.1 * 1e18) // todo: move to config

	switch source {
	case "special":
		if faucetType == "" {
			amountFaucet = big.NewInt(faucetconst.SpecialFaucetAmount)
		}
	}
	eCH, err := strconv.ParseBool(os.Getenv("ENABLED_CHROME_HEADLESS"))
	if err != nil {
		logger.AtLog.Logger.Error("ApiCreateFaucet.ParseBool", zap.Error(err))
		return "", err
	}

	chromePath := "google-chrome"

	// if u.Config.ENV == "develop" {
	// 	chromePath = ""
	// }
	var address string
	var contractAddress string
	address, txhash, contractAddress, err = getFaucetPaymentInfo(url, chromePath, eCH)
	fmt.Println("address, err: ", address, err)

	if err != nil {
		go u.sendSlack("", "ApiCreateFaucet.getFaucetPaymentInfo.url", url, err.Error())
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.getFaucetPaymentInfo"), zap.Error(err))
		return "", err
	}
	if address == "" && txhash == "" && contractAddress == "" {
		err := errors.New("The address or txhash is not found in the tweet URL.")
		return "", err
	}

	if txhash != "" {
		address = addressInput
	}
	var specFaucetType string
	if contractAddress == "" || (contractAddress != "" && faucetType != "dev") {
		specFaucetType, err = u.CheckValidFaucet(address, twName, txhash, faucetType)
		if err != nil {
			go u.sendSlack("", "ApiCreateFaucet.CheckValidFaucet.(address+twName)", address+","+twName, err.Error())
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.checkValidFaucet"), zap.Error(err))
			return "", err
		}
	} else {
		specFaucetType = "dev"
		address = addressInput
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
		Contract:    contractAddress,
		FaucetType:  specFaucetType,
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

func (u Usecase) CheckValidFaucet(address, twName, txhash, faucetType string) (string, error) {
	// find address to faucet:
	faucetItems, err := u.Repo.FindFaucetByTwitterNameOrAddress(twName, address)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByTwitterName"), zap.Error(err))
		return "", err
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
	specFaucetType := ""
	switch faucetType {
	case "dev":

	case "dapps":
		//check valid mint tx
		if txhash != "" {
			// check tx:
			tx, isPending, err := u.TcClient.GetClient().TransactionByHash(context.Background(), common.HexToHash(txhash))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash"), zap.Error(err))
				return specFaucetType, errors.New("tx not found")
			}
			if isPending {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash.isPending"), zap.Error(err))
				return specFaucetType, errors.New("tx is pending")
			}
			txReceipt, err := u.TcClient.GetClient().TransactionReceipt(context.Background(), common.HexToHash(txhash))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash"), zap.Error(err))
				return specFaucetType, err
			}
			if txReceipt.Status == 0 {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.TransactionByHash.Status"), zap.Error(err))
				return specFaucetType, errors.New("tx have failed status")
			}

			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.Sender"), zap.Error(err))
				return specFaucetType, err
			}

			if !strings.EqualFold(from.String(), address) {
				logger.AtLog.Logger.Error(fmt.Sprintf("CheckValidFaucet.Sender"), zap.Error(err))
				return specFaucetType, errors.New("requestor is not tx sender")
			}

			if strings.EqualFold(tx.To().String(), BNSAddress) {
				bnsAddress := common.HexToAddress(BNSAddress)
				bnsContract, err := tcbns.NewBNS(bnsAddress, u.TcClient.GetClient())
				if err != nil {
					return specFaucetType, err
				}
				haveEvent := false
				for _, v := range txReceipt.Logs {
					evt, err := bnsContract.ParseNameRegistered(*v)
					if err != nil {
						continue
					}
					if len(evt.Name) > 0 {
						specFaucetType = "bns"
						haveEvent = true
						break
					}
				}
				if !haveEvent {
					return specFaucetType, errors.New("tx is not mint artifact/bns")
				}
			}
			if strings.EqualFold(tx.To().String(), ArtifaceAddress) {
				artifactAddress := common.HexToAddress(ArtifaceAddress)

				artifactContract, err := tcartifact.NewNFT721(artifactAddress, u.TcClient.GetClient())
				if err != nil {
					return specFaucetType, err
				}
				haveEvent := false
				for _, v := range txReceipt.Logs {
					evt, err := artifactContract.ParseTransfer(*v)
					if err != nil {
						continue
					}
					if strings.EqualFold(evt.To.String(), address) {
						specFaucetType = "artifact"
						haveEvent = true
						break
					}
				}
				if !haveEvent {
					return specFaucetType, errors.New("tx is not mint artifact/bns")
				}
			}
		} else {
			return specFaucetType, errors.New("tx not found in tweet")
		}
	}

	if filteredTotalFaucet >= limitFaucet {
		// check times:
		err = errors.New("This Twitter account already claimed the faucet.")
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
		return specFaucetType, err
	}

	// if totalFaucet > 0 {
	// 	// last item:
	// 	lastItem := faucetItems[0]

	// 	t1 := lastItem.CreatedAt
	// 	t2 := time.Now()

	// 	diff := t2.Sub(*t1)

	// 	maxHours := float64(24)

	// 	fmt.Println("diff.Hours(): ", diff.Hours())

	// 	if diff.Hours() < maxHours {
	// 		err = errors.New(fmt.Sprintf("The faucet only allows one request per day. Please try again later in %0.1f hours.", maxHours-diff.Hours()))
	// 		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
	// 		return specFaucetType, err
	// 	}

	// }
	return specFaucetType, nil
}

// ////////
func getFaucetPaymentInfo(url, chromePath string, eCH bool) (string, string, string, error) {

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
		return "", "", "", err
	}

	spew.Dump(res)

	// if !strings.Contains(res, "@generative_xyz") {
	// 	return "", errors.New("Tweet not found. Please double-check and try again")
	// }
	var addressHex string
	var txHex string
	var contractAddress string
	res = strings.ToLower(res)
	if strings.Contains(res, "my tc address is:") {
		addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}") // payment address eth
		texts := strings.Split(res, "my tc address is:")
		addressHex = addressRegex.FindString(texts[1])
	}
	if strings.Contains(res, "my transaction id is:") {
		txRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{64}") // payment address eth
		texts := strings.Split(res, "my transaction id is:")
		txHex = txRegex.FindString(texts[1])
	}
	if strings.Contains(res, "contract address:") {
		addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}") // payment address eth
		texts := strings.Split(res, "contract address:")
		contractAddress = addressRegex.FindString(texts[1])
	}

	fmt.Println("result: ", addressHex, txHex)

	return addressHex, txHex, contractAddress, nil
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

	needRB, _ := u.Repo.FindFaucetByTx("0x61695550da1173eea02488474176c90ac062bc948b67d56f9aabeece53bbdd7f")

	if len(needRB) > 0 {
		for _, v := range needRB {
			v.Status = 0
			v.Tx = ""
			v.BtcTx = ""
			v.ErrLogs = "retry miss send tc"
			_, err := u.Repo.UpdateFaucet(v)
			if err != nil {
				go u.sendSlack(v.UUID, "ApiCreateFaucet.UpdateFaucet", "UpdateFaucet", err.Error())
				return err
			}
		}
		return nil
	}

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

	feeRate := 6

	feeRateCurrent, err := u.getFeeRateFromChain()
	if err == nil {
		feeRate = feeRateCurrent.FastestFee
	}

	faucetNeedTrigger, _ := u.Repo.FindFaucetByStatus(1)
	fmt.Println("faucetNeedTrigger len: ", len(faucetNeedTrigger))

	if len(faucetNeedTrigger) > 0 {
		// submit raw data:
		tempItem := faucetNeedTrigger[0]

		// check tx tc first:
		context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		status, err := u.TcClient.GetTransaction(context, tempItem.Tx)
		fmt.Println("GetTransaction status, err ", tempItem.Tx, status, err)
		if err == nil {
			if status > 0 {
				// pass:
				_, err = u.Repo.UpdateStatusFaucetByTxTc(tempItem.Tx, 3)
				if err != nil {
					go u.sendSlack(tempItem.UUID, "JobFaucet_CheckTx.UpdateFaucet", "UpdateFaucet", err.Error())
				}
				go u.sendSlack(tempItem.UUID, "JobFaucet_CheckTx.UpdateStatusFaucetByTxTc", "Update status 3 before Re-Trigger: ", tempItem.Tx)
				return nil
			}

		} else {
			go u.sendSlack(tempItem.UUID, "JobFaucet_CheckTx.GetTransaction", "CheckTxBefore Re-Trigger: ", err.Error())
		}

		txBtc, err := u.SubmitTCToBtcChain(tempItem.Tx, feeRate)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.SubmitTCToBtcChain"), zap.Error(err))
			go u.sendSlack(tempItem.UUID, "ApiCreateFaucet.Re-SubmitTCToBtcChain", "call send vs tcTx: "+tempItem.Tx, err.Error())
			return err
		}
		// update for tx:
		_, err = u.Repo.UpdateFaucetByTxTc(tempItem.Tx, txBtc, 2)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.UpdateFaucetByTxTc"), zap.Error(err))
			go u.sendSlack(tempItem.UUID, "ApiCreateFaucet.Re-SubmitTCToBtcChain.UpdateFaucetByTxTc", "update by tx err: "+tempItem.Tx+", btcTx:"+txBtc, err.Error())
			return err
		}
		go u.sendSlack(uuidStr, "ApiCreateFaucet.Re-SubmitTCToBtcChain", "okk=>tcTx/btcTx", "https://explorer.trustless.computer/tx/"+tempItem.Tx+"/https://mempool.space/tx/"+txBtc)
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
	// maxFaucet := big.NewInt(7 * 1e18)      // todo: move to config

	totalAmount := big.NewInt(0)

	t := 0

	var faucetsSent []*entity.Faucet

	// get list again:
	for _, item := range faucets {

		t += 1

		if t >= 200 {
			break
		}

		item.Address = strings.ToLower(item.Address)

		if !eth.ValidateAddress(item.Address) {
			fmt.Println("faucet valid address: ", item.Address)
			continue
		}
		if _, ok := destinations[item.Address]; ok {			
			continue
		}

		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			amount = big.NewInt(0).SetUint64(amountFaucet.Uint64())
		}
		// if amount.Uint64() > maxFaucet.Uint64() || amount.Uint64() == 0 {
		// 	amount = big.NewInt(0).SetUint64(amountFaucet.Uint64())
		// }

		totalAmount = big.NewInt(0).Add(totalAmount, amount)

		destinations[item.Address] = amount
		uuids = append(uuids, item.UUID)
		faucetsSent = append(faucetsSent, item)
	}
	
	if len(destinations) == 0 {
		return nil
	}

	uuidStr := strings.Join(uuids, ",")

	fmt.Println("destinations: ", destinations)

	privateKeyDeCrypt, err := encrypt.DecryptToString(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET"), os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("GenMintFreeTemAddress.Decrypt.%s.Error", "can decrypt"), zap.Error(err))
		return err
	}

	go u.sendSlack(fmt.Sprintf("%d", len(uuids)), "ApiCreateFaucet.SendMulti", "call send with total amount:", totalAmount.String())

	txID, err := u.TcClient.SendMulti(
		os.Getenv("TC_MULTI_CONTRACT"),
		privateKeyDeCrypt,
		destinations,
		totalAmount,
		0,
	)
	fmt.Println("txID, err ", txID, err)

	if err != nil {
		go u.sendSlack(uuidStr, "ApiCreateFaucet.SendMulti", fmt.Sprintf("call send %s err", totalAmount.String()), err.Error())
		return err
	}

	go u.sendSlack(uuidStr, "ApiCreateFaucet.SendMulti", "ok=> tx", txID)

	// update status 1 first:
	if len(uuids) > 0 {
		for _, item := range faucetsSent {
			item.Status = 1
			item.Tx = txID

			u.Repo.UpdateFaucet(item)
		}
	}

	// submit raw data:
	txBtc, err := u.SubmitTCToBtcChain(txID, feeRate)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.SubmitTCToBtcChain"), zap.Error(err))
		go u.sendSlack(uuidStr, "ApiCreateFaucet.SubmitTCToBtcChain", "call send vs tcTx: "+txID, err.Error())
		return err
	}
	
	go u.sendSlack(uuidStr, "ApiCreateFaucet.SubmitTCToBtcChain", "okk=>tcTx/btcTx", "https://explorer.trustless.computer/tx/"+txID+"/https://mempool.space/tx/"+txBtc)
	// update tx by uuids:
	if len(uuids) > 0 {
		for _, item := range faucetsSent {
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
	channel := "C052K111MK6"
	if strings.Contains(errorStr, "okk") || strings.Contains(preText, "okk") || strings.Contains(funcName, "okk") {
		channel = "C0582QV7MQD"
	}

	if _, _, err := u.Slack.SendMessageToSlackWithChannel(channel, preText, funcName, errorStr); err != nil {
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
			if strings.Contains(err.Error(), "not found") {
				now := time.Now()
				updatedTime := item.UpdatedAt
				if updatedTime != nil {

					duration := now.Sub(*updatedTime).Minutes()
					if duration >= 30 {
						u.sendSlack(item.UUID, "JobFaucet_CheckTx", fmt.Sprintf("long time to confirm okk? tcTx: https://explorer.trustless.computer/tx/%s, btcTx: https://mempool.space/tx/%s", item.Tx, item.BtcTx), fmt.Sprintf("%.2f mins ago", duration))
					}
				}
			}
			
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

// admin:
func (u Usecase) ApiAdminCreateFaucet(addressInput, url, txhash, faucetType, source string) (string, error) {

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

	}

	amountFaucet := big.NewInt(0.1 * 1e18) // todo: move to config

	faucetItem := &entity.Faucet{
		Address:     addressInput,
		TwitterName: twName,
		Status:      0,
		Tx:          "",
		Amount:      amountFaucet.String(),
		TwShareID:   sharedID,
		SharedLink:  url,
		UserTx:      txhash,
		FaucetType:  faucetType,
	}
	err := u.Repo.InsertFaucet(faucetItem)
	if err != nil {
		return "", err
	}

	go u.sendSlack("", "ApiAdminCreateFaucet.NewFaucet", twName+"/"+addressInput, "ok")

	return "The request was submitted successfully. You will receive TC after 1-2 block confirmations (10~20 minutes).", nil

}

func (u Usecase) ApiAdminCreateBatchFaucet(addresses []string, url, types string, amount float64) (string, error) {

	if amount == 0 || amount > 2 {
		amount = 0.1
	}

	uint64Value := amount * 1e18

	amountFaucet := big.NewInt(0).SetUint64(uint64(uint64Value))

	fmt.Println("amountFaucet: ", amountFaucet)

	for _, address := range addresses {
		faucetItem := &entity.Faucet{
			Address:    address,
			Status:     0,
			Tx:         "",
			Amount:     amountFaucet.String(),
			SharedLink: url,
			FaucetType: types,
		}
		err := u.Repo.InsertFaucet(faucetItem)
		if err != nil {
			return "", err
		}
	}

	go u.sendSlack("", "ApiAdminCreateBatchFaucet.NewFaucet", strings.Join(addresses, ","), "ok")

	return "The request was submitted successfully. You will receive TC after 1-2 block confirmations (10~20 minutes).", nil

}

func (u Usecase) ApiAdminCreateMapFaucet(addressAmountMap map[string]float64, url, types string) (string, error) {

	if len(addressAmountMap) == 0 {
		return "", nil
	}
	var totalAmount float64

	for address, amount := range addressAmountMap {

		totalAmount += amount

		uint64Value := amount * 1e18

		amountFaucet := big.NewInt(0).SetUint64(uint64(uint64Value))

		fmt.Println("amountFaucet: ", amountFaucet)
		fmt.Println("address: ", address)

		faucetItem := &entity.Faucet{
			Address:    address,
			Status:     0,
			Tx:         "",
			Amount:     amountFaucet.String(),
			SharedLink: url,
			FaucetType: types,
		}
		err := u.Repo.InsertFaucet(faucetItem)
		if err != nil {
			return "", err
		}
	}

	go u.sendSlack("", "ApiAdminCreateBatchFaucet.NewFaucet", createKeyValuePairs(addressAmountMap), fmt.Sprintf("total: %.4f", totalAmount))
	fmt.Println(fmt.Sprintf("total: %.4f", totalAmount))

	return "The request was submitted successfully. You will receive TC after 1-2 block confirmations (10~20 minutes).", nil

}

func createKeyValuePairs(m map[string]float64) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%f\"\n", key, value)
	}
	fmt.Println("b.String()", b.String())
	return b.String()
}
