package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	"math/big"
	"os"
	"rederinghub.io/external/coin_market_cap"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"sort"
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

func (u Usecase) GetChartDataERC20ForGMCollection(tcAddress string, gmAddress string, transferedETH []string, ens string, avatar string) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.erc20_1.gmAddress." + tcAddress + "." + gmAddress)
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection", zap.Error(err), zap.String("gmAddress", gmAddress))
			return result, nil
		}
	}

	keypepeRate := fmt.Sprintf("gm-collections.deposit.pepeRate.rate")
	var pepeRate float64
	cachedPEPERate, err := u.Cache.GetData(keypepeRate)
	if err == nil {
		pepeRate, _ = strconv.ParseFloat(*cachedPEPERate, 64)
	}
	if pepeRate == 0 {
		pRate, err := u.CoinMarketCap.PriceConversion(24478) //PEPE ID
		if err == nil && pRate != nil {
			pepeRate = pRate.Data.Quote.USD.Price
		}

		u.Cache.SetDataWithExpireTime(keypepeRate, pepeRate, 60*60) // cache by 1 hour
	}

	keyturboRate := fmt.Sprintf("gm-collections.deposit.turboRate.rate")
	var turboRate float64 = 0
	cachedTURBORate, err := u.Cache.GetData(keyturboRate)
	if err == nil {
		turboRate, _ = strconv.ParseFloat(*cachedTURBORate, 64)
	}
	if turboRate == 0 {
		tRate, err := u.CoinMarketCap.PriceConversion(24911) //TURBO ID
		if err == nil && tRate != nil {
			turboRate = tRate.Data.Quote.USD.Price
		}
		u.Cache.SetDataWithExpireTime(keyturboRate, turboRate, 60*60) // cache by 1 hour
	}

	pepe := "0x6982508145454ce325ddbe47a25d4ec3d2311933"
	turbo := "0xa35923162c49cf95e6bf26623385eb431ad920d3"
	moralisERC20BL, err := u.MoralisNft.TokenBalanceByWalletAddress(gmAddress, []string{pepe, turbo})
	//time.Sleep(time.Millisecond * 250)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection err1111", zap.Error(err), zap.String("gmAddress", gmAddress))
		return nil, err
	}

	pepeBalance := moralisERC20BL[pepe]
	turboBalance := moralisERC20BL[turbo]

	var items []*etherscan.AddressTxItemResponse
	usdtValue := float64(0)
	// pepe
	totalPepe := utils.GetValue(pepeBalance.Balance, 18)
	if totalPepe > 0 {
		usdtValue += utils.ToUSDT(fmt.Sprintf("%f", totalPepe), pepeRate)
		transferUsdtValue := float64(0)
		items = append(items, &etherscan.AddressTxItemResponse{
			From:      tcAddress,
			To:        gmAddress,
			Value:     pepeBalance.Balance,
			UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", totalPepe), pepeRate) + transferUsdtValue,
			Currency:  string(entity.PEPE),
			ENS:       ens,
			Avatar:    avatar,
		})
	}
	// Turbo
	totalTurbo := utils.GetValue(turboBalance.Balance, 18)
	if totalTurbo > 0 {
		usdtValue += utils.ToUSDT(fmt.Sprintf("%f", totalTurbo), turboRate)
		transferUsdtValue := float64(0)
		items = append(items, &etherscan.AddressTxItemResponse{
			From:      tcAddress,
			To:        gmAddress,
			Value:     turboBalance.Balance,
			UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", totalTurbo), turboRate) + transferUsdtValue,
			Currency:  string(entity.TURBO),
			ENS:       ens,
			Avatar:    avatar,
		})
	}

	if len(items) > 0 {
		resp := &structure.AnalyticsProjectDeposit{}
		//resp.CurrencyRate = ethRate
		//resp.Value = moralisEthBL.Balance
		resp.Currency = string(entity.ETH)
		resp.UsdtValue = usdtValue
		resp.Items = items
		u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
		return resp, nil
	}
	return nil, errors.New("not balance - " + gmAddress)
}

func (u Usecase) GetChartDataEthForGMCollection(tcAddress string, gmAddress string, transferedETH []string, oldData bool, ens string, avatar string) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.eth3.gmAddress." + tcAddress + "." + gmAddress)
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			logger.AtLog.Logger.Info("GetChartDataEthForGMCollection", zap.Any("result", result), zap.String("gmAddress", gmAddress))
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
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection", zap.Error(err), zap.String("gmAddress", gmAddress))
			return nil, err
		}
		u.Cache.SetDataWithExpireTime(keyRate, ethRate, 60*60) // cache by 1 hour
	}

	moralisEthBL, err := u.MoralisNft.AddressBalance(gmAddress)
	//time.Sleep(time.Millisecond * 250)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err2222", zap.Error(err), zap.String("gmAddress", gmAddress))
		//return nil, err
		moralisEthBL = new(nfts.MoralisBalanceResp)
		temp, err := u.EtherscanService.AddressBalance(gmAddress)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err3333", zap.Error(err), zap.String("gmAddress", gmAddress))
			return nil, err
		}
		moralisEthBL.Balance = temp.Result
	}

	//ethBL, err := u.EtherscanService.AddressBalance(gmAddress)
	//time.Sleep(time.Millisecond * 100)
	//if err != nil {
	//	return nil, err
	//}

	totalEth := utils.GetValue(moralisEthBL.Balance, 18)
	if totalEth > 0 {
		usdtValue := utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate)

		var items []*etherscan.AddressTxItemResponse
		if oldData {
			// get tx by addr
			ethTx, err := u.EtherscanService.AddressTransactions(gmAddress)
			time.Sleep(time.Millisecond * 100)
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataEthForGMCollection", zap.Error(err), zap.String("gmAddress", gmAddress))
				return nil, err
			}
			counting := 0
			for _, item := range ethTx.Result {
				if strings.ToLower(item.From) != strings.ToLower(tcAddress) {
					continue
				}
				items = append(items, &etherscan.AddressTxItemResponse{
					From:      tcAddress,
					To:        gmAddress,
					Value:     item.Value,
					UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(item.Value, 18)), ethRate),
					Currency:  string(entity.ETH),
					ENS:       ens,
					Avatar:    avatar,
				})
				counting++
			}
			if counting == 0 {
				return nil, errors.New("not balance - " + gmAddress)
			}
		} else {
			transferUsdtValue := float64(0)
			transferEthValue := float64(0)
			if len(transferedETH) > 0 {
				for _, v := range transferedETH {
					temp := utils.GetValue(v, 18)
					transferEthValue += temp
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", temp), ethRate)
				}
			}
			items = append(items, &etherscan.AddressTxItemResponse{
				From:      tcAddress,
				To:        gmAddress,
				Value:     fmt.Sprintf("%f", transferEthValue+totalEth),
				UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate) + transferUsdtValue,
				Currency:  string(entity.ETH),
				ENS:       ens,
				Avatar:    avatar,
			})
		}

		resp := &structure.AnalyticsProjectDeposit{}
		resp.CurrencyRate = ethRate
		resp.Currency = string(entity.ETH)
		resp.Value = moralisEthBL.Balance
		resp.UsdtValue = usdtValue
		resp.Items = items

		cachedExpTime := 1 * 60 * 60 // cache by 1 hour

		if oldData {
			cachedExpTime = 30 * 24 * 60 * 60 //a month
		}
		u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime)
		return resp, nil
	} else {
		transferUsdtValue := float64(0)
		if len(transferedETH) > 0 && !oldData {
			transferEthValue := float64(0)
			for _, v := range transferedETH {
				temp := utils.GetValue(v, 18)
				transferEthValue += temp
				transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", temp), ethRate)
			}
			var items []*etherscan.AddressTxItemResponse
			items = append(items, &etherscan.AddressTxItemResponse{
				From:      tcAddress,
				To:        gmAddress,
				Value:     fmt.Sprintf("%f", transferEthValue),
				UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate) + transferUsdtValue,
				Currency:  string(entity.ETH),
				ENS:       ens,
				Avatar:    avatar,
			})
			resp := &structure.AnalyticsProjectDeposit{}
			resp.CurrencyRate = ethRate
			resp.Currency = string(entity.ETH)
			resp.Value = moralisEthBL.Balance
			resp.UsdtValue = items[0].UsdtValue
			resp.Items = items

			cachedExpTime := 1 * 60 * 60 // cache by 1 hour

			if oldData {
				cachedExpTime = 30 * 24 * 60 * 60 //a month
			}
			u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime)
			return resp, nil
		}
	}
	return nil, errors.New("not balance - " + gmAddress)
}

func (u Usecase) GetChartDataBTCForGMCollection(tcWallet string, gmWallet string, transferedBTC []string, oldData bool) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.btc4.gmAddress." + tcWallet + "." + gmWallet)
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

	analyticItems := []*etherscan.AddressTxItemResponse{}
	if oldData {
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

		total := int64(0)
		for _, vout := range vouts {
			analyticItem := &etherscan.AddressTxItemResponse{
				From:      tcWallet,
				To:        vout.ScriptpubkeyAddress,
				Value:     fmt.Sprintf("%d", vout.Value),
				Currency:  string(entity.BIT),
				UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(fmt.Sprintf("%d", vout.Value), 8)), btcRate),
			}

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
		u.Cache.SetDataWithExpireTime(key, resp1, 24*60*60*30) // cache by a month
		return resp1, nil
	} else {
		/*_, bs, err := u.buildBTCClient()
		if err != nil {
			return nil, err
		}
		balance, confirm, err := bs.GetBalance(gmWallet)*/
		walletInfo, err := btc.GetBalanceFromQuickNode(gmWallet, u.Config.QuicknodeAPI)
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Millisecond * 50)
		if err != nil {
			return nil, err
		}
		if walletInfo.Balance > 0 {
			transferUsdtValue := float64(0)
			transferBtcValue := float64(0)
			if len(transferedBTC) > 0 {
				for _, v := range transferedBTC {
					temp := utils.GetValue(v, 8)
					transferBtcValue += temp
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", temp), btcRate)
				}
			}

			temp, _ := strconv.ParseFloat(fmt.Sprintf("%d", walletInfo.Balance), 64)
			item := &etherscan.AddressTxItemResponse{
				From:      tcWallet,
				To:        gmWallet,
				Value:     fmt.Sprintf("%f", temp+transferBtcValue),
				Currency:  string(entity.BIT),
				UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(fmt.Sprintf("%d", walletInfo.Balance), 8)), btcRate) + transferUsdtValue,
			}
			analyticItems = append(analyticItems, item)
			resp1 := &structure.AnalyticsProjectDeposit{
				Value:        fmt.Sprintf("%d", walletInfo.Balance),
				Currency:     string(entity.BIT),
				CurrencyRate: btcRate,
				UsdtValue:    item.UsdtValue,
				Items:        analyticItems,
			}
			u.Cache.SetDataWithExpireTime(key, resp1, 2*60*60) // cache by 2 hours
			return resp1, nil
		} else {
			transferUsdtValue := float64(0)
			transferBtcValue := float64(0)
			if len(transferedBTC) > 0 {
				for _, v := range transferedBTC {
					temp := utils.GetValue(v, 8)
					transferBtcValue += temp
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", temp), btcRate)
				}
				item := &etherscan.AddressTxItemResponse{
					From:      tcWallet,
					To:        gmWallet,
					Value:     fmt.Sprintf("%f", transferBtcValue),
					Currency:  string(entity.BIT),
					UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(fmt.Sprintf("%d", walletInfo.Balance), 8)), btcRate) + transferUsdtValue,
				}
				analyticItems = append(analyticItems, item)
				resp1 := &structure.AnalyticsProjectDeposit{
					Value:        fmt.Sprintf("%d", walletInfo.Balance),
					Currency:     string(entity.BIT),
					CurrencyRate: btcRate,
					UsdtValue:    item.UsdtValue,
					Items:        analyticItems,
				}
				u.Cache.SetDataWithExpireTime(key, resp1, 2*60*60) // cache by 6 hours
				return resp1, nil
			}
		}
		return nil, errors.New("not balance - " + gmWallet)
	}
}

func (u Usecase) JobGetChartDataForGMCollection() error {
	//clear cache for top 10 items
	u.ClearCacheTop10GMDashboard()

	//start
	now := time.Now().UTC()
	preText := fmt.Sprintf("[Analytics][Start] - Get chart data for GM Dashboard")
	content := fmt.Sprintf("Start at: %v", now)
	u.SendGMMEssageToSlack(preText, content)

	data, err := u.GetChartDataForGMCollection(false)
	if err != nil {
		log.Println("JobGetChartDataForGMCollection GetChartDataForGMCollection err", err)
	}

	//end
	end := time.Now().UTC()
	preText = fmt.Sprintf("[Analytics][End] - Get chart data for GM Dashboard")
	content = fmt.Sprintf("End at: %v with USDT: %f, contributors: %d", end, data.UsdtValue, len(data.Items))
	u.SendGMMEssageToSlack(preText, content)
	return err
}

