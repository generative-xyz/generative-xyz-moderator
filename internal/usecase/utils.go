package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"rederinghub.io/internal/entity"
)

type FeeRates struct {
	fastestFee  int
	halfHourFee int
	hourFee     int
	economyFee  int
	minimumFee  int
}

type FeeRateInfo struct {
	FeeRate     int
	MintFeeInfo map[string]entity.MintFeeInfo
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
		if feeRateObj.fastestFee > 0 {
			feeRateValue = int64(feeRateObj.fastestFee)
		}
	}

	return size * feeRateValue

}

func (u Usecase) getFeeRateFromChain(size int64) (*FeeRates, error) {

	response, err := http.Get("https://mempool.space/api/v1/fees/recommended")

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	feeRateObj := FeeRates{}

	err = json.Unmarshal(responseData, &feeRateObj)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}
	return &feeRateObj, nil

}

func (u Usecase) GetLevelFeeInfo(fileSize int64) (map[string]FeeRateInfo, error) {

	levelFeeFullInfo := make(map[string]FeeRateInfo)

	feeRateFromChain, err := u.getFeeRateFromChain(fileSize)
	if err != nil {
		return nil, err
	}

	fastestMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.fastestFee))
	if err != nil {
		return nil, err
	}
	fasterMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.halfHourFee))
	if err != nil {
		return nil, err
	}
	economyMintInfo, err := u.calMintFeeInfo(0, fileSize, int64(feeRateFromChain.hourFee))
	if err != nil {
		return nil, err
	}

	levelFeeFullInfo["fastest"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.fastestFee,
		MintFeeInfo: fastestMintInfo,
	}
	levelFeeFullInfo["faster"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.halfHourFee,
		MintFeeInfo: fasterMintInfo,
	}
	levelFeeFullInfo["economy"] = FeeRateInfo{
		FeeRate:     feeRateFromChain.hourFee,
		MintFeeInfo: economyMintInfo,
	}
	return levelFeeFullInfo, nil
}
