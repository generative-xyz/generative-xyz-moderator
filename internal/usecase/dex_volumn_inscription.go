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
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
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
	key := fmt.Sprintf("gm-collections.deposit.erc20.gmAddress." + tcAddress + "." + gmAddress)
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
	time.Sleep(time.Millisecond * 250)
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
		u.Cache.SetDataWithExpireTime(key, resp, 24*60*60) // cache by 1 day
		return resp, nil
	}
	return nil, errors.New("not balance - " + gmAddress)
}

func (u Usecase) GetChartDataEthForGMCollection(tcAddress string, gmAddress string, transferedETH []string, oldData bool, ens string, avatar string) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.eth2.gmAddress." + tcAddress + "." + gmAddress)
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
	time.Sleep(time.Millisecond * 250)
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
			if len(transferedETH) > 0 {
				for _, v := range transferedETH {
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(v, 18)), ethRate)
				}
			}
			items = append(items, &etherscan.AddressTxItemResponse{
				From:      tcAddress,
				To:        gmAddress,
				Value:     moralisEthBL.Balance,
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

		cachedExpTime := 24 * 60 * 60

		if oldData {
			cachedExpTime = 30 * 24 * 60 * 60 //a month
		}
		u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime) // cache by 1 day
		return resp, nil
	} else {
		transferUsdtValue := float64(0)
		if len(transferedETH) > 0 && !oldData {
			for _, v := range transferedETH {
				transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(v, 18)), ethRate)
			}
			var items []*etherscan.AddressTxItemResponse
			items = append(items, &etherscan.AddressTxItemResponse{
				From:      tcAddress,
				To:        gmAddress,
				Value:     moralisEthBL.Balance,
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

			cachedExpTime := 24 * 60 * 60

			if oldData {
				cachedExpTime = 30 * 24 * 60 * 60 //a month
			}
			u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime) // cache by 1 day
			return resp, nil
		}
	}
	return nil, errors.New("not balance - " + gmAddress)
}

func (u Usecase) GetChartDataBTCForGMCollection(tcWallet string, gmWallet string, transferedBTC []string, oldData bool) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.btc3.gmAddress." + tcWallet + "." + gmWallet)
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
			if len(transferedBTC) > 0 {
				for _, v := range transferedBTC {
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(v, 8)), btcRate)
				}
			}

			item := &etherscan.AddressTxItemResponse{
				From:      tcWallet,
				To:        gmWallet,
				Value:     fmt.Sprintf("%d", walletInfo.Balance),
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
			u.Cache.SetDataWithExpireTime(key, resp1, 6*60*60) // cache by 6 hours
			return resp1, nil
		} else {
			transferUsdtValue := float64(0)
			if len(transferedBTC) > 0 {
				for _, v := range transferedBTC {
					transferUsdtValue += utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(v, 8)), btcRate)
				}
				item := &etherscan.AddressTxItemResponse{
					From:      tcWallet,
					To:        gmWallet,
					Value:     fmt.Sprintf("%d", walletInfo.Balance),
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
				u.Cache.SetDataWithExpireTime(key, resp1, 6*60*60) // cache by 6 hours
				return resp1, nil
			}
		}
		return nil, errors.New("not balance - " + gmWallet)
	}
}

