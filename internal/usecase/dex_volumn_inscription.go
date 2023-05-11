package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"math/big"
	"os"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"strconv"
	"strings"
	"time"
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

func (u Usecase) GetChartDataEthForGMCollection(tcAddress string, gmAddress string, oldData bool) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.eth2.gmAddress." + tcAddress + "." + gmAddress)
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			return result, nil
		}
	}

	// try from cache
	keyRate := fmt.Sprintf("gm-collections.deposit.eth.rate")
	var ethRate float64
	cachedETHRate, err := u.Cache.GetData(keyRate)
	if err == nil {
		ethRate, _ = strconv.ParseFloat(*cachedETHRate, 64)
	}
	if ethRate == 0 {
		ethRate, err = helpers.GetExternalPrice(string(entity.ETH))
		if err != nil {
			return nil, err
		}
		u.Cache.SetDataWithExpireTime(keyRate, ethRate, 60*60) // cache by 1 hour
	}

	ethBL, err := u.EtherscanService.AddressBalance(gmAddress)
	time.Sleep(time.Millisecond * 100)
	if err != nil {
		return nil, err
	}

	totalEth := utils.GetValue(ethBL.Result, 18)
	if totalEth > 0 {
		usdtValue := utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate)

		var items []*etherscan.AddressTxItemResponse
		if oldData {
			// get tx by addr
			ethTx, err := u.EtherscanService.AddressTransactions(gmAddress)
			time.Sleep(time.Millisecond * 100)
			if err != nil {
				return nil, err
			}
			counting := 0
			for _, item := range ethTx.Result {
				if oldData {
					if strings.ToLower(item.From) != strings.ToLower(tcAddress) {
						continue
					}
				}
				items = append(items, &etherscan.AddressTxItemResponse{
					From:      tcAddress,
					To:        gmAddress,
					Value:     item.Value,
					UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(item.Value, 18)), ethRate),
					Currency:  string(entity.ETH),
				})
				counting++
			}
			if counting == 0 {
				return nil, errors.New("not balance - " + gmAddress)
			}
		} else {
			items = append(items, &etherscan.AddressTxItemResponse{
				From:      tcAddress,
				To:        gmAddress,
				Value:     ethBL.Result,
				UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(ethBL.Result, 18)), ethRate),
				Currency:  string(entity.ETH),
			})
		}

		resp := &structure.AnalyticsProjectDeposit{}
		resp.CurrencyRate = ethRate
		resp.Currency = string(entity.ETH)
		resp.Value = ethBL.Result
		resp.UsdtValue = usdtValue
		resp.Items = items

		u.Cache.SetDataWithExpireTime(key, resp, 24*60*60) // cache by 1 day
		return resp, nil
	}
	return nil, errors.New("not balance - " + gmAddress)
}