func (u Usecase) GetChartDataForGMCollection(useCaching bool) (*structure.AnalyticsProjectDeposit, error) {
	key := fmt.Sprintf("gm-collections.deposit")
	result := &structure.AnalyticsProjectDeposit{}
	//u.Cache.Delete(key)
	cached, err := u.Cache.GetData(key)
	if !useCaching || err != nil {
		if useCaching {
			return nil, err
		}
		ethDataChan := make(chan structure.AnalyticsProjectDepositChan)
		btcDataChan := make(chan structure.AnalyticsProjectDepositChan)
		erc20DataChan := make(chan structure.AnalyticsProjectDepositChan)

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
					temp, err := u.GetChartDataEthForGMCollection(wallet.UserAddress, wallet.Address, wallet.NativeAmount, false, wallet.ENS, wallet.Avatar)
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
			fromWallets := map[string]string{
				"0x2c7aFd015A4080C835139E94D0f624bE552b9c66": "",
				"0x46Ad79eFd29B4212eE2dB32153c682Db06614Ce5": "wwf88.eth",
				"0xD78D4be39B0C174dF23e1941aC7BA3e8E2a6b3B6": "",
				"0xBFB9AC25EBC9105c2e061E7640B167c6150A7325": "littlered.eth",
				"0xa3017BB12fe3C0591e5C93011e988CA4b45aa1B4": "",
				"0xa3EEE445D4DFBBc0C2f4938CB396a59c7E0dE526": "",
				"0xEAcDD6b4B80Fcb241A4cfAb7f46e886F19c89340": "",
				"0x7729A5Cfe2b008B7B19525a10420E6f53941D2a4": "trappavelli.eth",
				"0x4bF946271EEf390AC8c864A01F0D69bF3b858569": "",
				"0x21668e3B9f5Aa2a3923E22AA96a255fE8d3b9aac": "",
				"0x597c32011116c94994619Cf6De15b3Fdc061a983": "",
				"0xB18278584bD3e41DB25453EE3c7DeDfc84040420": "",
				"0xfA9A55607BF094f991884f722b7Fba3A76687e40": "",
				"0xCa2b4ad56a82bc7F8c5A01184A9D9c341213e0d3": "",
				"0x63cBF2D7cf7EF30b9445bEAB92997FF27A0bcc70": "",
				"0x64BE8226638fdF2f85D8E3A01F849E0c47AE9446": "",
				"0xbf22409c832E944CeF2B33d9929b8905163Ae5d4": "",
				"0xda9979247dC98023C0Ff6A59BC7C91bB627d4934": "",
				"0x9c0Da3467AeD02e49Fe051104eFb2255C2982C61": "",
				"0xCd2b27C0dc8db90398dB92198a603e5D5D0d5e30": "",
				"0xe9084DEDfcD06E63Dc980De1464f7786e2690c82": "",
			}
			for wallet, ens := range fromWallets {
				temp, err := u.GetChartDataEthForGMCollection(strings.ToLower(wallet), gmAddress, []string{}, true, ens, "")
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
					temp, err := u.GetChartDataBTCForGMCollection(wallet.UserAddress, wallet.Address, wallet.NativeAmount, false)
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
				temp, err := u.GetChartDataBTCForGMCollection(wallet, gmAddress, []string{}, true)
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

		go func(erc20DataChan chan structure.AnalyticsProjectDepositChan) {
			data := &structure.AnalyticsProjectDeposit{}
			var err error
			defer func() {
				erc20DataChan <- structure.AnalyticsProjectDepositChan{
					Value: data,
					Err:   err,
				}
			}()
			wallets, err := u.Repo.FindNewCityGmByType(string(entity.ETH))
			if err == nil {
				for _, wallet := range wallets {
					temp, err := u.GetChartDataERC20ForGMCollection(wallet.UserAddress, wallet.Address, wallet.NativeAmount, wallet.ENS, wallet.Avatar)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny("GetChartDataERC20ForGMCollection", zap.Any("err", err))
					}
				}
			}
		}(erc20DataChan)

		ethDataFromChan := <-ethDataChan
		btcDataFromChan := <-btcDataChan
		erc20DataFromChan := <-erc20DataChan

		result := &structure.AnalyticsProjectDeposit{}
		if ethDataFromChan.Value != nil && len(ethDataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, ethDataFromChan.Value.Items...)
			result.UsdtValue += ethDataFromChan.Value.UsdtValue
		}

		if btcDataFromChan.Value != nil && len(btcDataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, btcDataFromChan.Value.Items...)
			result.UsdtValue += btcDataFromChan.Value.UsdtValue
		}

		if erc20DataFromChan.Value != nil && len(erc20DataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, erc20DataFromChan.Value.Items...)
			result.UsdtValue += erc20DataFromChan.Value.UsdtValue
		}

		if len(result.Items) > 0 {
			result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
			result.MapTokensDeposit = make(map[string][]structure.TokensDeposit)
			for _, item := range result.Items {
				item.From = strings.ToLower(item.From)
				_, ok := result.MapItems[item.From]
				if !ok {
					result.MapItems[item.From] = &etherscan.AddressTxItemResponse{
						From:      item.From,
						To:        item.To,
						UsdtValue: item.UsdtValue,
						Currency:  item.Currency,
						Value:     item.Value,
						Avatar:    item.Avatar,
						ENS:       item.ENS,
					}
					result.MapTokensDeposit[item.From] = []structure.TokensDeposit{
						{
							Name:      item.Currency,
							Value:     item.Value,
							UsdtValue: item.UsdtValue,
						},
					}
				} else {
					result.MapItems[item.From].UsdtValue += item.UsdtValue
					if item.Avatar != "" {
						result.MapItems[item.From].Avatar = item.Avatar
					}
					if item.ENS != "" {
						result.MapItems[item.From].ENS = item.ENS
					}
					result.MapTokensDeposit[item.From] = append(result.MapTokensDeposit[item.From], structure.TokensDeposit{
						Name:      item.Currency,
						Value:     item.Value,
						UsdtValue: item.UsdtValue,
					})
				}
			}
			result.Items = []*etherscan.AddressTxItemResponse{}
			for _, item := range result.MapItems {
				result.Items = append(result.Items, item)
			}
			usdtExtra := 0.0
			usdtValue := 0.0
			for _, item := range result.Items {
				item.ExtraPercent = 0.0 //TODO u.GetExtraPercent(item.From)
				//item.UsdtValueExtra = item.UsdtValue/100*item.ExtraPercent + item.UsdtValue // TODO
				item.UsdtValueExtra = item.UsdtValue
				usdtExtra += item.UsdtValueExtra
				usdtValue += item.UsdtValue
			}
			for _, item := range result.Items {
				item.Percent = item.UsdtValueExtra / usdtExtra * 100
				item.GMReceive = item.Percent * 8000 / 100
			}
			result.UsdtValue = usdtValue
		}

		cachedData := &structure.AnalyticsProjectDeposit{}
		err := json.Unmarshal([]byte(*cached), cachedData)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection json.Unmarshal.cachedData", zap.Error(err))
			return nil, err
		}

		//the new data must be greater than the cached data (old)
		if result.UsdtValue >= cachedData.UsdtValue {
			u.Cache.SetDataWithExpireTime(key, result, 60*60*24*3)

			//backup to DB
			go u.BackupGMDashboardCachedData()
		}

		return result, nil
	}

	err = json.Unmarshal([]byte(*cached), result)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection json.Unmarshal.cachedData", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (u Usecase) GetExtraPercent(address string) float64 {
	user, err := u.Repo.FindUserByWalletAddress(address)
	if err == nil && user.UUID != "" {
		return 30.0
	}

	kll := map[string]bool{
		"0xe96943DA5E2C74c0450612D5feB3537d036CEFC0": true,
		"0xF3A66C660fa1A41f8FcC04504B506163c119552C": true,
		"0x0C2A603D6432D6beb49FCF34fD558eAF4cE68eC8": true,
		"0xADD8DB82416CC78411E5Fd18725eBDaE6DB14cbC": true,
		"0xf1aFd190f76E415B4716c899ECfFC4B80020B788": true,
		"0xfA9A55607BF094f991884f722b7Fba3A76687e40": true,
		"0x917Aa1FE4eE8154CE8c4EeB8f4768BA615245799": true,
		"0x508586efE86F3a41f9e692059A8172a0Fa120619": true,
		"0xB18278584bD3e41DB25453EE3c7DeDfc84040420": true,
		"0x00417C5d13292f5B07c3Dc7Fa82E430B002FE58e": true,
		"0x3fb65feeab83bf60b0d1ffbc4217d2d97a35c8d4": true,
		"0x5C30Bd517A981e38a802E26f12A17528303Cc0b7": true,
		"0x60b97528be1d072aC5b2405b6e292B43B3eb4AD0": true,
		"0x705562434ea1288BdbA5dC6A995DD07B3723dbA6": true,
		"0xE7F3443B16648A89dd3eAf9D179CB4090C02E66F": true,
		"0x75bC526d9655a93c89a459185021b194a71B070C": true,
		"0x5Bca4075dFC8065235cF75C6b15B410e62845Fec": true,
		"0xA4d922Fd0e7763623825Ddb5176724D6367560eE": true,
		"0x3cc65A13cEB912405C8cb7Bb0aC854dA25A6d540": true,
		"0x33074e245A66A2b70caDe6D84927a68F50c6d865": true,
		"0x98261DA497f986b8EE2bD038a1bbBb629578062E": true,
		"0xfb1904dB9aEbc93738f0C6Dc03DCa7376cEEf07e": true,
		"0x1bfc7b12b950A88371947EAFE123126ea3D5314A": true,
		"0x1f4bee6bf2Fe093Fc1b56F7A4C1C3F4dd98FbC0a": true,
		"0xaBd99a2d4FBdFCC22Ba508690B47aca266FaA8ed": true,
		"0x2218Eb0947140D776c5186317B6DbED151069caE": true,
		"0xC62ef665a6F45D2937B8Bf745E396E4db0297D0E": true,
		"0x45Ab04C54264F485eeF8DB0c20E531E9D37cD53A": true,
		"0x02b4DBF545C9E07691ce14daEe9553C88c80394b": true,
		"0x0D8B49bE1176b7c9436167A4FaA2C0F8547Aa7E7": true,
		"0xF0B131A9EAA9e2c8f1d26D200D47Bc1eDa50FB66": true,
		"0x8dBb75c576B71B43eea54398F8606aeC530181dc": true,
	}
	_, ok := kll[address]
	if ok {
		return 25.0
	}

	tcBalance, err := u.TcClientPublicNode.GetBalance(context.TODO(), address)
	if err == nil && tcBalance.Cmp(big.NewInt(0)) > 0 {
		return 20.0
	}

	manual := map[string]bool{
		"0x87907E5ac909454cc16DD503DC03ed8864eB7191": true,
		"0x08ADb99AF2B78AAb79e8bEe60dfd124C24C68926": true,
		"0x26c4C9E2a772Fe5F413063c62E7b2E9e14F5Dc5A": true,
		"0xbbE7148F4e5D7c607845b60C39A21173c0E0a77b": true,
		"0xf6890Ef5C3ABD88130fa80407067AD49E383cf55": true,
		"0x0b4Cf3ac10aF6E5e07c7920Cfc4b01DeA69Aa047": true,
		"0x918453d249A22b6A8535c81e21F7530CD6Ab59F1": true,
		"0xFa6284E0D78e9c3fe9934b548F167f8E82f63c38": true,
		"0x941beced3e87a15ba22e1a3705b547f50cfd2eb1": true,
		"0xF86a588f3604e3b12899f710E3D572F76ffB94B8": true,
		"0xdDF9a1F60aed0118739e317290B874aDAe269327": true,
		"0x73e60cD967E957bC6e074F93320FfA1d52697D5b": true,
		"0x913735E6D76b6b954CA799511244FF430CdA642f": true,
		"0xd0daddf983fce88bef3f10fc12280d0f0cd1208c": true,
		"0x4A29367c5Ae9F84eF03E447D1f7deE8e6b16229D": true,
		"0xf0E87FF63595e6043C31d78d4e16A48ff224880a": true,
		"0x8ce3b15b6f3f32b757f88bdcf69464e5afdc8452": true,
		"0x88b582774a9428226B33a76AA39821EC3899E883": true,
		"0x4bdfa798281e74438399a45451C377A0de6206b1": true,
		"0x93c2dc7c55af662d19d9ef01a14fee3682d7d3bd": true,
		"0x0B442613dcc852B0531c7C23f4fCE48C962472E3": true,
		"0xd0c72ada2bca3ae54e156416b35ee2ab14fd7e09": true,
		"0x9a6ef672c9af8c98201D3DDfFBa9de4A67Bb7Df8": true,
		"0xa3CFc7AF4e310eEaF32F325031Eb0352350f0822": true,
		"0xF8C3A3CfEF10fb2b64D92AA9953923a65ef0b9Bd": true,
		"0x509eCD4cCFc96Bc152156e986eC35Aaa680BD45A": true,
		"0x2389b28518c89c3b65989b6959a16a3940b03446": true,
		"0xf0358a8ccD61A35750c9210b80D3A6078b8c1309": true,
		"0xD616B44BA3B0FBdaF403E5fb2675013B37ecb35f": true,
		"0x0D8B49bE1176b7c9436167A4FaA2C0F8547Aa7E7": true,
		"0x0a90470Ef017c15f18F29C0B33e4935B0eC34929": true,
		"0x6F364eC121020E246993855666a536821D687621": true,
		"0x7f1d1C865A9De2a35ec0269d547D8c943B61a794": true,
		"0x6c8f1DA0f757Ce1a524F671e73Ee8097862b0a2A": true,
		"0xdAfbA80aDF54C6a67FdA035C5eb669b894306819": true,

		"0x360382fa386dB659a96557A2c7F9Ce7195de024E":                         true,
		"0x791Add699b675dF47bF6C350d2A9E75A93Ff8b4A":                         true,
		"0x5DC06223025CA67dDAaCD3d231Ef44A5C213B840":                         true,
		"bc1pw9l2j5tyxd5rj5pyu99afy564gej59pg0m0ucpdnfgj30hnl9zgsr80vuv":     true,
		"0x894486c6933b74bdd6b9e41b91fb460da49da84d":                         true,
		"0x9edBF347Ae013d202ec979bc1D766C54BD3a9879":                         true,
		"bc1pfw8xg22pwkgpf567vs5tk6dtmjmw2ypwv80s4v9wx3yyrqdvxlnqcaxhra":     true,
		"0xb0086CEd87CeD727E01c3b15E54D7Ef04301aC31":                         true,
		"bc1pctm6mc79qtwcepknq5ga9l4p5qfywp79rtfy87g9dvnf97rzhm7sme9a63":     true,
		"0xc6cd4a8412e5103d217f015d37372628e2a6b0a4":                         true,
		"0x8264f75a5f1c1924150414b8173b62DE7db2cD09":                         true,
		"0x1F22043648d5299Daf343ce520ef024456c95b55":                         true,
		"0x4354e8616c489c8b5e5b1972f50e4de2f981b95b":                         true,
		"0x3cf47f877b97c3e88ce467f7ee63c6a45c763ae6":                         true,
		"0xb14bc93be53b3608324c1b48706909842846b71f":                         true,
		"0x303eefb62956735f60d407fa38399c502c7680db":                         true,
		"0xd9b21df68FAD1b99E9f4B12C336E777A06A39213":                         true,
		"0x2acf7141dF952991839a769A887F631Fe404C77E":                         true,
		"0x0eD82aC45ccEBC6B1eAB8d0F44C8869602FD8617":                         true,
		"0xFF157c9a4CD63C8C0c4d40D068A8091Bf55c687B":                         true,
		"0x2a3FbcC6baf5744f5E37600c1aEEA0573Eaa0880":                         true,
		"0xC190380DE177e0f769269d154267418c0D792E35":                         true,
		"0xe7536eDAfc815524eeFF2b24664F47fE4e47b7B2":                         true,
		"0x15ABB22f3AdfE984B2404751B97Fe6CEA54E077f":                         true,
		"0x07a8DE33a717F93a9a30D4F496a69D02c3765835":                         true,
		"0x459832e4c2ec9f26947b9E09AE51B6299D8F8FD3":                         true,
		"0x6e1d7197C6fA06c84110F3B6A774C74d3F7421ca":                         true,
		"0xa8ee6B76b62013AF4626fbD8422e4eD7f7c46a7a":                         true,
		"0x597e4B9498b7ff26AF39463339BFa8FcC7e40A68":                         true,
		"0x0B70Fa903D5fcB0D5138AF5cf455c3d1fC3e32cF":                         true,
		"0x3660c5CA0d70A80D07872023903B913dFb41b048":                         true,
		"0xC6dB0618d627D33F94389304f6B5C2135f2Ef03A":                         true,
		"0xb4600572CA493546570a70AdC0AAB1dd64575Ee8":                         true,
		"bc1p5zhxqhslrjydxy767kag4pvnxkgxcvz9qd53xeyew4nwz0f5svjqxu2ydw":     true,
		"0xc2881096027ac72ffbc4892426afe4dd3be712d7":                         true,
		"0x9f9d43a3f5ff532c776e994ffa6f66532b79e51d":                         true,
		"0xE810E7becDb73e2D147321278A3a6c88Ae590732":                         true,
		"0x2AB33612eCeF16DDD6C78a147BFBb3902AD43913":                         true,
		"0x187375D8c65b7b6Db1DE453D6161C4457e2d9862":                         true,
		"0x736E5DF861C75A3af5cB95B996d3133d9d3E9f71":                         true,
		"0x886148A6Bd2c71Db59Ab3aAD230af9F3254173Ee":                         true,
		"0xb70b00cC8D5c1430F6c5123e73c7107f88635aDb":                         true,
		"0x6E510b192838d0C0F6C2841D589e21bf7Ba8bf6D":                         true,
		"0x7f04c4387423c5460f0a797b79B7De2A4769567A":                         true,
		"0x42761BaAb355d4C31e7fe472682715e3a8b4320f":                         true,
		"0xc35577B56E3ae215149668557820e307530292F5":                         true,
		"0x0A2D6BA3465b582D30C1153c6FD8A39A8894495C":                         true,
		"0x73d90dB79da0eE77C42cD6014180Cb69602a3a22":                         true,
		"0x045454cdde0821544bDAEBeF7F43048F88FEc784":                         true,
		"0x499052095eADbd71296586E58627108577a972ed":                         true,
		"0x8d58079a7C9e5320c8DCbfa2e13aB76D9e6Fb5d2":                         true,
		"0x007DbeD1B4a125c45DF88F3FFa350ff70c94DD9f":                         true,
		"0x7CfB6471D2A9913b7d27D3f3983751eFA028Eea9":                         true,
		"bc1p7rgewkpa5jhlg8wzxzyhpdqq8jl35e9mvrdakmkgfwtgusek07yqtxe3pw":     true,
		"0xbC7DA0f4d3681D33c8bE5Eb3405bE83EB34D697E":                         true,
		"0x1970d3148751978453f2e8Dc61E091a289e5B1c4":                         true,
		"0xC449D58e0B602c6c83474E8625864A598733fE08":                         true,
		"bc1pcpkfq2g08ewx5ya2vc2h2ux6ze3wmaem8l76ppj8gzua7ls5kflsv8pkfh":     true,
		"0x1aE3Aa819d529fAB09Ad82aeA3fC847c907665F0":                         true,
		"0x19748446A67B690EF1DD13Ee61A615E9028BC6E0":                         true,
		"0x9442d0Ee96A3DEf432d795214C4c9c0Af1Bfb8aC":                         true,
		"0x01552052B14Af22dc26F5DCEEC0b53dc4343EC16":                         true,
		"0x201d7bf55D33fdBC3E36C9AaFb17660643bCD1A2":                         true,
		"0x18f8040BF2DF552fde9227Af0F369780Ea6c3B1B":                         true,
		"0x903fC0d00Ab75B25C2F876D1f4242d86d0A19565":                         true,
		"bc1plq0z77muw3rl40d382fux8avpm97s40a7ln5ul7g9heqgr7c3cfqknyehm":     true,
		"0xf4d9f500A71Db552aF6e79D00B32080a90677Fd3":                         true,
		"0x8C278A82b891d8EE8c564A4D33ea12d6a4e6790E":                         true,
		"0x79C109c9832A24c31dc63e33f585Dae18FF1DDc8":                         true,
		"0x6c48e54F07177791Fc5d93Bdbe88F33734B27C88":                         true,
		"0xe1E677689092d2289D91daBb83639fEbC7aa22f5":                         true,
		"0x49F19D90b0DC28eE80d253A5Bfb65393Dc44019F":                         true,
		"bc1pquc9gcd7vlc7gn7rnj35g48fennu4p75j55h3srv8d7mupc8rq3s7p7k77":     true,
		"0x1769e601ee4B1Ff83154CFedc245472c0367aCa3":                         true,
		"bc1p82wse42c2m7xsvzrcxamac7uftdg8fx2qr5agxmr9yx4g542ck2sdektld":     true,
		"0xeC3185Dc088dbb28f5cB540a812Da28587EFd499":                         true,
		"0x2894df8c3c0d654674108613f2fcf0cf4e319b65":                         true,
		"0xcF664C7ef7cef9c792Fc3b0a165b0D83c69bbd29":                         true,
		"0x398Df2aAB9419aD606c16511d5ad7A10CCB57AbF":                         true,
		"0xd59eDC0EE02cE36593be6070D8CD835F81401A63":                         true,
		"0xAb48De856930c018238c226D166Feaed3541Ec7d":                         true,
		"bc1p6y5f96zunz9fgq5z0hy6927gedfe7ljwqhvlx0ey7gwuft2gcwcs38xpq8":     true,
		"bc1pue6mazdujj5puqtall73zv3rgp9s459y2dwz6ym5veetffghh4hqzd80lg":     true,
		"bc1pym9cvlpk2f4hj08xp03746wv769jh4mws9thf9vm9xx9yxs4px3qr354a7":     true,
		"0x27893E2Ba03bBdee0f2b8caA749a0783357Cd918":                         true,
		"0x37E6e346aeb4FCd11D15694351fa61431645Bf83":                         true,
		"0xf8598c16d3539d160A3088A0b278C3bf3CbF11F7":                         true,
		"0x3554df96F9072e1484cE60064b7e6D18C69d2445":                         true,
		"0x244B84eD1150D994774bE96f71cDeca7BB16d4BC":                         true,
		"0xB8976BE387D5BBB853305b597c72129Ed8486339":                         true,
		"0x232C5704f3a3c997aE281048885E5C4b3A712573":                         true,
		"0xfAE772E7498B0cDfC4EA929da72333DCC0Dd2517":                         true,
		"bc1qrxm4gqvkkwmag72eqjvmjlj97luwdrvzz8g4xm":                         true,
		"0xCabB179ca4f9360e4761121A2363a3AF5587B1aA":                         true,
		"0xD30f558fc6261C5B53E35a11001932902d7662F7":                         true,
		"0x44a9EcDb36918466D7Ac318fEd254c227Cb5Be7c":                         true,
		"0xF5f076b24bCFa039F1514BD03E388c047CE1e550":                         true,
		"bc1pjwwnxfsz9g2wf4f5mvv4fndafkd2ntjkn88374w7la89k57gd7ks7gvl4l":     true,
		"0xF12EC4Edd6ce98232e0894aa1aE3c39144846906":                         true,
		"0xb25A9908175e1459f84a2365d3393eF42469B9b0":                         true,
		"0x31289c0C737a21D56540a15D0847b428120b263c":                         true,
		"0x2666f0C8FB58d182f2Dd79475DCA4A07B3724607":                         true,
		"0x7e874922303CF8A978BA7e88D4BE459289c63b55":                         true,
		"0xda9979247dC98023C0Ff6A59BC7C91bB627d4934":                         true,
		"0xD0cFE8dA409D0FD31908E4913bd4547Ed988b178":                         true,
		"0x440E70C989972cD5458C87836D46512d38CF383A":                         true,
		"0xdedb005FF85a036770acB2c3Be0FB74d4FF4E465":                         true,
		"0x12BCF162bcaaB6C6F829dcaA5026D72aF956864c":                         true,
		"0xf88986E9448d0B6397cA40a3A83d0168646d69Cf":                         true,
		"0xdd9DCca486ac3DC82b31Eaa0a40d44d2D739b355":                         true,
		"0x30452B5016747c25f814c9bF9c811db9b2b20A7c":                         true,
		"bc1p3k3rn02wr3n3f6h2t9gnhq8tqms2w0v00y9kzze2lqez65js6xwqqlxyhh":     true,
		"0x9f77094cc3da711b02a5cb8980e513348e0fd516":                         true,
		"0x791ffea64feb5de305d354237940f185723600db":                         true,
		"0x93a692FE5477902C3cf5F6997c1cf59a3712cED4":                         true,
		"0xde2da94f2946f7fC0f1953B2C22AdE5A910aC9A4":                         true,
		"0x708Fa0A05E58D7f66089CbfA6bA010DD492585B9":                         true,
		"0x4159b60ecb28198e54a0ed1c6ea87c81a9788ebd":                         true,
		"0xdbCc6698d4686EE3fba49c2245072460594efE6E":                         true,
		"bc1qk53azuknclstk5pluw3ykvxhhxvgt90dzmkfv5":                         true,
		"0xBaCA135D297D17B584327F7BB436df5D271A2AC1":                         true,
		"0x8b4a8E8A0D0B544904b630c090FA8fc246A07bFB":                         true,
		"0x54793c8f80F642Ed4e200177bf845c61915EBc55":                         true,
		"0x427070e33d630B0476BedBA5d53Db224913eB946":                         true,
		"0xea8f8399a99193dbf348faab3294141dd0f25745":                         true,
		"bc1pf6s9w86tw23nc2k8hdjhy7lrqs2yqc6lz8eez2349nx2dzhscudsgr2tzc":     true,
		"0xbf22409c832E944CeF2B33d9929b8905163Ae5d4":                         true,
		"0xf06277d4F9c6215C208c1C86B99A649df99572Bd":                         true,
		"0x8513FAfE1813B6ec3BBD9fc4BAF5340bda8D670B":                         true,
		"0x01C2e7b1de06DA53Bd0EC82fdB59E5767b8c6dA1":                         true,
		"0xdB3821d232B8508091cf447db77824647F8C2695":                         true,
		"0xd308d28Ef6C988429AcC5213c8F2758C374ba5c7":                         true,
		"0x1f0f43310d87406Ff30a2df20c8cc569B90597D4":                         true,
		"0xdaB2d786cE627d93Ea4B7d7043C016EE4730Aba1":                         true,
		"0x94db5F225A1F6968Cd33c84580c0ADAe52a04eDF":                         true,
		"0x86262f79483ee7Dbc2570EccF5bDd0B19D0FF5C2":                         true,
		"0x6cf3C25E511fF34972B6D10FC9A70891338C59A2":                         true,
		"0x9D4718D87eDac3575858DA9AEE48CC3067769bE3":                         true,
		"0x3370CBe8b6ba8F30c4D3D8c31241eB6E60f14e92":                         true,
		"0x8d234ae9a86C57fe7B3a7a287c0751F6AA90a886":                         true,
		"0x61f9E0C45B5A8f729B0A52C8cb4CC67E23c6ca22":                         true,
		"0x162031B2F953D7BA0a5638b02F910Af1C5279754":                         true,
		"0xbd6fB052CC9e52fe23e6b1A8C4C6cAF4541884c8":                         true,
		"0xcC0Ae4aA8372EE5298B6AAd9D6EaFB4de248a60b":                         true,
		"0xD8AA69542EfbaE18a3BB1627e5EF0c714F888D79":                         true,
		"0xaCdC4CBc410Ba07260B23F46237640e34Dd578a1":                         true,
		"0xf90f4933d08d154ebf51910b29be5abeac539925":                         true,
		"0x0BffF40545a2250c3f11993e7B75dbbcB11e36ac":                         true,
		"0xA3A27ba8943134d9700a641FEe3F763cF2b4eDf3":                         true,
		"bc1px8vsju23f6433m3mx3tmkejh39zlr2qe2z26lal4f43tdt4kyhnsv6chkd":     true,
		"0xd210A01936a901c79254De8fC693f743E5757F8F":                         true,
		"0xeB8CAAFcb2B61Dc354A3f311d5f3341Bc2993EC5":                         true,
		"0xf873BeBDD61AB385D6b24C135BAF36C729CE8824":                         true,
		"0x68492881985afd26b04dcC4963ca5128CCd3DD9D":                         true,
		"bc1p6lg9g9ffpz4a247kxuguhm6ug0cspw3shz3w245c098e0lul35kq6vrxa3":     true,
		"0x02A04748db3572fDfdC0fF6f6E257085EBf9b22E":                         true,
		"0x37Db1629458c7ACd1ECC0b6702AC0C6636341F99":                         true,
		"0x290501163C0cc17817e36F5BcD2D25439a326DAb":                         true,
		"0xCAFd432b7EcAfff352D92fcB81c60380d437E99D":                         true,
		"0xD31538fE582D1511b6DE46E95D358e709e409326":                         true,
		"0x5AfA7d532F90Da0146c27cDAd8c2960183181604":                         true,
		"0x9f8A88d1fd15EDd4408d73e5917E7781090C2616":                         true,
		"0x134F62c8e31Ce0CDDc2532f053c0cB92302d57D9":                         true,
		"0xCbC69338c98887dDf680DA64BF2495855E3a7218":                         true,
		"0xEcE44413e685659eE964757364d70B062505fEef":                         true,
		"0x4537c8cf4707475761aD8D9315Ea5EE50DE26464":                         true,
		"0x0B985EF83B04839dae6Ac4832C54432DbD826bA6":                         true,
		"0xFc6243De43EFA8dEEA2099E5643A41a8E9aF73C5":                         true,
		"0x78807cdF90Bcd3e25b520bD2C11eBB881507e5f0":                         true,
		"0x61642C663f48BC4F8f2b8FA1C5dF9aD5abc7F165":                         true,
		"0x83229468BBF208edD0df4bef04191284A067A90F":                         true,
		"0xec8EEe9BAc5569eE02314cEF96f6017643a1210A":                         true,
		"0xc3771f93fd002fa3fd195fb175507bb701428b85":                         true,
		"0x43105a71c046116ed80ef4ea323a6a476c5545b6":                         true,
		"0x2f12a7cc6c9d99b83974bd603a9d3c328494770b":                         true,
		"0x360399bf1ad25c537f8e348b60dc1356328b4309":                         true,
		"0x64809bdf6622856382c27a7cc73e8c8106d33c33":                         true,
		"0x26f5242d398728d9e542d8fdfb63e64ed8ab8eab":                         true,
		"0x0a52a605927bfaab5d781eb27dbf44534d9601db":                         true,
		"0xa4995f4491ec0f4ef5dfde9169608ac51817e282ed9a24d5563ff6ddfe8a3650": true,
		"0x17476d0Ed31f81d95b5ba8960b2D0b4dE4675e64":                         true,
		"Q": true,
		"0x56D947fA1AE6d6293b307b489540C156AF628F0F":                                                                    true,
		"0xf5596f2dDE208d4E25641db4899e49E948BEE989":                                                                    true,
		"0x26c965fdb7bf614fcf9f5684a0737d4072a79f4c":                                                                    true,
		"0x56cdea7b7ae215ae01a4a9709d9bea639789a80e0f0be9c4ecde8708466a42e1 ":                                           true,
		"0x360399bf1Ad25C537F8e348b60DC1356328B4309":                                                                    true,
		"0x6249A9C03B148aCd773d806A50E1aA983e34bE0F":                                                                    true,
		"0xeCFf34C282c48095bB00aF2CAE746F236f5523a8":                                                                    true,
		"0x4b6f1a108dA670A547005a6b89eFC805AB0f155B":                                                                    true,
		"0x565841c1e403e9Db5976395fBA0e64dbb685eaA6":                                                                    true,
		"0xe47733da657ab715f98ef5a28d13186a61daf9bc9cf954c47886e94bbd77f9e0":                                            true,
		"0xcbf0dbca6bb4b6d4e2be6ec862bc552ef08dc2a841be1f7ce0e114b57398f996":                                            true,
		"0xd7CdB08EaDd3D730540817a348b62D15BF45DF10":                                                                    true,
		"0xbec8177b84ff18e9dc83a077f9e54c63a905265e69f86585d65a97e1e5841bca":                                            true,
		"0xA19265bADD946329A8A8F84f25403E44Ab185aB8":                                                                    true,
		"0xc9a969908738219e5805EC233a0417B47ED245EA":                                                                    true,
		"0xF7Fb8B89AF9B28298C572eF8c07cDf25d1579E38":                                                                    true,
		"0xe4c1acfc81da2e03e0d166c7e95031802c3819c1":                                                                    true,
		"bc1ptjpc2welap7eghqe79udlmv288rsrkmrwcumvc2cfdfl48gv6lfsgff8y9":                                                true,
		"0x5E00473077C622E4Af698D790b5a1670ceC5b8E9":                                                                    true,
		"0x94397d6cf674943B23F5aAa792eB457f40ba7f4C":                                                                    true,
		"bc1p30q5f3xvlzpu5qpjl59hjfs8ku6aznp3tp3mft7t0qrvagjyjppsy07lp5":                                                true,
		"0x5a3DDf88B35E355AC29f0742BF060DB6a37B8BC1":                                                                    true,
		"0x500e241c5CF4D03b8ceC5D8fC35b0107Db199DFd":                                                                    true,
		"0x6c591c83aADeAF04635beB57C442077aFDd54C98":                                                                    true,
		"0xcEE6CB978E2A2d9Fd2F47CD9b1793aFd8654E95A":                                                                    true,
		"0xfb4af5976c9911da0772c5f7a8c299e969b210ec9d5104443b75d164da4e38c5":                                            true,
		"0xCCb500F042A25EA50d4076CEE6f0d6C7fCd488d3":                                                                    true,
		"0x619a24E19F64D8599C3AD7B2A49B3314C27F3F40":                                                                    true,
		"0x48d11882D1a28F1e0cFeF25c1EaDCA0239b1Af09":                                                                    true,
		"0x1Ff2c19161e89eE481399905a4329e65BE890a6d":                                                                    true,
		"0x2fa85613f965205c24f53ee896ef6d66b9fa66da":                                                                    true,
		"0x4a5465C105FB89057333CD94C3af63f9910C2FA8":                                                                    true,
		"0x65053640483BD6efE6083A02BadF5bc62dEde3Db":                                                                    true,
		"0x5Ce306Be65984dB51a290dfF096d1A49B428B985":                                                                    true,
		"0x6186290b28d511bff971631c916244a9fc539cfe":                                                                    true,
		"0x5B150a94e6EcE78eBe7792F35f527b24F5a3D9a0":                                                                    true,
		"0x16C9404Ee5706f595ee11563c7918815549e6d1B":                                                                    true,
		"0x45d2baFe56c85433e0b9f9b50dd124ea3041f223":                                                                    true,
		"0x6A15A0A92fEB8078518c229634dF3688c1d46f69":                                                                    true,
		"0xF39b92a1fff76E7F400867C4c60D5d63d8EA9C4C":                                                                    true,
		"0xe5d009bbcE5a5D7ab9c6c2B1b0A56F9B98297Cff":                                                                    true,
		"0x668E961736454a2444ADB485340cB7F0844DDd3d":                                                                    true,
		"0xeBB4975A3fbEaF7FBBcA67f2D68Ec4Ba955Cdecc":                                                                    true,
		"0x00BcaFF8b59a1249a56772B835D1055973173D6d":                                                                    true,
		"0xe4c0B5AEB5006D93906a88C430C8e3602B1488D8":                                                                    true,
		"0xfaD7C751a252E6Aac03aC817Ee8C4D57c9Fe438b":                                                                    true,
		"0x64BE8226638fdF2f85D8E3A01F849E0c47AE9446":                                                                    true,
		"0x445130702EB79291b597C98386b2e7b6d7bB26B1":                                                                    true,
		"bc1q30cfkxkcsae0qs6nl0qq3sepccd8kgm0ct00yw":                                                                    true,
		"0x2b79C35558C931aCd60352e8aCADDD9aa07FF643":                                                                    true,
		"0xEdcF5D6dDc5C8b64dd6927caB50a5B7FB3e50aBD":                                                                    true,
		"0x1dbec7D882d8d9a3a2D21A5483426D4436E77777":                                                                    true,
		"bc1pd2h4pals4h0nfrnv2afemp8dk93n7tn3qptq4flrd8vernfz9rss5j75l3":                                                true,
		"0xEAcDD6b4B80Fcb241A4cfAb7f46e886F19c89340":                                                                    true,
		"0x2E053E05fD89288aEc116E917311f7f472b1A4A3":                                                                    true,
		"0x30C3C993fBc5a98630477c77c7c60bb8708248c9":                                                                    true,
		"0xcFD68C149030AEe9b153C33a8CB6ABA3765228A8":                                                                    true,
		"0x0d9b6e0f938f070879023054e59d0da681085c11":                                                                    true,
		"0x874CfF564b18565c486559Be102219FCe8641781":                                                                    true,
		"0x0C70B20eEF88E15f601Cc26Ae5dF4D038bb293eA":                                                                    true,
		"bc1p0d6ec62376ywl5l54x475cd8wxtlxzcuepeysa754cv2hejxvkesl0yjku":                                                true,
		"0x8fa899ac9C6b1fcBe157B3B0B626493C0e133Ab4":                                                                    true,
		"0x1cDDcD54D41d45dA496b2c74D3E606BB876EA54E":                                                                    true,
		"0x645ff0771422669fe52f5a643f03ba2176359f8c":                                                                    true,
		"bc1pwg7snkvndgwlae6p93p3w7ypns5l9sndw756utsq3rvaeh20j47sd8sfk3":                                                true,
		"bc1psrspq2u8e2jyy2dzgwxah8ty5xzdltfyhrdv9cpp63t5ympq22tqsectth":                                                true,
		"0xCDd85101c9fecCAE48174c8b56f80Af058fC6722":                                                                    true,
		"0x6245f1c86AF1D0F87e5830b400033b1369d41c34":                                                                    true,
		"bc1pw8cf2xp3y9h9up8ttc7tje6g80yf2kn9p80dknmergc0fqc2pp5qjj0s20":                                                true,
		"bc1qcvamd38l0lh8fagnn2ug8kzzzjranc7gv8ea9w":                                                                    true,
		"0x6239FF7d3ddC9E2F43ae780677866b5d2Dbede48":                                                                    true,
		"0x5a972a0d6ac4b3eea098fb5b1b34c26350c0ddb716bf8dcf5856760b99cb8128":                                            true,
		"0xA2dCB52F5cF34a84A2eBFb7D937f7051ae4C697B":                                                                    true,
		"0.00239773645652168":                                                                                           true,
		"0x69CA096066ADDc51046097D029c64E8254e4A43E":                                                                    true,
		"0x25306cb01f39cD5954c7ca253b56c6Eb70a6979d":                                                                    true,
		"bc1pt6s87n7vgp2225ggs69ymavgqgt43dy7xr0ruxer9xq23argejvsp3n47e":                                                true,
		"0xe72856cfb35380e0566bc8c64afe9692bc4efc24":                                                                    true,
		"0x495fb6C5215607Ab3b0DD1141b6c0c14a51EE766":                                                                    true,
		"0x9fA27eAB9d6eDc40748B64419D30EDC00ca50758":                                                                    true,
		"0xB12221Ed2A4CbB74b87824987Afb23d5D54c3b56":                                                                    true,
		"0xcb616EAaF3034882Cdb5dccF4D2b5B9e8E742577":                                                                    true,
		"0x9b237218651ecB29fA8298c3acF3464cFe6F0C96":                                                                    true,
		"0x1867331344F18b0846095e3DbAc189c156E2A2fa":                                                                    true,
		"0x0c07a7c8655DE7eCe5b2e61BAA3AEfeC049fF170":                                                                    true,
		"0xf2f6912Ee6B5BB19AD81E429B3EA211DF8D328aA":                                                                    true,
		"0xfB39970b889B6d463413698cf5B430261C4422cC":                                                                    true,
		"0x3f4a56E332e038dF29642B4d8CA93bEAa48Ab45D":                                                                    true,
		"0x0Dab54Ea132DA4262AB827aFC6273C612CAf97c9":                                                                    true,
		"bc1qjpz9cna08j8x84udq6q9z3p7uh97m68vdy78av":                                                                    true,
		"0x1Ad91f6c3a124dF9fB561fBAdAC2E6Cdc04A8BC2":                                                                    true,
		"0xFdf04bCD21a772399E174fA98AaE7E06f1655863":                                                                    true,
		"0x29D7D656F14Fb439a0ad7b76659c0a40c7064bDe":                                                                    true,
		"0x8647CaFE62478d08c4D70777C95fC85f95BD9bDF":                                                                    true,
		"0xb280A579a74a0aBdB5794D98FA63d9603cf3Fe06":                                                                    true,
		"0xE0acAfCAAf0d70207ef1aa061eA652edC7Bd22c3":                                                                    true,
		"0x21F628AfFABecC5eA4621c718F33c90cEC2a878c":                                                                    true,
		"0x16de09B46836952F4f0A5Cb8aef9F5A0A3ec73b6":                                                                    true,
		"0x87216bfc1A2c387cE2B5eE1F03F729A6D77A048b":                                                                    true,
		"0xA9D3127532BD15DCa2260fA7663cb724fa75Cf18":                                                                    true,
		"0x075e0a47e2Cc1171F482e817963DA6a23d3EBB91":                                                                    true,
		"0x8fF0e1538fC5D2B4074270080B846e903F78c84E":                                                                    true,
		"0xa56077dE721CbF8c6fa79D668a6C61fFf12CC37d":                                                                    true,
		"0x36cC2e55060F973E0DfAB855551511AC807A0c42":                                                                    true,
		"0xE89577666072d6A1fE646cdb417DEa50c0601661":                                                                    true,
		"0xca2F4D1A4D7499F2e9695FFB59071FD1F8601C33":                                                                    true,
		"0x998d36760B8C9496bfFD032f6F5d394c4fAda862":                                                                    true,
		"0x15960984538a468955bE768942C9125Ab18d0eCB":                                                                    true,
		"0x60CD4F61458850318843a1f50101DfE8dA914E25":                                                                    true,
		"0xc085c5193A6718f75a63a5728bE6C16153a3537B":                                                                    true,
		"0xf4674e10A65e979d3701Bbc46c786bAA87e8cBF3":                                                                    true,
		"0x51cF362E5e9C5c7e989dF2008C9a855ca4ca4F45":                                                                    true,
		"0x118ffbc4F196123113af66aE4113b54D7188a2Ab":                                                                    true,
		"0x63000327CE0De8B88b3608A1DEb26edD076229bb":                                                                    true,
		"0x99af385b6c478f8fed6a588dc23903fa7dbdbd34":                                                                    true,
		"0xAC693E7E769FAcc6F238306E10A865438a0a0E24":                                                                    true,
		"0xa3d2f547e109f8f84f617a3b192245db4f426a04":                                                                    true,
		"0xE3DA7e4bA0e7a49441276dCC888a11d7071e8C6E":                                                                    true,
		"0x005CA6C20FaE2020D4ac8b857c3574B3Df670518":                                                                    true,
		"0xbffba9998fe7ca283298576dffb16cf85075f47a":                                                                    true,
		"0xc043e2e5eac6097ce1c68bf64ec7bf273d87dadb":                                                                    true,
		"0x8e4dab90691a3cd1dc702f74cd123caf134b2791":                                                                    true,
		"0x7805334E9DCf2983B7FF7394B4a09f3340c17a13":                                                                    true,
		"0x17d187359f10AA692Adb597b7B8AF6d65F045C2E":                                                                    true,
		"0x5af26b73aF4b29006d88Ac986017cDa5126864Ab":                                                                    true,
		"0x5f0b2b0A5A7b363c92dE85835c80f6DF46889926":                                                                    true,
		"0xAd2f4f3F138cC7653DB40743CC69D3F736B089b3":                                                                    true,
		"0x02628CF23591c4b877ACccDA06d8285A2BEdc3c0":                                                                    true,
		"0xdadC093aF00CdDFE7b995bA810DB324bA3872a05":                                                                    true,
		"0xEF238505381b34683A07500D92547B60a1c3e0E0":                                                                    true,
		"0x95eE8B32AAC20a4Ca806D801774Fc5e31feCE7D4":                                                                    true,
		"0xd7e20d59e36198b4dfae1f38c28ce3bc2fa6ffb2":                                                                    true,
		"0xa70ddAD0bc272fb53FcfBD3399dF863d4FD1225D":                                                                    true,
		"bc1plxkyxesf9yj97p5u5enmhc75u3g94km6dph727fdsp6p7lrrs6zqwq8f7e":                                                true,
		"bc1pjenlryuwx69l6xrut5ffczuq5gv6d95wrdv5tg5xram9aa5a59fs34u978":                                                true,
		"0x4C4bCa4797B044c378edDfB7479676e1E18298E9":                                                                    true,
		"0xcb38cf6163886E25d457115db723b1203Aaa4F2a":                                                                    true,
		"0xb97C918897Cf24915aB3055789356279414d4A71":                                                                    true,
		"0xB6B18fEF653B15D81751aEDBD6894084ACD7FAd1":                                                                    true,
		"bc1pr90y8jzj6lrlv4vhufza43rzuvvmgtgqs0pqqg8wpy78f03037js4d252q":                                                true,
		"0x9D677028E0e592D48fc8a8Ddc910693301a4a450":                                                                    true,
		"0xFd1581fD77e9fcC57dAA9a7Ab92fD07d4CF3f59a":                                                                    true,
		"0x42d50AF7a1BE61207B2a62CC06869dd4294a8b55":                                                                    true,
		"0x7Ba3D1C4F46516fb975FB012F7db04381188e907":                                                                    true,
		"0xc6cc7f25ba045b8c08fb84aa1494b106fb6824a5":                                                                    true,
		"bc1pqra6hxf3gwkwefngf2fl69nrxdv57vclc8na74nppapdp7eexq6sfh88gw":                                                true,
		"0x12792eb7331DcC0eA9a2e29A8e12e41Fe1c22304":                                                                    true,
		"0x83be849f0f51D23F17e49838f2811E3178226d6a":                                                                    true,
		"0x124178CbA048AE0e7EF2C6DBc712B8E175F2c851":                                                                    true,
		"0xc0Ed15F86Fd4226F4893A5925e8170D770acE587":                                                                    true,
		"0x37252a3aa5f9DA0B445D272773eB90aac10eF76b":                                                                    true,
		"0x727b10A419cB83898e8f1C36bB31EF5Ef5eC606A":                                                                    true,
		"0xE4C8EFd2ed3051b22Ea3eedE1AF266452b0E66E9":                                                                    true,
		"0x1986f4bcc6b78d40e499e928a910dd7bde857734":                                                                    true,
		"0x02a5c980029cB470Ac89Df2E2de1CF453aEE6558":                                                                    true,
		"0x0D85EFF1963C3615044D0B407eA25ED0AdF06EE2":                                                                    true,
		"0x0FF165D9420fA705c8803aC1eC09725D887FB74b":                                                                    true,
		"bc1pjextd5p4xf3l9pz3f59m3n0ta97xhy5khcmxkl5k6kfpq7ltkmgqkzcl4n":                                                true,
		"0xDf30B739B82b0B463d402E716aE2Ff05a3203011":                                                                    true,
		"0x76750300a6bb602ffAC4FBb1f82346E20349309C":                                                                    true,
		"0x5D030B5CeC6933b9b0F91B09a0fD3f7b8ef2eB3E":                                                                    true,
		"0x7611353E3160D38aBfD1B5b79aB7461C30AC026b":                                                                    true,
		"0x874A80Dc85c4D5e060b0C93B51Ca995fE8c243bB":                                                                    true,
		"0x60D4646cD2676B1A2397A7B00F292522B628aBb8":                                                                    true,
		"0x1f6186A390190865c4932c466214D28E27a45e49":                                                                    true,
		"0x014B328D5D55751d7e8fD45cB0683e02D884Dc28":                                                                    true,
		"0xec41De09b42378AEa473365edD46FCa19ed096ec":                                                                    true,
		"0x43c7C3943A181774FD1791742EF6b42d671E30c3":                                                                    true,
		"0x14b3C6Fe2A0df86441f786db497a8173eFdD14f9":                                                                    true,
		"0x982264d588B2696a7e2D19393da6C04E3e6f5Ce7":                                                                    true,
		"0xE6C9819eFBcC89c8768796461Bc464a181871B9e":                                                                    true,
		"bc1psy0lpv22ttndv3esahl6jawldylkwjvg76zfx6rd8uzljr235fustupaqs":                                                true,
		"bc1pwnw24vgedkac32wmggnv6rnax2qrhx3hyflvu0snuvw2t7uqeu3qmcwujl":                                                true,
		"0x21357bdA562d6861E87bB7dBf850F6AAfE940919":                                                                    true,
		"0x7C4687272e5A6cda8db19eAFead3Eb24e0aBb715":                                                                    true,
		"0x73b643C692cC61496964CE2E0D5FF28bb3cdf604":                                                                    true,
		"0x67f72412a592d066a2e688e62664116deabeab29":                                                                    true,
		"0x3A409EfF50A47aEeF294E3f0BB3874490dD99abc":                                                                    true,
		"0xC66D22D1e94fDf125FBcDd5250B37dcF2FabFa5b":                                                                    true,
		"0x093dFF4b9517498642aB8F2510aC4f845a9620DB":                                                                    true,
		"0x8efd9addd8de6a4e64664d1893dec51f8c3339e9":                                                                    true,
		"bc1qw8d9yx4wcut4nj6uaqarxpcu54m7jspet34nqq":                                                                    true,
		"0xc200B58a4B620A77Df3C4415Ed07afacd2fef43a":                                                                    true,
		"0x197b118f84f0fC0Aa22b1DcffCE292EA75Ce53e3":                                                                    true,
		"0xA620570661caf299CCB1C547C66dB5b5c4D982bD":                                                                    true,
		"0xD5eDa0f579e31Df795AD2B0d01D65b7b36a0b0E0":                                                                    true,
		"0x1d5C0696D4Ea8387f02FB5289B9aBe69f08ADAFE":                                                                    true,
		"bc1quwtz05el0cl3kl3t6h6cvjjup562xyq3ujxy7l":                                                                    true,
		"0x3709575aDb4EABf12692Cb58C6ce7610f78622C9":                                                                    true,
		"0xE99a22F492247E857FafD8a74AF32c79262A1ec8":                                                                    true,
		"0xA6060C1F9a9e7E89Af6ae2198721049ee241B47b":                                                                    true,
		"0xe882c0cbb836c34447dc58f74c08f373a6d37cd8":                                                                    true,
		"0xdFA413375306E2169AdCBbE8551f69739328E6dd":                                                                    true,
		"0x952c23f8F067A5e7e165ff0E42491f51D87DBc95":                                                                    true,
		"0x4d48df14a698984C60e705a11dD696F368cb3A21":                                                                    true,
		"0x04621637c3a72db27d6f1Fc1352480Fd4D1e7B06":                                                                    true,
		"0x522DD4BEfC804EB041028f2b5c9499a0d798E445":                                                                    true,
		"0xFE09C68917E9bd70f1dC739128A3B53BD5be8e6D":                                                                    true,
		"0xfB7e0d2d6fDC661F72efbbad15735FaB88966bbA":                                                                    true,
		"0x96b289AA02D0BB8137249fB0C73235fE1364e6B1":                                                                    true,
		"0x63cBF2D7cf7EF30b9445bEAB92997FF27A0bcc70":                                                                    true,
		"0x8FCcF9F119cB5dA9629072277f65B5abFA35F017":                                                                    true,
		"bc1pcfm4rv6z9qlmaycjf60x8xh3lttz9txt3jnh9p9eugssmv762s8stawekj":                                                true,
		"0xd465d8a36dCed7d3617EAD555EA225B51516a631":                                                                    true,
		"0x8Aa2dd81fb8f986fa5dca27c6Aed9E2170c8851e":                                                                    true,
		"0x739ba0465198fF7Ba88cd7A2C4684734243D7677":                                                                    true,
		"0xcc3415f9fd1e170413a776266bcef6ecc02511f5":                                                                    true,
		"0x46337De97dE52C5A8eF04b539CE68837f81ED3aF":                                                                    true,
		"0x7eA4464B13fA31dC9fF5140F5216A0c72F54C9eE":                                                                    true,
		"0x0f553d211b1ACF5bA914BBfD5d6D70DDD9B7425e":                                                                    true,
		"0x59E36f05da46Ec9597d28F41f45d9e79290CFfeC":                                                                    true,
		"0x44ee8bcf09dbe52ec9336d357ac5d584ce3dde2e":                                                                    true,
		"bc1pnddgdrlsk7d9lxutqxvc8qzathw45wxxk6dfzu7asz2uj0jv7pjs0xzt7r":                                                true,
		"0xcC3415f9Fd1e170413a776266BceF6Ecc02511F5":                                                                    true,
		"0x50ac5CFcc81BB0872e85255D7079F8a529345D16":                                                                    true,
		"0xD1E1cb24961D43d6Cc25dC001B6332d6fA67888c":                                                                    true,
		"0xD3c7485ACcDf1ed393812B86d397c08ee4Ae5893":                                                                    true,
		"0x0Dc65Fd872C3bE4DD1ACA5EE8621C408402f0729":                                                                    true,
		"0x19Bcb9b9d7ecD218469957EAfa140731c693d9A9":                                                                    true,
		"0xF15b97b82981C43884E264B4E8D908F65c1c63D3":                                                                    true,
		"0x0Eccf3295b5DfcFCe69044EA157D55Dfd0f0B955":                                                                    true,
		"0x38394235964A44c9f2EBbC805429933bE2e63C74":                                                                    true,
		"0x84f9aeCa6F9eCacA340Bd9535B4f642EE74360Ae":                                                                    true,
		"bc1p7239t08kehkc0hp5c40q323yj2gusjxxuxq906t0hp2c67dgprfq6dvvds":                                                true,
		"0x0967774931B50f96863a0F2Dcbad2E17AA392A78":                                                                    true,
		"0xBf593D86d40Fa16fFE71138961e68faDB9F74cfa":                                                                    true,
		"0xF0B131A9EAA9e2c8f1d26D200D47Bc1eDa50FB66":                                                                    true,
		"0x8ff19817bB771D6EC28376fC5D872FC8727d7Ff1":                                                                    true,
		"0x0a0FD254a278C9099D4BDaDb9Da81B9dE4Ddc61a":                                                                    true,
		"0xe41be092c12Dc99eea3f8E91e0CB4Bc520d9705A":                                                                    true,
		"0x8bb9B5d417a7fFF837c501Ff87491047E1598104":                                                                    true,
		"bc1pj5mcp2t52mk0503ql2e6a3wjlqwjwjsd0nv59gpp7v82l4t83pyq2f2krm":                                                true,
		"0x9819686b54cb0ccc8161e84e22db7ae1ec54960c":                                                                    true,
		"bc1pv44ua5k4mez673qdg9g52pkf62y9ueunn450u054wkttltmen4tq0apf43":                                                true,
		"0xa8B44d98F6092529deb51fFF42a69Ebd00122409":                                                                    true,
		"bc1q8m6srg9k09xzu5kepdfmvz5dj8fs9zcnnc35sm":                                                                    true,
		"0xDCcF17655800aBf354743a2F65EF3b8c8575c3e3":                                                                    true,
		"0x38B7cD11a2F344986b7F89683d781C7cac798fcF":                                                                    true,
		"0xf597dcb772d821B0fF0B222c84c82e16294cAE99":                                                                    true,
		"0x4f3e48DA06bcE4D17ebE7A07118A04b6d1e555c3":                                                                    true,
		"bc1p2clxzsgyq2wds4m0gs33s5qeu884zheguwc6tcuhcrpp2vyjnh4qw8cpdg":                                                true,
		"0x3d99b50cb54b445Bb9f29CF6C995a4f4148A5fb8":                                                                    true,
		"0xca85bBbbB5060C18d21408C2F80CCcFE718b6986":                                                                    true,
		"0x9a19bCABC43B0469Be59ED2497927170ffC968BA":                                                                    true,
		"0x775C8c679B5cAc282894DA4E18C40fC3E5f7F488":                                                                    true,
		"0x6993fc084A8CD64729D503019C863e9486C14966":                                                                    true,
		"0x4ed39Ad181820F8729cD761D098fa8bbc2BE6d40":                                                                    true,
		"0x66208416CB1dAF88716af68E405fcC8994c059D2":                                                                    true,
		"0x9e702dd2e1Ecc7593838576b100882002D7f6cA3":                                                                    true,
		"0xf416526650C9596Ed5A5aAFEd2880A6b3f9daEfc":                                                                    true,
		"0x4cd7Ef1F6719c75E30C4289Ff857504C59466Af4":                                                                    true,
		"0xad236b7b8d774BEfC6113794F7432803D498Da23":                                                                    true,
		"0x69F8122e685155D0338119a7C307D2c02473B1F1":                                                                    true,
		"0x6d2113692857536b8993924D6B66d8409247fBb8":                                                                    true,
		"0xcDbab8ff4DE3c9Ab024C4f17CB10b82E25ae16A3":                                                                    true,
		"0x5Ae0AFC22ee9a5F49a13A749e39E20Fcc8021f3F":                                                                    true,
		"0x322b4D1Dea0213047Ff23Dd7687b6E0FCC78e546":                                                                    true,
		"0xa425da3eC37DD33D72236C33C1ec225Ce56323f9":                                                                    true,
		"0x8081A75141dBC89f70d31eece08FF12cfe105e43":                                                                    true,
		"0x535A95120E925EfD8EA5eA891836FC75F919E82B":                                                                    true,
		"0xFC522edEd95e7E4e93C7e7b3E5634635BC82edf2":                                                                    true,
		"0x121AcD7E5e24d4E426724E39D0FF449b4C89c601":                                                                    true,
		"0xFa4dc543531B5b176EbD03209d4b18b575f76A52":                                                                    true,
		"0xDeA09Ab1F4Fc1f2395Cef3c7eBC17f6C27a6AE8C":                                                                    true,
		"0x650211A2809779e609cAe1Ef0F864345bFdB903d":                                                                    true,
		"0x1402dE507719544277BE20043b6c6717F38d362e":                                                                    true,
		"0x1C6aEb6F93F36d9fA6D2D3123698d5663Ec5E2f7":                                                                    true,
		"bc1p2f8km4w9evekramcq9xe94jn5cuuz3rqsuuqxz69pm0x0mesdx6sgn2us3":                                                true,
		"0x21688D7C443844EBF63427c4e2AbFAAabc50a405":                                                                    true,
		"0xc25a6189D807A2551A6E5bd1D41F4bb52288E1Ae":                                                                    true,
		"0xb93EE4C03BA63AAC8D08008c0DE91b86BB0244b4":                                                                    true,
		"0xc8b506B8150Ea6Ee0cbd7dC71034B12d902141d7":                                                                    true,
		"0xC0fDd3CceFb7591F9C57aBD38E847181200E274c":                                                                    true,
		"0x5a8C44d9F71346012439FF06c3964c45A45eF878":                                                                    true,
		"0xD3AeEB09b5ba94F7DbB61eFC556f11679294AD25":                                                                    true,
		"0x02E1A0d3e5b5Bf8E6C518C2Bc6aF10C365fBD93a":                                                                    true,
		"0x42eb484bf06464351CB1907190fdf093346B89F3":                                                                    true,
		"bc1pps52mvns3tmapdmp2p60r5mhgr6wnsq9ps02fvrwyqug79m0e67q3wegkj":                                                true,
		"0x35F8c88322964b4e23816F6d14cB988a5B91Cbda":                                                                    true,
		"0x1c4Cfb61Dc235f96ae7606A56eD60a74f4585C82":                                                                    true,
		"0x0358E04cfa98a22C890f0dB844D5c121585F116b":                                                                    true,
		"0x2fc4ab5A82D4740d4015Be233a2fb4bC61Ecc073":                                                                    true,
		"0xb8a8b61F912872E0DD66397b41aDa445D71BA6f3":                                                                    true,
		"0x8bfe553Ab794615B6009fbA5F077509A453DCBFd":                                                                    true,
		"0x54f9cB1FA035772Ff4c6e276B7FD758096653D1d":                                                                    true,
		"0x05b27974688b1078A105c7f5BAdab65a5A41431C":                                                                    true,
		"0xA37FD8dCb5bD4bAEC2d1C96A1E3Df8A7901e7364":                                                                    true,
		"0xd2DE82B236c5b31609caeda354CAfd70E6b343B2":                                                                    true,
		"0xabc6E77893F72C5eF1E6b2a63DF609b4ac7DBFCe":                                                                    true,
		"0x678283f7dbd0C76d89F26B57F906725607ede785":                                                                    true,
		"0x84576F263698d4bad6f341887b64a1666E00D745":                                                                    true,
		"0xD9Bc1fC51b8D8369980505aa96e975DA03346B4A":                                                                    true,
		"0xa1E84210239baD5571171a8fe304A90E7Ffe5189":                                                                    true,
		"0x7c6FF9F13779896CcB81b7e29793526bC6bD0849":                                                                    true,
		"0x92313BE60e87808948D0F31dc3B029d7aab4C5e1":                                                                    true,
		"0xf3512206636Bb52462D2da391c80535D1ad6D4F6":                                                                    true,
		"0x2DAfbc6E546D8F40a94dd0DaC3fee0aEe02c1311":                                                                    true,
		"0x4B642f06710edE89aA332820A909F3e97311671F":                                                                    true,
		"0xDdfeE28a4AC343517aafFBda263A446F812B5229":                                                                    true,
		"0x6133571B04154E852368D6b9BceA281525C5beF8":                                                                    true,
		"0xcE4fE64177947D885bD1Fa6f9F1D54230937571F":                                                                    true,
		"bc1p5kvkn8s63aw4qvt0zjj3gff2pks78n4cs5sdd5yzmf3l5refs50q0hd6qw":                                                true,
		"0x890c343365C5B0380e6f532b82437cC5d0B31199":                                                                    true,
		"0x02343D7c62E99Ca155C27aa3a2edAff06F47E1B6":                                                                    true,
		"0xEb741575c5fCB58F24EaFa8f5654EA10eAde36Bd":                                                                    true,
		"0xAD79d347D42b9709E58972Bd0E4A790157227492":                                                                    true,
		"0x8FE2DFcf2F758dB5aC625991A04FCF591c0fA9eD":                                                                    true,
		"0xB75b6976373a351708177f04a7a22A7B2a20DDcB":                                                                    true,
		"0x6afCd57Ab22Fa837C28fCd168725Ee43463ED7ea":                                                                    true,
		"0x4A13A55003EEaf84421ad16799bBbe68D851BbFe":                                                                    true,
		"0xb40266773f74dc193E59f791c6E2F99F7d8DCB90":                                                                    true,
		"0xc886dB8b8CD260f5ee38Ba3d8f8E9324EE27EA33":                                                                    true,
		"0x3478B07E6e6a39Fd668B036136C5ae5f62135Be1":                                                                    true,
		"0x336F6BECa25Aed6bc4B4F9E6eF5B3Eb415aEcA67":                                                                    true,
		"0xdaDcf812B2eA669F8620Ac125Cdd672656CF8a62":                                                                    true,
		"0x5De2D0474D5fE1AB7eA9174B7075b6bdf84Fc8E2":                                                                    true,
		"0x1236B30Da2177bead2FB21A204e97663e9aB80c1":                                                                    true,
		"bc1plu54v662pe89xa86w95hs7ggzgkdm238cunzw8chks7gycmsr7zqsp6m6v":                                                true,
		"0x86d7B47a29ae3dCebC6bd1B6071115CC42AfABAD":                                                                    true,
		"0xfBEf51fBaB7EF54B523364CE229555E378908139":                                                                    true,
		"bc1p3uhs5dkunux004xavay0hj75w5rl2hd80hasyzp9d9jzjl0zzmxqgeffac":                                                true,
		"0x78fb6f60B6D795dedd90c187CCBb7fBc90fa61ba":                                                                    true,
		"0x694207F27181889453e8607B79Af3ACBE413806A":                                                                    true,
		"0x61fF15Db4fF75b180CDc3351F71FCC288606366F":                                                                    true,
		"0x6D745Ff5859A2d770d5Aa0261b0155A25489d837":                                                                    true,
		"0xD803556e336807Fab0b52772e7085d8193Fe8Da0":                                                                    true,
		"0x549c01f812E609b80548c9837b2DaCB0CB0ADc64":                                                                    true,
		"0xB3897952Ce7a4CeA159D232dAe8D02ca8273372E":                                                                    true,
		"https://magiceden.io/ordinals/item-details/0cbbb307e286c2bc7e545e39f7f57c02d14fb2075fa3f5e2fd7fb20e9706e142i0": true,
		"0x6027fb3d162174976b84803caeb96d3c1f4371cb":                                                                    true,
		"0x36d9c5e40bEcdC523D8388CAF8345eB03E0FfF21":                                                                    true,
		"0x1d13949903346D96014eD0264f0109b4aE5d4185":                                                                    true,
		"0xe73A56c786ac755cb0729A5D429b9f4129F743F3":                                                                    true,
		"0x69cfe6F383dd218DB7Eb5A1e10345B9d7BB2dEb6":                                                                    true,
		"0x03E6aA09C8F3D5448BaB141166f1fd9B98E62709":                                                                    true,
		"0xdf5d4cE14b38D75DeB7208f2517B10E9e41eC8d8":                                                                    true,
		"0x5BE3bD1d263643BE7582DbB7cB453873a0EC6657":                                                                    true,
		"0x9A2AEFa2FfaC7959B6C04b05584950dEEaFDe9eA":                                                                    true,
		"0x3e917eAdd649784aaFBEB5Cf4C9132fce42f9f90":                                                                    true,
		"0x4b0d93ecdADD048EC4cb10794fa48624E23a8c2d":                                                                    true,
		"0xEC5C5E6D9d8D7718BE716137BB4a55d84B1eB602":                                                                    true,
		"0x7937a23DD2632395B9616919DA176aABda00A549":                                                                    true,
		"0x4344B8F69f14F227E197f7811983564b29Da48A2":                                                                    true,
		"0x2275B8be039e99F88835CE49B16285aD0e61d485":                                                                    true,
		"0xc06ac1f7e7ee52c90d5730808e865f8d758aaa00":                                                                    true,
		"0x484053B5828874b34Fcc755E1857f95F2948101E":                                                                    true,
		"0x4DDC70aAEA98A5D4c17aACaCC9D3863DEe399FaA":                                                                    true,
		"0xF1c0341935E4dae98F8607736A303e5eD158704F":                                                                    true,
		"0xF173429B82A485A800e2a644592CCB0e59706Cf3":                                                                    true,
		"0x557FE7607c0bf066e925b1adD51e2756b36F744B":                                                                    true,
		"0x2451945F2A788c7a83Ca7C57C6e7f0278849BAA3":                                                                    true,
		"0xe1e2E153e785118Fa7b50e93e6b8cf887F025Bc7":                                                                    true,
		"bc1pm9mtpey58mr69h7fryw3st99n2mjazrrp7mdw0xnve9dmmmadclspru8a6":                                                true,
		"0xAD9A35D9B4C256ee79bDd022189D18c4426D3d53":                                                                    true,
		"bc1pjkl7ujhk3877y8yk4sdy970mysjrly5k6ne7am8l6pe2g09nqfsquu059m":                                                true,
		"0x0C8C9b24ed088B73B367dFf21dC36e33565b0583":                                                                    true,
		"0x22426E727390fcb542e829b3dE242D9D9A234cF8":                                                                    true,
		"0xdfBaeeF21396BF205D4B7D23345155489072Cf9B":                                                                    true,
		"My TC address is: 0x80b0a2e108c483b3e7ebdc4b1dfd8f4208266021":                                                  true,
		"0x80b0a2e108c483b3e7ebdc4b1dfd8f4208266021":                                                                    true,
		"0x7318F778c2ABeA6Ac5F89c0Aba1803e12f9ddde3":                                                                    true,
		"0xe1127dd3bfc88ae61211c99e36655d92d6aee28d":                                                                    true,
		"0xCBA711BEF21496Cfd66323d9AEA8C8EFd0F43e9d":                                                                    true,
		"0x6bf189130C6C97ad73074e027Fc9128fADcF5832":                                                                    true,
		"0xf2C2140595806d7FD78C41098B032162BE61D487":                                                                    true,
		"0x2823150E733267870d76839CBb9D3c53c9ebd658":                                                                    true,
		"0x50D462Aa9e091A6D6106A2A6c0063831B04e96E8":                                                                    true,
		"0x7869A0A12ab067C992d05EBb7FcCC8A7Fc9d5103":                                                                    true,
		"0xBE218041d6f39CfADD7fA0f1CBf582D8d898BD57":                                                                    true,
		"bc1pgeu6cferzqp637ea23xtek5vfgdq6rwce4kelkx7sq5nxxtquy9qgy84e8":                                                true,
		"0xd3db9d11c09cecd2e91bde73f710de6094179fa0":                                                                    true,
		"0x55c564139d9C1d59A3a85e8F91a12914cb32Ca2A":                                                                    true,
		"0xe60F9490D33D0c557f6d7211f63A97bBC7350C16":                                                                    true,
		"0x931048F06D0EA84a84bdAfebD058dcF9B9762501":                                                                    true,
		"0xAA3044894d1BD6c2845AdCe181c818520f40b7E9":                                                                    true,
		"bc1qa44azjpyz6qk2m9prnqq453649jxw6epauf3yt":                                                                    true,
		"0x9418984098CE76041CeE56694923393543bd7aBe":                                                                    true,
		"0x2b40825F8677de768c6B061A600C52cEa0031377":                                                                    true,
		"0xD767bf953D355104737748f22355C0433aA9f3A6":                                                                    true,
		"0x2c8349E7aA2BbCC2C3EE5bf9475fd4fB885A8f1f":                                                                    true,
		"0x1F2cA89033dE3526B987686eFD4928485953374F":                                                                    true,
		"0xE2f42ACd9A98173cAE4Ec023E73204F0ECF0E432":                                                                    true,
		"0x53FAE7bD5101772bcc1D4040f61A4FE4D02DdA5E":                                                                    true,
		"0x58cA09A20B25C3546Cf2b5f411cde799E9A159Cb":                                                                    true,
		"0x7089943B03C65A7099BAe54CA9510a55E881F23E":                                                                    true,
		"0x3532Bf9FF4900f32DA9037b2cd0188A0419BB7b1":                                                                    true,
		"0xE5A70A52c720116b13413E262482cd998320539a":                                                                    true,
		"0xe0C6a117FC7F5dEC85FA544470497dF82BF047d5":                                                                    true,
		"0x8f23EC6Cb459529F147068298D98CAff2E62568F":                                                                    true,
		"bc1qt5lyqks0m8q9he4hy0u6yehmlk2wscued0vnk7":                                                                    true,
		"bc1p6v8y0nmuy8gjlqvrcykcfxeyh39c4luxzg6r8q8e480gxtc4qqwqhyrguw":                                                true,
		"0x1403587BeE7B4B9Bcb2756A6647851b831215371":                                                                    true,
		"bc1p0mqs9kl9xak22qrklp0hewmefx37pljsss5ucdprg96vnkrwrthqr9hz5g":                                                true,
		"0x780dED412E072A56bbCB7249Acba42290bA4262d":                                                                    true,
		"0x17b8909a6fFA2C0255BB0dE6Be9E5A4A768690c2":                                                                    true,
		"bc1pz55y3zs82v2fyw7pgpsrq0z5eh3tm2e9q5q775zl2rxw0t3kqvussx9a0p":                                                true,
		"0xcF690F05D9b3E88164371182F2eDa3E3349175D4":                                                                    true,
		"0x4be41356e363135702909192943D990651452B68":                                                                    true,
		"0xE0BCBE59Bb04bEC47cBA0AF01d2086dE967c0279":                                                                    true,
		"0x9dd174ac950ba3dd4b9ad37e76372f3e4718337d":                                                                    true,
		"0xA6c53c8842FdABEa235EA93c6C5555B5F5aDe118":                                                                    true,
		"0xD9820cA2C09b98d5DF5B36C3879906A236d0B63F":                                                                    true,
		"0xE306223f8069991449CAA1f46c4f097Ba60817B7":                                                                    true,
		"0x9e6484512a752034abCCc0FED8038066406Ad8D1":                                                                    true,
		"bc1p4v8rwzx0lwvt750g2tlzk9xfch5krkehpm963p9u0hn0rknd930qppdtmv":                                                true,
		"0xbA5A17CE9C6F3b6287C9B45556b25D542dC726DC":                                                                    true,
		"0x34333Ed158f50d7AE32873C83b908Af52c48a8a9":                                                                    true,
		"0xb9E99E151b436D9bDd46CF86caA345683E51a3b0":                                                                    true,
		"0x8115949804c22e03c95349a67c038c55730a4ba3":                                                                    true,
		"0x8551084392E882b4AeD21903492a7fA59338387C":                                                                    true,
		"0x3582138E18b5acA135B73d4091A022401e04f2B3":                                                                    true,
		"0xD5E7C8051bB55471e65c77735246037B88887794":                                                                    true,
		"0xDBbce16eDeE36909115d374a886aE0cD6be56EB6":                                                                    true,
		"0xb8057538733a5bEd77D72D245e32892aE4cE7F5A":                                                                    true,
		"0xFdb2cE90eCbaa90502c3AC7F22298023fE3984E2":                                                                    true,
		"0x35Ab979d9CdD994B99C04DFcFfE245A77eD4e9a4":                                                                    true,
		"0x00Bf0fFF79BdD5e34a6881d145De2F21f9483B02":                                                                    true,
		"0x251384107FCAA63b45e15b08804AF40db90A922e":                                                                    true,
		"0xbc38d1699d5472BC4D9CC2CC738B0b1052291334":                                                                    true,
		"0xef34138b7D8654b0bBbDEd60EEe651473fC89BA6":                                                                    true,
		"0x1749F4aA848Aeab728f1D183E24dF13cf22864d0":                                                                    true,
		"0x5Ef4f751F6590315422FA58432D7D5D85A7e2572":                                                                    true,
		"0xe815264a16e15dd1d8c9404905432a0D60FAd92b":                                                                    true,
		"0x5E8803F5605923026884052b6693D00C2A3cE1E4":                                                                    true,
		"bc1qqszneqg66mq52ac3tdgszdjnyc9mkg3ef9ajtc":                                                                    true,
		"0xab7ad8d0a44971a26ef78395648837998d78f87d":                                                                    true,
		"0x5C3080E822D167B5ef15E27718bcfE465f37a94B":                                                                    true,
		"0x9E77040200585B0f902A22e270E250644e06554a":                                                                    true,
		"0x1c37cB1b6fb4fEEa025bc4c432c39287FF112371":                                                                    true,
		"0xD37DAbe2CccA72166F41c236b91Ca36494FE5919":                                                                    true,
		"0xB9295F4468c22D5FDD718131b747fB2Fc7eeEb71":                                                                    true,
		"0x14bb7E16442CD0FD05130251d8B3202C6A8213f9":                                                                    true,
		"bc1pc8v25ux03amkyr5kvlvtz8s7j8kc5jza6cmar5ym7ejs3aegvqwsq2qdss":                                                true,
		"0xF1D8FEc8C8B40167Ede1B1d007a9B73D41a949e4":                                                                    true,
		"0x3c6d37f99f6ad9810191144b1f7e11fe6bd500b2":                                                                    true,
		"0x76f89850e596E201B0B490872Cfa3784b7eb9FAD":                                                                    true,
		"0x1A4695Aa5d9895Ea163F8dcdD46687D3C003cCF3":                                                                    true,
		"bc1pq2je7f6sscuqgvpzjs7tvmz8s75hpxjhcfz6qfjymqsdjaxecggqxxshzs":                                                true,
		"0xC35f7Fd97E9f17DCcabF8e523c8F24b25685167f":                                                                    true,
		"0x3A52749372DAEB2625eFbb496672599BadDC98a4":                                                                    true,
		"0x8C21aAB953367e055EafbC21684cba1c360014a2":                                                                    true,
		"0xdD8502c9aBd9eF8ce5016096E7838e4A377a07a3":                                                                    true,
		"0xec3bbE3a2b68859CB52b6eeC85f085Cb9d009D7a":                                                                    true,
		"0x43FC6cE5fF85E761F024E4c315bad1425b6a2617":                                                                    true,
		"0x57C78847C68CBC61D1B2aa28cc7dd925785a63f6":                                                                    true,
		"0x7FAda7c77032e2812F87d656B7fd5076A9dCBB50":                                                                    true,
		"0xCb396477694118363ebfeB6040995615432fb88E":                                                                    true,
		"0xABf180aF2aaFbd32098ff898928BCfcdF4277D9a":                                                                    true,
		"0x0E4B99CD9202Cb057a7cFC1EF825553f29D1eC07":                                                                    true,
		"bc1p5tz5wvnwfcws4866tuarz3us3csz5fycnphgv9lj2ndwd5lakyssakhnak":                                                true,
		"0x78a06b4D01030F5407C927f2A949ff5C03a6518c":                                                                    true,
		"0x4fdA2dAC61633Bd19d165e7c72AD2Fea581ceF07":                                                                    true,
		"0x2c99D76ce6f94ECC79858D5ABA5f8b209A0b44a9":                                                                    true,
		"0x4dA792b5058F59162E1B619749a0Ce4E984D4841":                                                                    true,
		"0x0FB4BBFCABE6D9239F8f7Ec2b6EF83F29B9F4278":                                                                    true,
		"0x21f6302E04744eF90bc9ee65B7CA556E078df9b6":                                                                    true,
		"0x6b6f9c634769dfb2e7A44A43d5658a28922CEA04":                                                                    true,
		"0xB6825fe2feE68e8C0d2781d0D963fbFCf6da0487":                                                                    true,
		"0xc8AB2D82724fAB1e8B6B448219fE1168BAF0b0c":                                                                     true,
		"0xd0daddf983fce88bef3f10fc12280d0f0cd120":                                                                      true,
		"bc1pvv89f7g8wwkx4dpjh3v4hqa8gyf3dl8cxpamw9qwecgkwa3zn8ds4v99t9":                                                true,
		"0xDc79fF68F761c84Fc4Ba13c4e9c695a6c078D637":                                                                    true,
		"0x01180F770161351e946f6665befA13beb56898fA":                                                                    true,
		"0xB3557ba0d49bb21b43b2a5AC1dB4b5258b8E6640":                                                                    true,
		"0x01Cb6466c3576B83CDc707f63D0Ba9d34BA76c3E":                                                                    true,
		"0x6Bc55429118fecDd60170Dd41793d7E3494cF737":                                                                    true,
		"0xfF2c7125049B649b5764221A0fCac138F0B3ffF6":                                                                    true,
		"0x35295fc8c0d41C5683E23a71963359031fF48bd3":                                                                    true,
		"0x6e052818116298f292514E3f78789439002F4249dd":                                                                  true,
		"0xc8AB2D82724fAB51e8B6B448219fE1168BAF0b":                                                                      true,
		"bc1pv8vrrr98e7k0pv7yj490khe5494snre6wjcw3e7lu3ulpl7er6ws8v6mlu":                                                true,
		"0xC98665DBB9520B073474Bd1b9E64d285C5fEB13a":                                                                    true,
		"0xB552740655d9719ebc9a058b977a1AB69Adeb622":                                                                    true,
		"0x1B4DfFa580d972525DaBe0338c7169881BD1d3":                                                                      true,
		"0xacaa84c31cd2666c759018dddd3df296362d7a2":                                                                     true,
		"0x2F99f21476D0706070443803a688ECD6C60feC9f":                                                                    true,
		"0xD283AC6b8dD58CDE7EdE41002219D976115Dca36":                                                                    true,
		"0x23e304897cc82F2aA1B78D8cD0a5fdE9bB272d03":                                                                    true,
		"bc1p5emyjpkrvvw0s9r4sxfta698pgtvkrzwr7p0t7qngqqtcr7w22us4e609j":                                                true,
		"0x8BcCCD0D7C57201c3EcF8C4e115A628851D2326B":                                                                    true,
		"0x8C4C29135a819687a72413b2113A95a682E83a27":                                                                    true,
		"1G23ukpkFurXQyeohVgdf3zMjN9cuvKaxf":                                                                            true,
		"0x58047ac4BDF9c18F83177FbdeA9358A8F824b083":                                                                    true,
		"0x38fBf1121086c6c713EdE32AEAF92D935E0fB4A1":                                                                    true,
		"0xAB0987c5c5bE39Aa7A892309fE3f949d4ebC631A":                                                                    true,
		"0xd0926e79c8b37A2DE3992415376f36442C21C018":                                                                    true,
		"0xA79Ac8e75cAbBBA83b03Fc94EB50F0bB7850E682":                                                                    true,
		"0x7BFD55db29C30a9F13a1bb8df403415Ce6103f39":                                                                    true,
		"0xc5F8ef8870E58B9d40f73Df6C31Ab74Ebc4b5E8e":                                                                    true,
		"0xfed74f78700bB468e824b6BfE4A2ED305a9D86ba":                                                                    true,
		"0x65a0E5b625EDda8F482C71C22aAef504E0B0A0AE":                                                                    true,
		"0xA2b5AD8b73f4790C4FFF0921EA9Dbf78Abf5254C":                                                                    true,
		"0x32D0a0542E62950f8D48504489405450e9c0AEe4":                                                                    true,
		"0x5b1f0DEeDca8E61474515202AcC5C7564B08291d":                                                                    true,
		"bc1q8tk98ln6crf6fpug80nk4dtsemvkd6p2457hr5":                                                                    true,
		"0x59895b132680761285415f118e03beb27a8111c4":                                                                    true,
		"0x303c861ac471b7ECfd777EB29b88367D92BA548E":                                                                    true,
		"0xeE77369a0822c08df9F7d09E3A53A104b32c5784":                                                                    true,
		"0xD63e6532454c67E32aA7F249ed6FCD3674F7e962":                                                                    true,
		"0xf0eb29FC0a4B7114B5D91Ce947796606E776541e":                                                                    true,
		"0xbb778D61779161bbD134C7Ea52Db78423672189A":                                                                    true,
		"0x9F77094cC3dA711b02A5cb8980E513348e0fD516":                                                                    true,
		"0x75fE3E054C67f0090303E04E8a5381DD792bdA1D":                                                                    true,
		"0xde6F3Ec16305780794F53733f7266f3295545973":                                                                    true,
		"0xBF737A06cb8D36C304988f749DEdAc4b443B46a6":                                                                    true,
		"0x499E76c1B1120d7B54609bF2483dA35341352A0b":                                                                    true,
		"0x4ab150645Be04D295F9F28B5eE23AA28Fec2fDAe":                                                                    true,
		"0x563153823D702516F92fc24edD9358D6973f60F9":                                                                    true,
		"0x8B853c7d34e1ac8Fb4b880887A097E3931cC1534":                                                                    true,
		"0x9C985145E86999E94c5c0F9EdA35fD30D0bd273C":                                                                    true,
		"0x81a55494572fa5A5474fF2DCd506C0416A8f8EA7":                                                                    true,
		"0x08b3cd07d27CFd5057D6BF5F64df16aFA912E703":                                                                    true,
		"0x4aA8219C910DCD33b7CB2913f5e4Ae955F3345CC":                                                                    true,
		"0x0286a22F655F84c36Ff6C80eB566a5a4A8F07541":                                                                    true,
		"0x2aE8512b8F0399fd4348B2F4b9a50D03a5a62AF5":                                                                    true,
		"0xc75265Cc1928eC230D53914F3e9bC36845BaD820":                                                                    true,
		"0x84bbd969db8f05a20d08933e5f609edb7c769503":                                                                    true,
		"0x57289AE13Ec410Edc3D6075cA07741b30CD1C99E":                                                                    true,
		"0x52Fa10126103995fA594a76CF9DBea8E01789c98":                                                                    true,
		"0x780423e2216405500be0670Affacc5c80eFd3c23":                                                                    true,
		"0xA0bDb2157E09b032c0f1A0972986C6b9b834A569":                                                                    true,
		"0x6031c003b204D74A88Ec6E84c7ab40F3b81FD9Ed":                                                                    true,
		"0xDe54227dC7cb1dE999979f21548096D92B64827f":                                                                    true,
		"0x84Bbd969db8f05a20D08933E5f609eDb7c769503":                                                                    true,
		"0xfa714ce7ab8cb29625a06fb7e39f9864ec3a05d1":                                                                    true,
		"0x052c76b1641ca932FC23278B5A2617f3Dc5cC485":                                                                    true,
		"0xF8D8a27790a6aDDb892fb95c1236E68f788Ff2E5":                                                                    true,
		"0xA67b03DBB74c8a8DBb8Fb93C5d554275817acbF5":                                                                    true,
		"0x895a258cf950d5b0b8538115b6906770fb82b0df":                                                                    true,
		"0x2Af40bC2F3C1A8aCF1a55b5874F9b65055D3e11B":                                                                    true,
		"0x2d15d84e790226F3AA0672110517eA6304e29cd5":                                                                    true,
		"0x0Bc483B0B506EfBeBb4237B467818faC94235804":                                                                    true,
		"bc1pz5z7n0sk28wny8tj2h732acfup4jxc53est252fj32g9pq7wlqvsrsjs6z":                                                true,
		"0x755Ae635904b36fBd8D61021f719cbFfa79c2fD7":                                                                    true,
		"0x2c7776827b7DFC4Cf0F3b48D3d7FE4896946F5b9":                                                                    true,
		"0xd7D743254BF6BAE1a509f96D0369eB7f45A6F190":                                                                    true,
		"0x10713cF63A5cbD1183Ef2d8c9A8D3192984e8126":                                                                    true,
		"0xb6692a25C624464f53F6e7fa978185A9659F1c78":                                                                    true,
		"0x43105A71c046116ed80eF4eA323A6A476c5545b6":                                                                    true,
		"0x5f7c27402603D2607b07C19FcD41C84A1f5557ac":                                                                    true,
		"0xD966aC73Fd52fE3964129079eAc60D04f46fE8EC":                                                                    true,
		"0x4aa8F01b133618611fA9718e6c4f25074B65220e":                                                                    true,
		"0x808A023B72260170c95d831F589A1ae0DCa1e43E":                                                                    true,
		"0x68C37e85Bc414F5D7d8E30580fe16E1079eC85dD":                                                                    true,
		"0x769a38BD7f53D5d281dB932bF2E3F4F071fEE2f4":                                                                    true,
		"0x5e624A7Ad13b5c01d547B1A95A386D1f6147Bf56":                                                                    true,
		"0x323F4F140E2A38A7604b9cB1A0928dfa5566ec91":                                                                    true,
		"0x57A88770cE4583332d4b1eA3EE96F2Cc8bF0171c":                                                                    true,
		"0x5E3A437412b980528211227F62A2fa6Ea71B4375":                                                                    true,
		"0x81B1573614F7382d0c17bcB6936cA86C44A6dd3A":                                                                    true,
		"0x0D46cB2116979BCb77ec0b9BeF99a3768384658f":                                                                    true,
		"0x23e175932E864fa17F22483F7b7BB94CDc97FF6D":                                                                    true,
		"0x065C84641EE62d032ea5F20C49d59817C87A2747":                                                                    true,
		"0xf3d0Cb1b3Db63A5f98a0EcB05d0D9F263058B7CB":                                                                    true,
		"0xD14E3d828a9B6C4623E525E2Dda7549f1AfF1828":                                                                    true,
		"0x0Ef8582381874780e4CDbbeaEf8Bfa1F9cd34DAe":                                                                    true,
		"0x13c2a9bD09858Dc078f66934Fb85dC96bE4bd668":                                                                    true,
		"0x99186F10Fb1C9821FA9955CE99d31c3454cF7d1A":                                                                    true,
		"0x701F9D1ccf4107e548Ee31693e746361129659fD":                                                                    true,
		"0xe96943DA5E2C74c0450612D5feB3537d036CEFC0":                                                                    true,
		"0xF3A66C660fa1A41f8FcC04504B506163c119552C":                                                                    true,
		"0x917Aa1FE4eE8154CE8c4EeB8f4768BA615245799":                                                                    true,
		"0x60b97528be1d072aC5b2405b6e292B43B3eb4AD0":                                                                    true,
		"0xB18278584bD3e41DB25453EE3c7DeDfc84040420":                                                                    true,
		"0x75bC526d9655a93c89a459185021b194a71B070C":                                                                    true,
		"0xfb1904dB9aEbc93738f0C6Dc03DCa7376cEEf07e":                                                                    true,
		"0x330AA0a042347313B68Be4CB629323488CF19D20":                                                                    true,
		"0x1bfc7b12b950A88371947EAFE123126ea3D5314A":                                                                    true,
		"0x45Ab04C54264F485eeF8DB0c20E531E9D37cD53A":                                                                    true,
		"0x77F5Fe97D52e4e3C736a4d26869Ec54f7Ae0b033":                                                                    true,
		"0xbF8Bd659Cc3062904098a038FF914Db67b7550D6":                                                                    true,
		"0xC00604A96e6Fa9f978E977124ad00638fEA52D0d":                                                                    true,
		"0xe2B397253Aea10eF34D2Ed41A28FDB4eA77cF8aC":                                                                    true,
		"0xE56C470a342eedbb2e001e70d1C5765C042EB349":                                                                    true,
		"0xF83abc519E046c5391d219fabF1A3C87dd5924D3":                                                                    true,
		"0x8A19d26887722c5073f90180a966a5BF88EF5210":                                                                    true,
		"0xEeC74Dcc37e01ef9405F52Ad1F8c4A01Ef6ab6cB":                                                                    true,
		"0x0D1580C763957b10529B2dbd770D884b06D4b246":                                                                    true,
		"0x92740bf77826b2a0B4D24947cD28404d9195399F":                                                                    true,
		"0x928ab1b6720FFCF398FF2DeC47C44916a7Bb7AA5":                                                                    true,
		"0xF7BFdFAb2765ec2E364FE6C9CC18f00111b6a26D":                                                                    true,
		"0x37091b108c1CF433E0E33805a8316470D3Ea01fF":                                                                    true,
		"0x5Cc2E85dF7af0409A422025a572bFb3ee4b9aFb6":                                                                    true,
		"0xF1B96fAf44D0a04F7A39ed90b4B3A2942403b109":                                                                    true,
		"0x7ce2E4c550f81f0352c12b2b961FC59324b4d025":                                                                    true,
		"0xc1B206BD2d9fe92aB61C9ca838C977bE0CbaDd7c":                                                                    true,
		"0xBC689f7292ae3f08ff8374bC973b9F79ce0A7144":                                                                    true,
		"0x8044F81869F5Fb9AEa222811f7993d5075bb3B92":                                                                    true,
		"0xCeFf445cb5A6016103C7b24dC8cE30C4692E1C2C":                                                                    true,
		"0x7143aCD5d3D933ad92177F441Bf5ebd949605e31":                                                                    true,
		"0x64C5480bf1BA9C1E7Fc3695a349E26446d777E7f":                                                                    true,
		"0xc6B968376B7c48788D59b2610663E9913BBf132E":                                                                    true,
		"0x8A4Cf3A98D5D86BbdA7bC522Dc7F614e1fa8bA20":                                                                    true,
		"0xDbbEdFb2f0bfB0b71baB63708A1a3C5a1509f37A":                                                                    true,
		"0xBE9594815B8C43a89162bAbA9bcd1943B51Dc05e":                                                                    true,
		"0x9DD8Fb2d577B11Ae20C33289cDAA91D27eA937A6":                                                                    true,
		"0x82530e41ea7c0d3cb0279e70dd5abda176ffd93b":                                                                    true,
		"0xAb3d2d2cf939DbEa55e63e7f2206636535DEC7E9":                                                                    true,
		"0xd00F74Ae1497bB132B52F58AB8B53970D55EDC4E":                                                                    true,
		"0x2CE7b28f1D6080C091009BE9efFa19eA6C04adEA":                                                                    true,
		"0xaFc699A932307d063940c639b072dfE4F040786e":                                                                    true,
		"0x995399a28F951Ec14a3fca942823066ac0Abc5F8":                                                                    true,
		"0x6a4a5B7F9913575A1006ACd77F6CfD6CC97df091":                                                                    true,
		"0xbAbaCd48A04d4459c220Fd25c4F5cbA7970D93D3":                                                                    true,
		"0x9bC1cf05830382D0314921CE93d744867669b10C":                                                                    true,
		"0x8e2cB1470BeFd8D2a861Dd3Cdbda8908490002B6":                                                                    true,
		"0xa694858857a0D50Da11A8176FDCECc3F96964968":                                                                    true,
		"0x8CC200570E368aF3e01793aAD4067E4654D1409A":                                                                    true,
		"0x2F9622d8c2B3a68cE6598D1AE231681AC75D71e6":                                                                    true,
		"0xeb15a2b191a982845f80eb696d931d05385a54bc":                                                                    true,
		"0xa1622c2c356b8a23174a1603f457cec28158afd8":                                                                    true,
		"0xbcacfd0e6ee95de4b37681d34e53dbf377d0f208":                                                                    true,
		"0x93cc4f9c4ddeef142d7aac607f3cb833af3a5800":                                                                    true,
		"0xe6dce1b5FCc9Fc48fDe98fb40ae1899BDB4c4cA1":                                                                    true,
		"0x7a4A103A5C7E07C0158075cd45E1C47ACD1BBC6B":                                                                    true,
		"0x4Dbc72984DB81143D2f939766f35Cc94Da292Bc3":                                                                    true,
		"0x7f0d1690437e0ecde0999d2a572aba7e4bccb15a":                                                                    true,
		"0xca12718015f266bb78e9597d555f014f72f76842":                                                                    true,
		"0x9fa440137565c9b11f5610d6c0eb6d305e8831a0":                                                                    true,
		"0x8201756E2008304faF0c83C1e693FC50c7CE7909":                                                                    true,
		"0x22587cd068e05d845c108205c1c698e882d407a1":                                                                    true,
		"0xdd312ecb79b28008f9861e364e5ae211cabe8527":                                                                    true,
		"0xF32E0558e639ACe32dEEb5e11fe4d0510Fb28887":                                                                    true,
		"0xF7c96fA35B0162B8637d8849773C8C0217B38F46":                                                                    true,
		"0xdccDdB7432B7371745233410C943F4882Df07937":                                                                    true,
		"0xEb62143E84d138C81c93B32855b735ccc8A1eb1D":                                                                    true,
		"0xc8b8dadb02e3b5c8bb49f70e498aa729f1def845":                                                                    true,
		"0x40bf0275c7bb44e0def1a3d4c651b1863dd94d96":                                                                    true,
		"0x03630a389CeD4b57cCDb7954bDfc656424B6cb66":                                                                    true,
		"0xc5345c212f140bb4aa2f493630d6fb6ec5de4174":                                                                    true,
		"0xd2f26af8df03be22b8608ca67d9f3178c0dd0591":                                                                    true,
		"0x508688a097d23a14a8782ab3a3d47b4141f97c5a":                                                                    true,
		"0x685678fe8Dabf9a1bDA9ab06C8eB2771CC44E98b":                                                                    true,
		"0xa245c29985153e076326267b3A26fA0B6EEA4EBA":                                                                    true,
		"0xf4a993362b8af739b030b3a52834b9942de4b664":                                                                    true,
		"0xdd2614a63ecffb220aa73c439af1a3df1848c98f":                                                                    true,
		"0x011d46ac42bd2e61138c1439e886f0ac1499140f":                                                                    true,
		"0xfe9d82cefc9e44694e5a21e8c1b5bc27670e6925":                                                                    true,
		"0x68a9360E07a5fe96a2209A64Fa486bB7B2dF217B":                                                                    true,
		"0xcd08a15b4032f4165620e3346f4539cf0986fe64":                                                                    true,
		"0x6542f55754ee2e10bfc5fec44aa0ac1530a620a8":                                                                    true,
		"0x8966141754fdf90F4e8bBBB5dC384aa4F14048f9":                                                                    true,
		"0xD577002B765e048Fda0b64fAd500c9B2Cb6fA2E4":                                                                    true,
		"0x0275BEaf895813C7913D3e301D69515CE93fEB9d":                                                                    true,
		"0x8FdDE47b396B437Bb23282fE6C19ea330B6525Ff":                                                                    true,
		"0x5ae60FD8Adc5990Ce91fe2066ba603F55e484a10":                                                                    true,
		"0xA7B70f9B5224FB81f69F7E413f3a682DAe1AA936":                                                                    true,
		"0x0d352FAa7ba9003F772589e15A39f3aB80e88ABB":                                                                    true,
		"0x0fB8B55EE84F2D59F92AdAb1A163F6787a2A3147":                                                                    true,
		"0xCE64998009B70710d8E2a422A22c87524836667F":                                                                    true,
		"0xB7e2163B711692cbc6224646538b42E651Cd085a":                                                                    true,
		"0x2B60448C6CCd8BC79700D2b2f0632B0c855b357b":                                                                    true,
		"0x9DE97b63B677d4856dE8B4e8F2F1a55d8787Eb71":                                                                    true,
		"0x981741Bdd9c413D45E1f7F1F0c1300F5469FED87":                                                                    true,
		"0x041f21a9DF5c6116fD1ACBD99Ee928cd69f497eB":                                                                    true,
		"0xe18e10bb8740836B52077221bcA8d601F4196aa7":                                                                    true,
		"0xCfC64CE7D1fdBCab9D846EC5ca1847CC0D828753":                                                                    true,
		"0x310dD78dD87353Cb66A26fb2FD724eB95F4058Dc":                                                                    true,
		"0xFCfD06C217b0C1bf5b6BdbDe7FE0C15c21201DDd":                                                                    true,
		"0x38579412DD68870F340c840b4eadcd511Bb87FB6":                                                                    true,
		"0x44eBB53349577a4B48954267930CDF9e08B065EC":                                                                    true,
		"0x13D91dd132cF416a3293221B4786Ec75c7504A6e":                                                                    true,
		"0xecfbddc9fd292c23979c1b298eb839e7dc84f868":                                                                    true,
		"0xF7f78de8ef7534bC76F1360Cf67d7661c9427db3":                                                                    true,
		"0x67B368d90D50d618d66818671d8Dd02263875712":                                                                    true,
		"0x703A189F5794085f4dfb8016aC0d81bF076597E6":                                                                    true,
		"0x4def593b79ca25f40d1757a8ebdcc9b35f2351ec":                                                                    true,
		"0x4a29367c5ae9f84ef03e447d1f7dee8e6b16229d":                                                                    true,
		"0x7462aCdEDdFBE39540E3654d4c81c04E41A9E784":                                                                    true,
		"0xeb6a9f6b820c3fe5deee4e586033b1b800699662":                                                                    true,
		"0xed17F05A7dD17Dd05a60C01099f1c7B5480f7Af9":                                                                    true,
		"0x57b33f88d6bd74e2d5008209fa7b3d812387ee07":                                                                    true,
		"0xB8A89619Ddb9a9aAf41C12Cf135D48759Ae483d6":                                                                    true,
		"0x51494ff727ce562135b58505415dfb548df4bd50":                                                                    true,
		"0xaa9CAF94a514bA21ceaef898024b870C8B96EdF8":                                                                    true,
		"0x22c465eD9F1bBa7C8F5F9bD4282e9645BA25A262":                                                                    true,
		"0x38e037135Cf65959C7d39D5Ed6C42F2a54e5Df56":                                                                    true,
		"0x9e645A3a5790c95Fd975463E9404723cd53C66e5":                                                                    true,
		"0xC8Fb61bCa1DEab05aeD66a2E6b91a5a2078Bad8F":                                                                    true,
		"0xBe01bC7266Bdb41787e6e13B2B4804a83D177B72":                                                                    true,
		"0xA1eF6e844c432731e5Eae28832b5A7eE66209419":                                                                    true,
		"0x16E1258fD25683f8FAeDa7488B2a6e85cF413926":                                                                    true,
		"0xe814c24037C80B1De49CAB13F2Ab61c5E54C24c1":                                                                    true,
		"0x5446278D248e37f788a68E7a3A582B998605df06":                                                                    true,
		"0x43207ce2b040d554aAE23bbEe357F7c72227Fd9B":                                                                    true,
		"0x8D855947fD9c77c77Bcc737D1b9612b881dC7497":                                                                    true,
		"0x02218f563b4600c034e72fffeab0258bd03eacae":                                                                    true,
		"0x7cBA6c9F61BF1fF8902648040593c29927486a70":                                                                    true,
		"0x2d5076d56be5fe93e1bf8b3dbeec43dc915a9ab6":                                                                    true,
		"0x99b1DFbD043Da1202EEBA16e2b68040cba246d69":                                                                    true,
		"0x93206c3a76b87c015e0202b14a4c7af482ad6c5e":                                                                    true,
		"0xd428BBbD5A0eBEc47ebDBa71B6C72B1e4b013Ecc":                                                                    true,
		"0x3b2f53c130e5b7b4cc9b6947bb7fa73c59ec96a8":                                                                    true,
		"0xfD1803a7961e2302664fDB9390246dB0255C5A12":                                                                    true,
		"0x6594b44f810ef7f24bacf04cb487b62e06a34dac":                                                                    true,
		"0x9E0616Aa26cC7cE07dcA71F3800Be182f15879CE":                                                                    true,
		"0xFd7477FF27Cfe71bcF112445911E183e7D00f1e2":                                                                    true,
		"0x1bf555e99b9056e75f5752ebf6593c4929bf5d50":                                                                    true,
		"0x8C05B197f30723fcBf95B5B2A0a6e1e79d882AC1":                                                                    true,
		"0xC969CE1f930d6264e7B234a212B88f6AcaEB9027":                                                                    true,
		"0x5D38D1CEF7103212933f8d52249199E841DAe4E0":                                                                    true,
		"0x8CcBb046eEC43CE314782Ca18765c9a424089216":                                                                    true,
		"0xf073625f1248DB830Ca9d638EF55fe01f74172e2":                                                                    true,
		"0x9c9a76BBe192dCF6D08495939BE83664E1288315":                                                                    true,
		"0xE3Fb114BE211cf896b13D2f7783Cc1972487e6b9":                                                                    true,
		"0x62035f0D3F6eF33dC45b160b119F887233061197":                                                                    true,
		"0x35E000A0b0dAD478FBdfa7c617522F7e7C036BB9":                                                                    true,
		"0xa0A266736BE7a0f9aDA37F4660CE943847Ce24eE":                                                                    true,
		"0x8a90E9DECb493fD5521c4A8aDb6Dc385D14Ca8Ff":                                                                    true,
		"0xd42913216ACcf7ebF6DA4dBFb4F36f86Ce8eb86F":                                                                    true,
		"0x2106dF7C24125C1C8f7F1F6DEAaEDe93Fd9a2E74":                                                                    true,
		"bc1p59emk2dcfdu7mff30cg5wvc98lzwu7d85xrg4erkathwk686d9fspqvn2p":                                                true,
		"0xb1b7c542196d532dcdE45FbAbF1AaAeDab620357":                                                                    true,
		"0x954C8952CFB768617141A2CBeB2b23FAAf7b89Cf":                                                                    true,
		"0x9f341c18fbdc6802e09b3ecf97cb040591b80e5d":                                                                    true,
		"0x7cE0BA3Bbe9c690dc6eB210e9C8223BB4594915A":                                                                    true,
		"0x0a0D11db196C1770562214088DcEC6DdBE368d4d":                                                                    true,
		"0x9cBc43AeB5B537B11c11C333705E3155480833c4":                                                                    true,
		"0x2eAc439143437Ef9241867A2d2CEaAc2b95C00c7":                                                                    true,
		"0xC035f22499c26a0B2C71A3dD155CD4b131540AAb":                                                                    true,
		"0x8A7Bc8EF535760570c3A99fa79Ed1A8bCea26861":                                                                    true,
		"0xDCa23632231a1a27E56F9f4A56fbc5558A46D6D8":                                                                    true,
		"0xbCeac80526dB65D439988CB74eAb6093DdB53364":                                                                    true,
		"0x7508d0d98f094bb1facce7d13114bebc7dbf8e3b":                                                                    true,
		"bc1pl4wg5r6fea9qnkjrs7pyg9w76q4kqzxyn4rnfgxzp7cpe94rldhsuc306p":                                                true,
		"bc1p397k4sppa59v9q6uzu0qrqlwhla4tugm7zzejyc7utpqxxae3jxss24tg8":                                                true,
		"bc1p8eecj5z2w99n6sd7whc6qdz8lek8c50z2s9c9q5nleuzx4vdn95qq2x50k":                                                true,
		"bc1p88gd82wu55358tgxr63ryl727vwtaeemettj5jmp6wz5v5ykuujsnms4e9":                                                true,
		"bc1pmpqjrxkfwgqggk3tkngagzls2q478tr3zckpwsucuc7kz4frpq4suzlpmm":                                                true,
		"0x4C880bF84229494159882d469050aF4D323B5Df0":                                                                    true,
		"0x53f2ba1B842Bf7500B2B7Dd767c89bC356417BcD":                                                                    true,
		"0x02b71643925B74BfF44fc42c701e5299Dfd262Bb":                                                                    true,
		"0xE7F3443B16648A89dd3eAf9D179CB4090C02E66F":                                                                    true,
		"0xa18c029b6310bA96d6F9e125EE6Bb7FD48D4cD1f":                                                                    true,
		"0xD827084f06670B692A66043b3675C57430c3e630":                                                                    true,
		"0x567CB330Bd9875A6dFB6EAeFF1eBa80ab141e935":                                                                    true,
		"0x766acD4F63FB51be0AAF1Ee8525616dF2f0b766d":                                                                    true,
		"0x5F8A23905d4C467109D1188f5F6900e7f05eB5E6":                                                                    true,
		"0xc675628d1BDb07B916f9ffC8C5434A8B18DaBC73":                                                                    true,
		"0x505836DE537E3B8d28d092449c6Bf6F7A3B80c30":                                                                    true,
		"0x73D05c2Ea70dFC3B220444c94567dbc84Bb0d24C":                                                                    true,
		"0xBDa3735c8BF15Eea9520CF165666E0CBcB13134A":                                                                    true,
		"0x32CC128F35aD103041ECCb742Bb36e36Cf9158e4":                                                                    true,
	}
	_, ok = manual[address]
	if ok {
		return 20.0
	}

	allow, err := u.Repo.GetProjectAllowList("999998", address)
	if err == nil && allow.UUID != "" {
		return 10.0
	}

	return 0.0
}

