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

	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiCreateFaucet(url string) (string, error) {

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

	address, twName, err := getFaucetPaymentInfo(url, chromePath, eCH)
	fmt.Println("address, err: ", address, err)

	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.getFaucetPaymentInfo"), zap.Error(err))
		return "", err
	}
	// find address to faucet:
	faucetItemByTwitterName, _ := u.Repo.FindFaucetByTwitterName(twName)
	faucetItem, _ := u.Repo.FindFaucetByAddress(address)
	if faucetItem != nil {

		if faucetItemByTwitterName == nil {
			err = errors.New("The account not match.")
			return "", err
		}
		if faucetItem.TwitterName != faucetItemByTwitterName.TwitterName || faucetItem.Address != faucetItemByTwitterName.Address {
			err = errors.New("The account not match.")
			return "", err
		}

		if faucetItem.Status > 0 {
			err = errors.New("The transaction already exists.")
			logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.FindFaucetByAddress"), zap.Error(err))
			return "", err
		}
	} else {
		faucetItem = &entity.Faucet{
			Address:     address,
			TwitterName: twName,
			Status:      0,
			Tx:          "",
			Amount:      amountFaucet.String(),
		}
		err = u.Repo.InsertFaucet(faucetItem)
		if err != nil {
			return "", err
		}
	}
	// transfer now:
	privateKeyDeCrypt, err := encrypt.DecryptToString(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET"), os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("GenMintFreeTemAddress.Decrypt.%s.Error", "can decrypt"), zap.Error(err))
		return "", err
	}

	tx, err := u.TcClient.Transfer(privateKeyDeCrypt, faucetItem.Address, amountFaucet)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.Transfer"), zap.Error(err))
		return "", err
	}
	// submit raw data:
	txBtc, err := u.SubmitTCToBtcChain(tx, 10)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.SubmitTCToBtcChain"), zap.Error(err))
		// return "", err
	}

	_, err = u.Repo.UpdateFaucetByUUid(faucetItem.UUID, tx, txBtc, 1)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ApiCreateFaucet.UpdateFaucetByUUid"), zap.Error(err))
		return "", err
	}

	return "https://explorer.trustless.computer/" + tx, nil
}

//////////
func getFaucetPaymentInfo(url, chromePath string, eCH bool) (string, string, error) {

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

	//https://twitter.com/2712_at1999/status/1643190049981480961

	if !strings.Contains(res, "DappsOnBitcoin") {
		return "", "", errors.New("tweet data invalid")
	}

	addressRegex := regexp.MustCompile("(0x)?[0-9a-fA-F]{40}")                // payment address eth
	twNameRegex := regexp.MustCompile(`https://twitter.com/(\w+)/status/\d+`) // payment address eth

	addressHex := addressRegex.FindString(res)
	if len(addressHex) == 0 {
		return "", "", errors.New("address not found")
	}

	// Find the first match in the tweet URL
	matchTwName := twNameRegex.FindStringSubmatch(url)

	twName := ""

	if len(matchTwName) >= 2 {
		twName = matchTwName[1]
		fmt.Println("twName:", twName) // Output: 2712_at1999
	} else {
		return "", "", errors.New("No username found in the tweet URL")
	}

	fmt.Println("result: ", addressHex)

	return addressHex, twName, nil

}
func ByTestId(s string) string {
	return "[data-testid='" + s + "']"
}