func (u Usecase) GetChartDataBTCForGMCollection(tcWallet string, gmWallet string, oldData bool) (*structure.AnalyticsProjectDeposit, error) {
	return nil, errors.New("rate limit")
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.btc1.gmAddress." + tcWallet + "." + gmWallet)
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			return result, nil
		}
	}

	// try from cache
	keyRate := fmt.Sprintf("gm-collections.deposit.btc.rate")
	var btcRate float64
	cachedETHRate, err := u.Cache.GetData(keyRate)
	if err == nil {
		btcRate, _ = strconv.ParseFloat(*cachedETHRate, 64)
	}
	if btcRate == 0 {
		btcRate, err := helpers.GetExternalPrice(string(entity.BIT))
		if err != nil {
			return nil, err
		}
		u.Cache.SetDataWithExpireTime(keyRate, btcRate, 60*60) // cache by 1 hour
	}

	resp, err := u.MempoolService.AddressTransactions(gmWallet)
	time.Sleep(time.Millisecond * 500)
	if err != nil {
		return nil, err
	}

	vouts := []mempool_space.AddressTxItemResponseVout{}
	for _, item := range resp {
		if item.Status.Confirmed {
			if oldData {
				isContinue := true
				for _, v := range item.Vin {
					if strings.ToLower(v.Prevout.Scriptpubkey_address) == strings.ToLower(tcWallet) {
						isContinue = false
					}
				}
				if isContinue {
					continue
				}
			}
			vs := item.Vout
			for _, v := range vs {
				if strings.ToLower(v.ScriptpubkeyAddress) == strings.ToLower(gmWallet) {
					vouts = append(vouts, v)
				}
			}
		}
	}

	analyticItems := []*etherscan.AddressTxItemResponse{}
	total := int64(0)
	for _, vout := range vouts {
		value := fmt.Sprintf("%d", vout.Value)
		analyticItem := &etherscan.AddressTxItemResponse{
			From:     tcWallet,
			To:       vout.ScriptpubkeyAddress,
			Value:    value,
			Currency: string(entity.BIT),
		}

		itemTotalEth := utils.GetValue(value, 8)
		itemUsdtValue := utils.ToUSDT(fmt.Sprintf("%f", itemTotalEth), btcRate)
		analyticItem.UsdtValue = itemUsdtValue

		total += vout.Value
		analyticItems = append(analyticItems, analyticItem)
	}

	amount := fmt.Sprintf("%d", total)

	amountF := utils.GetValue(amount, float64(8))
	usdt := utils.ToUSDT(fmt.Sprintf("%f", amountF), btcRate)

	resp1 := &structure.AnalyticsProjectDeposit{
		Value:        fmt.Sprintf("%d", total),
		Currency:     string(entity.BIT),
		CurrencyRate: btcRate,
		UsdtValue:    usdt,
		Items:        analyticItems,
	}
	u.Cache.SetDataWithExpireTime(key, resp1, 24*60*60) // cache by 1 day
	return resp1, nil
}