func (u Usecase) GetPriceCoinBase(coinID int) (*coin_market_cap.PriceConversionResponse, error) {
	key := fmt.Sprintf("gm-collections.coinbase.price.rate." + strconv.Itoa(coinID))
	cached, err := u.Cache.GetData(key)
	result := &coin_market_cap.PriceConversionResponse{}
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			return result, err
		}
	}
	result, err = u.CoinMarketCap.PriceConversion(coinID)
	if err == nil {
		u.Cache.SetDataWithExpireTime(key, result, 60*30)
	}
	return result, nil
}

func (u Usecase) GetBitcoinBalance(addr string) (*structure.BlockCypherWalletInfo, error) {
	key := fmt.Sprintf("gm-collections.quicknode.bitcoin.balance" + addr)
	result := &structure.BlockCypherWalletInfo{}

	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			return result, err
		}
	}

	result, err = btc.GetBalanceFromQuickNode(addr, u.Config.QuicknodeAPI)
	if err == nil {
		u.Cache.SetDataWithExpireTime(key, result, 60*5)
	}
	return result, nil
}

func (u Usecase) ClearCacheTop10GMDashboard() {
	dashboardItems, err := u.GetChartDataForGMCollection(true)
	if err != nil {
		logger.AtLog.Logger.Error("ClearCacheTop10", zap.Error(err))
		return
	}

	items := dashboardItems.Items
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].UsdtValue > items[j].UsdtValue
	})

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("gm-collections.deposit.eth3.gmAddress." + items[i].From + "." + items[i].To)
		err := u.Cache.Delete(key)
		if err != nil {
			logger.AtLog.Logger.Error("ClearCacheTop10", zap.Error(err), zap.String("cachedKey", key), zap.Any("item", items[i]))
		}

		keyErc2 := fmt.Sprintf("gm-collections.deposit.erc20_1.gmAddress." + items[i].From + "." + items[i].To)
		err = u.Cache.Delete(keyErc2)
		if err != nil {
			logger.AtLog.Logger.Error("ClearCacheTop10", zap.Error(err), zap.String("cachedKey", key), zap.Any("item", items[i]))
		}
	}
}

