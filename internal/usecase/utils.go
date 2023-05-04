package usecase

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/erc20"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type FeeRates struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

type FeeRateInfo struct {
	FeeRate     int                           `json:"rate"`
	MintFeeInfo map[string]entity.MintFeeInfo `json:"mintFees"`
}

func (u Usecase) networkFeeBySize(size int64) int64 {

	feeRateValue := int64(entity.DEFAULT_FEE_RATE)

	response, err := http.Get("https://mempool.space/api/v1/fees/recommended")

	if err != nil {
		fmt.Print(err.Error())
		return size * feeRateValue
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return size * feeRateValue
	} else {

		feeRateObj := FeeRates{}

		err = json.Unmarshal(responseData, &feeRateObj)
		if err != nil {
			logger.AtLog.Logger.Error("err", zap.Error(err))
			return size * feeRateValue
		}
		if feeRateObj.FastestFee > 0 {
			feeRateValue = int64(feeRateObj.FastestFee)
		}
	}

	return size * feeRateValue

}

func (u Usecase) getFeeRateFromChain() (*FeeRates, error) {

	response, err := http.Get("https://mempool.space/api/v1/fees/recommended")

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println("responseData", string(responseData))

	feeRateObj := &FeeRates{}

	err = json.Unmarshal(responseData, &feeRateObj)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}
	return feeRateObj, nil

}

func (u Usecase) GetLevelFeeInfo(fileSize, customRate, mintPrice int64) (map[string]FeeRateInfo, error) {

	var btcRate, ethRate float64

	btcRate, err := helpers.GetExternalPrice("BTC")
	if err != nil {
		return nil, err
	}

	ethRate, err = helpers.GetExternalPrice("ETH")
	if err != nil {
		return nil, err
	}

	levelFeeFullInfo := make(map[string]FeeRateInfo)

	feeRateFromChain, err := u.getFeeRateFromChain()
	if err != nil {
		return nil, err
	}

	fastestFee := feeRateFromChain.FastestFee
	fasterFee := feeRateFromChain.HalfHourFee
	ecoFee := feeRateFromChain.EconomyFee // not use HourFee

	// fmt.Println("fastestFee", fastestFee)
	// fmt.Println("fasterFee", fasterFee)
	// fmt.Println("ecoFee", ecoFee)

	min := 4 // todo: move config

	// set mint 5 if fee < 5.
	if fastestFee < min {
		fastestFee = min
	}
	if fasterFee < min {
		fasterFee = min
	}
	if ecoFee < min {
		ecoFee = min
	}

	fastestMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(fastestFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	fasterMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(fasterFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	economyMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(ecoFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}

	levelFeeFullInfo["fastest"] = FeeRateInfo{
		FeeRate:     fastestFee,
		MintFeeInfo: fastestMintInfo,
	}
	levelFeeFullInfo["faster"] = FeeRateInfo{
		FeeRate:     fasterFee,
		MintFeeInfo: fasterMintInfo,
	}
	levelFeeFullInfo["economy"] = FeeRateInfo{
		FeeRate:     ecoFee,
		MintFeeInfo: economyMintInfo,
	}

	if customRate > 0 {
		customRateMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(customRate), btcRate, ethRate)
		if err != nil {
			return nil, err
		}
		levelFeeFullInfo["customRate"] = FeeRateInfo{
			FeeRate:     int(customRate),
			MintFeeInfo: customRateMintInfo,
		}
	}

	return levelFeeFullInfo, nil
}

func (u Usecase) NotifyWithChannel(channelID string, title string, userAddress string, content string) {
	//slack
	preText := fmt.Sprintf("[App: %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlackWithChannel(channelID, preText, title, c); err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err), zap.String("channelID", channelID), zap.String("title", title), zap.String("userAddress", userAddress))
	}
}

func (u Usecase) Notify(title string, userAddress string, content string) {

	//slack
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, c); err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
	}
}

func (u Usecase) IsWhitelistedAddress(ctx context.Context, userAddr string, whitelistedAddrs []string) (bool, error) {

	logger.AtLog.Logger.Info("whitelistedAddrs", zap.Any("whitelistedAddrs", whitelistedAddrs))
	if len(whitelistedAddrs) == 0 {
		logger.AtLog.Logger.Info("whitelistedAddrs.Total", zap.Any("len(whitelistedAddrs)", len(whitelistedAddrs)))
		return false, nil
	}
	filter := nfts.MoralisFilter{}
	filter.Limit = new(int)
	*filter.Limit = 1
	filter.TokenAddresses = new([]string)
	*filter.TokenAddresses = whitelistedAddrs

	logger.AtLog.Logger.Info("filter.GetNftByWalletAddress", zap.Any("filter", filter))
	resp, err := u.MoralisNft.GetNftByWalletAddress(userAddr, filter)
	if err != nil {
		logger.AtLog.Logger.Error("u.MoralisNft.GetNftByWalletAddress", zap.Error(err))
		return false, err
	}

	logger.AtLog.Logger.Info("resp", zap.Any("resp", resp))
	if len(resp.Result) > 0 {
		return true, nil
	}

	delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userAddr)
	if err != nil {
		logger.AtLog.Logger.Error("u.DelegateService.GetDelegationsByDelegate", zap.Error(err))
		return false, err
	}

	logger.AtLog.Logger.Info("delegations", zap.Any("delegations", delegations))
	for _, delegation := range delegations {
		if containsIgnoreCase(whitelistedAddrs, delegation.Contract.String()) {
			return true, nil
		}
	}
	return false, nil
}

func (u Usecase) IsWhitelistedAddressERC20(ctx context.Context, userAddr string, erc20WhiteList map[string]structure.Erc20Config) (bool, error) {
	client, err := helpers.EthDialer()
	if err != nil {
		return false, err
	}

	for addr, whitelistedThres := range erc20WhiteList {
		erc20Client, err := erc20.NewErc20(common.HexToAddress(addr), client)
		if err != nil {
			continue
		}

		blance, err := erc20Client.BalanceOf(nil, common.HexToAddress(userAddr))
		if err != nil {
			continue
		}

		pow := new(big.Int)
		pow = pow.Exp(big.NewInt(1), big.NewInt(whitelistedThres.Decimal), nil)
		confValue := big.NewInt(whitelistedThres.Value)

		confValue = confValue.Mul(confValue, pow)

		//bigInt64 := big.
		tmp := blance.Cmp(confValue)

		spew.Dump(whitelistedThres.Value, whitelistedThres.Decimal)
		spew.Dump(confValue.String())
		spew.Dump(blance.String())
		if tmp >= 0 {
			return true, nil
		}
	}

	return false, nil
}

// // containsIgnoreCase ...

func containsIgnoreCase(strSlice []string, item string) bool {
	for _, str := range strSlice {
		if strings.EqualFold(str, item) {
			return true
		}
	}

	return false
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func Commas(s string) string {
	n := len(s)
	if n < 4 {
		return s
	}
	return Commas(s[:n-3]) + "," + s[n-3:]
}
