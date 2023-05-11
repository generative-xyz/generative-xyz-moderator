package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"math/big"
	"os"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"strings"
)

func (u Usecase) GetChartDataOFProject(req structure.AggerateChartForProject) (*structure.AggragetedCollectionVolumnResp, error) {

	pe := &entity.AggerateChartForProject{}
	err := copier.Copy(pe, req)
	if err != nil {
		return nil, err
	}

	res := []entity.AggragetedProject{}
	if helpers.IsOrdinalProject(*req.ProjectID) {
		res, err = u.Repo.AggregateVolumnCollection(pe)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = u.Repo.AggregateVolumnCollectionTC(pe)
		if err != nil {
			return nil, err
		}
	}

	resp := []structure.AggragetedCollection{}
	for _, item := range res {
		tmp := structure.AggragetedCollection{
			ProjectID:    item.ID.ProjectID,
			ProjectName:  item.ID.ProjectName,
			Timestamp:    item.ID.Timestamp,
			Amount:       item.Amount,
			Erc20Address: item.ID.Erc20Address,
		}

		resp = append(resp, tmp)
	}

	return &structure.AggragetedCollectionVolumnResp{Volumns: resp}, nil
}

func (u Usecase) GetChartDataOFTokens(req structure.AggerateChartForToken) (*structure.AggragetedTokenVolumnResp, error) {

	pe := &entity.AggerateChartForToken{}
	err := copier.Copy(pe, req)
	if err != nil {
		return nil, err
	}

	res, err := u.Repo.AggregateVolumnToken(pe)
	if err != nil {
		return nil, err
	}

	resp := []structure.AggragetedTokenURI{}
	for _, item := range res {
		tmp := structure.AggragetedTokenURI{
			TokenID:   item.ID.TokenID,
			Timestamp: item.ID.Timestamp,
			Amount:    item.Amount,
		}

		resp = append(resp, tmp)
	}

	return &structure.AggragetedTokenVolumnResp{Volumns: resp}, nil
}

func (u Usecase) GetChartDataEthForGMCollection(tcAddress string, gmAddress string) (*structure.AnalyticsProjectDeposit, error) {
	ethRate, err := helpers.GetExternalPrice(string(entity.ETH))
	if err != nil {
		return nil, err
	}

	ethBL, err := u.EtherscanService.AddressBalance(gmAddress)
	if err != nil {
		return nil, err
	}

	ethTx, err := u.EtherscanService.AddressTransactions(gmAddress)
	if err != nil {
		return nil, err
	}

	totalEth := utils.GetValue(ethBL.Result, 18)
	if totalEth > 0 {
		usdtValue := utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate)
		counting := 0
		for _, item := range ethTx.Result {
			if strings.ToLower(item.From) != strings.ToLower(tcAddress) {
				continue
			}
			item.From = tcAddress
			item.To = gmAddress
			itemTotalEth := utils.GetValue(item.Value, 18)
			item.UsdtValue = utils.ToUSDT(fmt.Sprintf("%f", itemTotalEth), ethRate)
			counting++
		}
		if counting == 0 {
			return nil, errors.New("not balance")
		}

		resp := &structure.AnalyticsProjectDeposit{}
		resp.CurrencyRate = ethRate
		resp.Currency = string(entity.ETH)
		resp.Value = ethBL.Result
		resp.UsdtValue = usdtValue
		resp.Items = ethTx.Result

		return resp, nil
	}
	return nil, errors.New("not balance")
}

func (u Usecase) GetChartDataBTCForGMCollection(tcWallet string, gmWallet string) (*structure.AnalyticsProjectDeposit, error) {
	btcRate, err := helpers.GetExternalPrice("btc")
	_ = btcRate
	_ = err

	resp := &structure.AnalyticsProjectDeposit{}
	return resp, nil
}