func (u Usecase) GetChartDataForGMCollection(useCaching bool) (*structure.AnalyticsProjectDeposit, error) {
	key := fmt.Sprintf("gm-collections.deposit")
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if err != nil {
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
			wallets, err := u.Repo.FindNewCityGmByType(string(entity.ETH))
			if err == nil {
				for _, wallet := range wallets {
					temp, err := u.GetChartDataEthForGMCollection(wallet.UserAddress, wallet.Address, false)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny("GetChartDataEthForGMCollection", zap.Any("err", err))
					}
				}
			}
			err = nil

			// for old
			gmAddress := os.Getenv("GM_ETH_ADDERSS")
			if gmAddress == "" {
				gmAddress = "0x360382fa386dB659a96557A2c7F9Ce7195de024E"
			}
			fromWallets := []string{
				"0x2c7aFd015A4080C835139E94D0f624bE552b9c66",
				"0x46Ad79eFd29B4212eE2dB32153c682Db06614Ce5",
				"0xD78D4be39B0C174dF23e1941aC7BA3e8E2a6b3B6",
				"0xBFB9AC25EBC9105c2e061E7640B167c6150A7325",
				"0xa3017BB12fe3C0591e5C93011e988CA4b45aa1B4",
				"0xa3EEE445D4DFBBc0C2f4938CB396a59c7E0dE526",
				"0xEAcDD6b4B80Fcb241A4cfAb7f46e886F19c89340",
				"0x7729A5Cfe2b008B7B19525a10420E6f53941D2a4",
				"0x4bF946271EEf390AC8c864A01F0D69bF3b858569",
				"0x21668e3B9f5Aa2a3923E22AA96a255fE8d3b9aac",
				"0x597c32011116c94994619Cf6De15b3Fdc061a983",
				"0xB18278584bD3e41DB25453EE3c7DeDfc84040420",
				"0xfA9A55607BF094f991884f722b7Fba3A76687e40",
				"0xCa2b4ad56a82bc7F8c5A01184A9D9c341213e0d3",
				//"0xfA9A55607BF094f991884f722b7Fba3A76687e40",
				"0x63cBF2D7cf7EF30b9445bEAB92997FF27A0bcc70",
				"0x64BE8226638fdF2f85D8E3A01F849E0c47AE9446",
				"0xbf22409c832E944CeF2B33d9929b8905163Ae5d4",
				"0xda9979247dC98023C0Ff6A59BC7C91bB627d4934",
				"0x9c0Da3467AeD02e49Fe051104eFb2255C2982C61",
				"0xCd2b27C0dc8db90398dB92198a603e5D5D0d5e30",
				"0xe9084DEDfcD06E63Dc980De1464f7786e2690c82",
			}
			for _, wallet := range fromWallets {
				temp, err := u.GetChartDataEthForGMCollection(wallet, gmAddress, true)
				if err == nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
					data.CurrencyRate = temp.CurrencyRate
				}
				if err != nil {
					u.Logger.ErrorAny("GetChartDataEthForGMCollection", zap.Any("err", err))
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
			wallets, err := u.Repo.FindNewCityGmByType(string(entity.BIT))
			if err == nil {
				for _, wallet := range wallets {
					temp, err := u.GetChartDataBTCForGMCollection(wallet.UserAddress, wallet.Address, false)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny("GetChartDataBTCForGMCollection", zap.Any("err", err))
					}
				}
			}

			// for old data
			gmAddress := os.Getenv("GM_BTC_ADDRESS")
			if gmAddress == "" {
				gmAddress = "bc1pqkvfsyxd8fw0e985wlts5kkz8lxgs62xgx8zsfyhaqr2qq3t2ttq28dfta"
			}
			fromWallets := []string{
				"bc1pcry79t9fe9vcc8zeernn9k2yh8k95twc2yk5fcs5d4g8myly6wwst3r6xa",
				"bc1qyczv69fgcxtkpwa6c7k3aaveqjvmr0gzltlhnz",
				"bc1plurxvkzyg4vmp0qn9u0rx4xmhymjtqh0kan3gydmrrq2djdq5y0spr8894",
				"bc1pft0ks6263303ycl93m74uxurk7jdz6dnsscz22yf74z4qku47lus38haz2",
				"bc1q0whajwm89z822pqfe097z7yyay6rfvmhsagx56",
			}

			for _, wallet := range fromWallets {
				temp, err := u.GetChartDataBTCForGMCollection(wallet, gmAddress, true)
				if err == nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
					data.CurrencyRate = temp.CurrencyRate
				}
				if err != nil {
					u.Logger.ErrorAny("GetChartDataBTCForGMCollection", zap.Any("err", err))
				}
			}

		}(btcDataChan)

		ethDataFromChan := <-ethDataChan
		btcDataFromChan := <-btcDataChan

		result := &structure.AnalyticsProjectDeposit{}
		if ethDataFromChan.Value != nil && len(ethDataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, ethDataFromChan.Value.Items...)
			result.UsdtValue += ethDataFromChan.Value.UsdtValue
		}

		if btcDataFromChan.Value != nil && len(btcDataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, btcDataFromChan.Value.Items...)
			result.UsdtValue += btcDataFromChan.Value.UsdtValue
		}

		if len(result.Items) > 0 {
			result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
			for _, item := range result.Items {
				_, ok := result.MapItems[item.From]
				if !ok {
					result.MapItems[item.From] = &etherscan.AddressTxItemResponse{
						From:      item.From,
						To:        item.To,
						UsdtValue: item.UsdtValue,
						Currency:  item.Currency,
						Value:     item.Value,
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
			usdtValue := 0.0
			for _, item := range result.Items {
				item.ExtraPercent = u.GetExtraPercent(item.From)
				item.UsdtValueExtra = item.UsdtValue/100*item.ExtraPercent + item.UsdtValue
				usdtExtra += item.UsdtValueExtra
				usdtValue += item.UsdtValue
			}
			for _, item := range result.Items {
				item.Percent = item.UsdtValueExtra / usdtExtra * 100
				item.GMReceive = item.Percent * 8000 / 100
			}
			result.UsdtValue = usdtValue
		}

		u.Cache.SetDataWithExpireTime(key, result, 1800)
		return result, nil
	}

	err = json.Unmarshal([]byte(*cached), result)
	if err != nil {
		return nil, err
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
