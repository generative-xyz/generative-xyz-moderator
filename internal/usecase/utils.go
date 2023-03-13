package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"rederinghub.io/internal/entity"
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

func (u Usecase) GetLevelFeeInfo(fileSize, customRate int64) (map[string]FeeRateInfo, error) {

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

	fastestMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.FastestFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	fasterMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.HalfHourFee), btcRate, ethRate)
	if err != nil {
		return nil, err
	}
	economyMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.HourFee), btcRate, ethRate)
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
		customRateMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(customRate), btcRate, ethRate)
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