func (u Usecase) GetChartDataForGMCollection() (*structure.AnalyticsProjectDeposit, error) {
	ethDataChan := make(chan structure.AnalyticsProjectDepositChan)
	btcDataChan := make(chan structure.AnalyticsProjectDepositChan)

	go func(ethDataChan chan structure.AnalyticsProjectDepositChan) {
		data := &structure.AnalyticsProjectDeposit{}
		var err error
		defer func() {
			ethDataChan <- structure.AnalyticsProjectDepositChan{
				Value: data,
				Err:   err,
			}
		}()
		wallets, err := u.Repo.FindNewCityGmByType("eth")
		if err != nil {
			for _, wallet := range wallets {
				temp, err := u.GetChartDataEthForGMCollection(wallet.UserAddress, wallet.Address)
				if err != nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
				}
			}
		}

		// for old
		gmAddress := os.Getenv("GM_ETH_ADDERSS")
		if gmAddress == "" {
			gmAddress = "0x360382fa386dB659a96557A2c7F9Ce7195de024E"
		}
		fromWallets := []string{
			"0xD78D4be39B0C174dF23e1941aC7BA3e8E2a6b3B6",
		}
		for _, wallet := range fromWallets {
			temp, err := u.GetChartDataEthForGMCollection(wallet, gmAddress)
			if err != nil && temp != nil {
				data.Items = append(data.Items, temp.Items...)
				data.UsdtValue += temp.UsdtValue
				data.Value += temp.Value
			}
		}
	}(ethDataChan)

	go func(btcDataChan chan structure.AnalyticsProjectDepositChan) {
		data := &structure.AnalyticsProjectDeposit{}
		var err error
		defer func() {
			btcDataChan <- structure.AnalyticsProjectDepositChan{
				Value: data,
				Err:   err,
			}
		}()
		wallets, err := u.Repo.FindNewCityGmByType("btc")
		if err != nil {
			for _, wallet := range wallets {
				temp, err := u.GetChartDataBTCForGMCollection(wallet.UserAddress, wallet.Address)
				if err != nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
				}
			}
		}
	}(btcDataChan)

	ethDataFromChan := <-ethDataChan
	btcDataFromChan := <-btcDataChan

	result := &structure.AnalyticsProjectDeposit{}
	if ethDataFromChan.Err == nil {
		result.Items = append(result.Items, ethDataFromChan.Value.Items...)
		result.UsdtValue += ethDataFromChan.Value.UsdtValue
	}

	if btcDataFromChan.Err == nil {
		result.Items = append(result.Items, btcDataFromChan.Value.Items...)
		result.UsdtValue += btcDataFromChan.Value.UsdtValue
	}

	if len(result.Items) > 0 {
		for _, item := range result.Items {
			_, ok := result.MapItems[item.From]
			if !ok {
				result.MapItems[item.From] = &etherscan.AddressTxItemResponse{
					From:      item.From,
					UsdtValue: item.UsdtValue,
				}
			} else {
				result.MapItems[item.From].UsdtValue += item.UsdtValue
			}
		}
		result.Items = []*etherscan.AddressTxItemResponse{}
		for _, item := range result.MapItems {
			result.Items = append(result.Items, item)
		}
		usdtExtra := 0.0
		for _, item := range result.Items {
			item.ExtraPercent = u.GetExtraPercent(item.From)
			item.UsdtValueExtra = item.UsdtValue/100*item.ExtraPercent + item.UsdtValue
			usdtExtra += item.UsdtValueExtra

		}
		for _, item := range result.Items {
			item.Percent = item.UsdtValueExtra / usdtExtra
			item.GMReceive = item.Percent * 80000
		}
	}

	return result, nil
}

func (u Usecase) GetExtraPercent(address string) float64 {
	user, err := u.Repo.FindUserByWalletAddress(address)
	if err == nil && user.UUID != "" {
		return 30.0
	}

	tcBalance, err := u.TcClientPublicNode.GetBalance(context.TODO(), address)
	if err == nil && tcBalance.Cmp(big.NewInt(0)) > 0 {
		return 20.0
	}

	allow, err := u.Repo.GetProjectAllowList("999998", address)
	if err == nil && allow.UUID != "" {
		return 10.0
	}

	return 0.0
}