func (u Usecase) BackupGMDashboardCachedData() {
	key := fmt.Sprintf("gm-collections.deposit")
	cached, err := u.Cache.GetData(key)
	if err == nil && cached != nil {
		dataEntity := &entity.CachedGMDashBoard{
			Value: cached,
			Key:   key,
		}

		dataEntity.SetID()
		dataEntity.SetCreatedAt()

		inserted, err := u.Repo.Create(context.Background(), dataEntity.TableName(), dataEntity, nil)
		if err != nil {
			logger.AtLog.Logger.Error("BackupGMDashboardCachedData", zap.Error(err), zap.String("key", key))
			return
		}

		logger.AtLog.Logger.Info("BackupGMDashboardCachedData", zap.String("key", key), zap.Any("inserted", inserted))
		return
	}

	logger.AtLog.Logger.Error("BackupGMDashboardCachedData", zap.Error(err), zap.String("key", key))
}

func (u Usecase) RestoreGMDashboardCachedData(UUID string) {
	key := fmt.Sprintf("gm-collections.deposit")
	cached, err := u.Repo.FindOne(entity.CachedGMDashBoard{}.TableName(), UUID)
	if err != nil {
		return
	}

	data := &entity.CachedGMDashBoard{}
	err = helpers.Transform(cached, data)
	if err != nil {
		return
	}

	resp := &structure.AnalyticsProjectDeposit{}
	bytes, err := json.Marshal(data.Value)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return
	}

	u.Cache.SetDataWithExpireTime(key, resp, 60*60*1)
}

func (u *Usecase) SendGMMEssageToSlack(preText string, content string) {
	slackChannel := os.Getenv("SLACK_GM_DASHBOARD")
	if slackChannel == "" {
		slackChannel = os.Getenv("SLACK_ALLOW_LIST_CHANNEL")
	}

	//send message to slack
	title := ""

	if _, _, err := u.Slack.SendMessageToSlackWithChannel(slackChannel, preText, title, content); err != nil {
		logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
	}
}
