package usecase

import (
	"fmt"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
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

	//gmAddress := os.Getenv("GM_ETH_ADDERSS")
	//ethBL, err := u.EtherscanService.AddressBalance(gmAddress)
	//if err != nil {
	//	return nil, err
	//}

	ethTx, err := u.EtherscanService.AddressTransactions(gmAddress)
	if err != nil {
		return nil, err
	}

	//totalEth := utils.GetValue(ethBL.Result, 18)
	//usdtValue := utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate)

	for _, item := range ethTx.Result {
		item.From = tcAddress
		itemTotalEth := utils.GetValue(item.Value, 18)
		itemUsdtValue := utils.ToUSDT(fmt.Sprintf("%f", itemTotalEth), ethRate)
		item.UsdtValue = itemUsdtValue
		item.ExtraPercent = 0
		// TODO item.UsdtValueExtra = item.UsdtValue/100*item.ExtraPercent + item.UsdtValue
		// TODO percent := itemUsdtValue / usdtValue
		// TODO item.Percent = float64(percent)
	}

	resp := &structure.AnalyticsProjectDeposit{}
	//resp.CurrencyRate = ethRate
	//resp.Currency = string(entity.ETH)
	//resp.Value = ethBL.Result
	//resp.UsdtValue = usdtValue
	resp.Items = ethTx.Result

	return resp, nil
}

func (u Usecase) GetChartDataBTCForGMCollection() (*structure.AnalyticsProjectDeposit, error) {
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
		// TODO call database to get list and loop

		data, err = u.GetChartDataEthForGMCollection("", "")
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

		data, err = u.GetChartDataBTCForGMCollection()
	}(btcDataChan)

	dataFromChan := <-ethDataChan
	if dataFromChan.Err != nil {
		return nil, dataFromChan.Err
	}

	btcDataFromChan := <-btcDataChan
	if btcDataFromChan.Err != nil {
		return nil, btcDataFromChan.Err
	}

	return dataFromChan.Value, nil
}