func (u Usecase) JobGetChartDataForGMCollection111() error {
	_, err := u.GetChartDataForGMCollection(false)
	if err != nil {
		log.Println("JobGetChartDataForGMCollection GetChartDataForGMCollection err", err)
	}
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
				"0x2c7aFd015A4080C835139E94D0f624bE552b9c66": "0x2c7aFd015A4080C835139E94D0f624bE552b9c66",
				"0x46Ad79eFd29B4212eE2dB32153c682Db06614Ce5": "wwf88.eth",
				"0xD78D4be39B0C174dF23e1941aC7BA3e8E2a6b3B6": "0xD78D4be39B0C174dF23e1941aC7BA3e8E2a6b3B6",
				"0xBFB9AC25EBC9105c2e061E7640B167c6150A7325": "littlered.eth",
				"0xa3017BB12fe3C0591e5C93011e988CA4b45aa1B4": "0xa3017BB12fe3C0591e5C93011e988CA4b45aa1B4",
				"0xa3EEE445D4DFBBc0C2f4938CB396a59c7E0dE526": "0xa3EEE445D4DFBBc0C2f4938CB396a59c7E0dE526",
				"0xEAcDD6b4B80Fcb241A4cfAb7f46e886F19c89340": "0xEAcDD6b4B80Fcb241A4cfAb7f46e886F19c89340",
				"0x7729A5Cfe2b008B7B19525a10420E6f53941D2a4": "trappavelli.eth",
				"0x4bF946271EEf390AC8c864A01F0D69bF3b858569": "0x4bF946271EEf390AC8c864A01F0D69bF3b858569",
				"0x21668e3B9f5Aa2a3923E22AA96a255fE8d3b9aac": "0x21668e3B9f5Aa2a3923E22AA96a255fE8d3b9aac",
				"0x597c32011116c94994619Cf6De15b3Fdc061a983": "0x597c32011116c94994619Cf6De15b3Fdc061a983",
				"0xB18278584bD3e41DB25453EE3c7DeDfc84040420": "0xB18278584bD3e41DB25453EE3c7DeDfc84040420",
				"0xfA9A55607BF094f991884f722b7Fba3A76687e40": "0xfA9A55607BF094f991884f722b7Fba3A76687e40",
				"0xCa2b4ad56a82bc7F8c5A01184A9D9c341213e0d3": "0xCa2b4ad56a82bc7F8c5A01184A9D9c341213e0d3",
				//"0xfA9A55607BF094f991884f722b7Fba3A76687e40":,
				"0x63cBF2D7cf7EF30b9445bEAB92997FF27A0bcc70": "0x63cBF2D7cf7EF30b9445bEAB92997FF27A0bcc70",
				"0x64BE8226638fdF2f85D8E3A01F849E0c47AE9446": "0x64BE8226638fdF2f85D8E3A01F849E0c47AE9446",
				"0xbf22409c832E944CeF2B33d9929b8905163Ae5d4": "0xbf22409c832E944CeF2B33d9929b8905163Ae5d4",
				"0xda9979247dC98023C0Ff6A59BC7C91bB627d4934": "0xda9979247dC98023C0Ff6A59BC7C91bB627d4934",
				"0x9c0Da3467AeD02e49Fe051104eFb2255C2982C61": "0x9c0Da3467AeD02e49Fe051104eFb2255C2982C61",
				"0xCd2b27C0dc8db90398dB92198a603e5D5D0d5e30": "0xCd2b27C0dc8db90398dB92198a603e5D5D0d5e30",
				"0xe9084DEDfcD06E63Dc980De1464f7786e2690c82": "0xe9084DEDfcD06E63Dc980De1464f7786e2690c82",
			}
			for wallet, ens := range fromWallets {
				temp, err := u.GetChartDataEthForGMCollection(wallet, gmAddress, []string{}, true, ens, "")
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
			for _, item := range result.Items {
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
				} else {
					result.MapItems[item.From].UsdtValue += item.UsdtValue
					if item.Avatar != "" {
						result.MapItems[item.From].Avatar = item.Avatar
					}
					if item.ENS != "" {
						result.MapItems[item.From].ENS = item.ENS
					}
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

		u.Cache.SetDataWithExpireTime(key, result, 60*30)
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

	// TODO kll
	kll := map[string]bool{
		"": true,
	}
	_, ok := kll[address]
	if ok {
		return 25.0
	}

	tcBalance, err := u.TcClientPublicNode.GetBalance(context.TODO(), address)
	if err == nil && tcBalance.Cmp(big.NewInt(0)) > 0 {
		return 20.0
	}

	// TODO manual
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
