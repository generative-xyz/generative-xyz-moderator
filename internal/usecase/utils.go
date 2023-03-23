package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/erc20"
	"rederinghub.io/utils/helpers"
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
			u.Logger.Error(err)
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

	fmt.Println("responseData", string(responseData))

	feeRateObj := &FeeRates{}

	err = json.Unmarshal(responseData, &feeRateObj)
	if err != nil {
		u.Logger.Error(err)
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

	fmt.Println("fastestFee", feeRateFromChain.FastestFee)
	fmt.Println("halfHourFee", feeRateFromChain.HalfHourFee)
	fmt.Println("hourFee", feeRateFromChain.HourFee)

	fastestMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(feeRateFromChain.FastestFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	fasterMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(feeRateFromChain.HalfHourFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	economyMintInfo, err := u.calMintFeeInfo(mintPrice, fileSize, int64(feeRateFromChain.HourFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}

	levelFeeFullInfo["fastest"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.FastestFee,
		MintFeeInfo: fastestMintInfo,
	}
	levelFeeFullInfo["faster"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.HalfHourFee,
		MintFeeInfo: fasterMintInfo,
	}
	levelFeeFullInfo["economy"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.HourFee,
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
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlackWithChannel(channelID, preText, title, c); err != nil {
		u.Logger.Error(err)
	}
}

func (u Usecase) Notify(title string, userAddress string, content string) {

	//slack
	preText := fmt.Sprintf("[App: %s][traceID %s] - User address: %s, ", os.Getenv("JAEGER_SERVICE_NAME"), "", userAddress)
	c := fmt.Sprintf("%s", content)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, c); err != nil {
		u.Logger.Error(err)
	}
}

func (u Usecase) IsWhitelistedAddress(ctx context.Context, userAddr string, whitelistedAddrs []string) (bool, error) {

	u.Logger.Info("whitelistedAddrs", whitelistedAddrs)
	if len(whitelistedAddrs) == 0 {
		u.Logger.Info("whitelistedAddrs.Total", len(whitelistedAddrs))
		return false, nil
	}
	filter := nfts.MoralisFilter{}
	filter.Limit = new(int)
	*filter.Limit = 1
	filter.TokenAddresses = new([]string)
	*filter.TokenAddresses = whitelistedAddrs

	u.Logger.Info("filter.GetNftByWalletAddress", filter)
	resp, err := u.MoralisNft.GetNftByWalletAddress(userAddr, filter)
	if err != nil {
		u.Logger.Error("u.MoralisNft.GetNftByWalletAddress", err.Error(), err)
		return false, err
	}

	u.Logger.Info("resp", resp)
	if len(resp.Result) > 0 {
		return true, nil
	}

	delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userAddr)
	if err != nil {
		u.Logger.Error("u.DelegateService.GetDelegationsByDelegate", err.Error(), err)
		return false, err
	}

	u.Logger.Info("delegations", delegations)
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
