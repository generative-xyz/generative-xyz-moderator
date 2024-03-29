package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
	"os"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/request"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
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

func (u Usecase) GetChartDataERC20ForGMCollection(newcity entity.NewCityGm, transferedETH []string, ens string, avatar string) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.erc20_1.gmAddress." + newcity.UserAddress + "." + newcity.Address)
	result := &structure.AnalyticsProjectDeposit{}
	if newcity.UpdatedAt != nil {
		if time.Now().Add(time.Minute * -90).Before(*newcity.UpdatedAt) {
			u.Cache.Delete(key)
		}
	} else {
		if time.Now().Add(time.Minute * -90).Before(*newcity.CreatedAt) {
			u.Cache.Delete(key)
		}
	}

	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			logger.AtLog.Logger.Info(fmt.Sprintf("GetChartDataERC20ForGMCollection get from cached for userAddr %s wallet %s", newcity.UserAddress, newcity.Address),
				zap.Any("result", result), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

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
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection PriceConversion PEPE", zap.Any("err", err))
		} else {
			if pRate != nil {
				pepeRate = pRate.Data.Quote.USD.Price
				u.Cache.SetDataWithExpireTime(keypepeRate, pepeRate, 30*60)
			}
		}
	}

	keyturboRate := fmt.Sprintf("gm-collections.deposit.turboRate.rate")
	var turboRate float64 = 0
	cachedTURBORate, err := u.Cache.GetData(keyturboRate)
	if err == nil {
		turboRate, _ = strconv.ParseFloat(*cachedTURBORate, 64)
	}
	if turboRate == 0 {
		tRate, err := u.CoinMarketCap.PriceConversion(24911) //TURBO ID
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection PriceConversion TURBO", zap.Any("err", err))
		} else {
			if tRate != nil {
				turboRate = tRate.Data.Quote.USD.Price
				u.Cache.SetDataWithExpireTime(keyturboRate, turboRate, 30*60)
			}
		}
	}

	pepe := "0x6982508145454ce325ddbe47a25d4ec3d2311933"
	turbo := "0xa35923162c49cf95e6bf26623385eb431ad920d3"
	moralisERC20BL, err := u.MoralisNft.TokenBalanceByWalletAddress(newcity.Address, []string{pepe, turbo})
	//time.Sleep(time.Millisecond * 250)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection err1111", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
		return nil, err
	}

	logger.AtLog.Logger.Info(fmt.Sprintf("GetChartDataERC20ForGMCollection moralis balance pepe, turbo for user %s wallet %s", newcity.UserAddress, newcity.Address), zap.Any("moralisERC20BL", moralisERC20BL), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

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
			From:      newcity.UserAddress,
			To:        newcity.Address,
			Value:     fmt.Sprintf("%f", totalPepe),
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
			From:      newcity.UserAddress,
			To:        newcity.Address,
			Value:     fmt.Sprintf("%f", totalTurbo),
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
		u.Cache.SetDataWithExpireTime(key, resp, 1*60*60)

		logger.AtLog.Logger.Info("GetChartDataERC20ForGMCollection len(items) > 0", zap.Any("result", resp), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
		return resp, nil

	} else {
		if newcity.UpdatedAt != nil {

			logger.AtLog.Logger.Info("GetChartDataERC20ForGMCollection len(items) = 0,  newcity.UpdatedAt != nil", zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address))

			if time.Now().Add(time.Hour * -12).After(*newcity.UpdatedAt) {
				// cache empty for inactive wallet
				resp := &structure.AnalyticsProjectDeposit{}
				err := u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
				if err != nil {
					logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection len(items) = 0,  newcity.UpdatedAt != nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
				}
			}
		} else {

			logger.AtLog.Logger.Info("GetChartDataERC20ForGMCollection len(items) = 0,  newcity.UpdatedAt == nil", zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

			if newcity.Status == 1 && time.Now().Add(time.Hour*-12).After(*newcity.CreatedAt) {
				// cache empty for inactive wallet
				resp := &structure.AnalyticsProjectDeposit{}
				err := u.Cache.SetDataWithExpireTime(key, resp, 12*60*60) // cache by 1 day
				if err != nil {
					logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection len(items) = 0,  newcity.UpdatedAt == nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
				}
			}
		}
	}
	return nil, errors.New("not balance - " + newcity.Address)
}

func (u Usecase) GetChartDataEthForGMCollection(newcity entity.NewCityGm, transferedETH []string, oldData bool, ens string, avatar string) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.eth3.gmAddress." + newcity.UserAddress + "." + newcity.Address)
	result := &structure.AnalyticsProjectDeposit{}
	if !oldData {
		if newcity.UpdatedAt != nil {
			if time.Now().Add(time.Minute * -90).Before(*newcity.UpdatedAt) {
				u.Cache.Delete(key)
			}
		} else {
			if time.Now().Add(time.Minute * -90).Before(*newcity.CreatedAt) {
				u.Cache.Delete(key)
			}
		}
	}
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			logger.AtLog.Logger.Info(fmt.Sprintf("GetChartDataEthForGMCollection get from cached for userAddr %s wallet %s", newcity.UserAddress, newcity.Address),
				zap.Any("result", result),
				zap.String("walletAddress", newcity.UserAddress),
				zap.String("gmAddress", newcity.Address),
				zap.String("key", key))

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
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection get price for eth", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address))
			return nil, err
		}
		u.Cache.SetDataWithExpireTime(keyRate, ethRate, 60*60) // cache by 1 hour
	}

	moralisEthBL, err := u.MoralisNft.AddressBalance(newcity.Address)
	//time.Sleep(time.Millisecond * 250)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err2222", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address))
		//return nil, err
		moralisEthBL = new(nfts.MoralisBalanceResp)
		temp, err := u.EtherscanService.AddressBalance(newcity.Address)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err3333",
				zap.Error(err),
				zap.String("gmAddress", newcity.Address))
			return nil, err
		}
		moralisEthBL.Balance = temp.Result
		if moralisEthBL.Balance == "" {
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err4444",
				zap.Error(err), zap.String("walletAddress", newcity.UserAddress),
				zap.String("gmAddress", newcity.Address))
			return nil, err
		}
	}

	//ethBL, err := u.EtherscanService.AddressBalance(gmAddress)
	//time.Sleep(time.Millisecond * 100)
	//if err != nil {
	//	return nil, err
	//}

	logger.AtLog.Logger.Info(fmt.Sprintf("GetChartDataEthForGMCollection moralisEthBL for user %s wallet %s", newcity.UserAddress, newcity.Address),
		zap.Any("moralisEthBL", moralisEthBL),
		zap.String("walletAddress", newcity.UserAddress),
		zap.String("gmAddress", newcity.Address),
		zap.String("key", key))

	totalEth := utils.GetValue(moralisEthBL.Balance, 18)
	if totalEth > 0 {
		usdtValue := utils.ToUSDT(fmt.Sprintf("%f", totalEth), ethRate)

		var items []*etherscan.AddressTxItemResponse
		if oldData {
			// get tx by addr
			ethTx, err := u.EtherscanService.AddressTransactions(newcity.Address)
			time.Sleep(time.Millisecond * 100)
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataEthForGMCollection", zap.Error(err), zap.String("gmAddress", newcity.Address))
				return nil, err
			}
			counting := 0
			for _, item := range ethTx.Result {
				if strings.ToLower(item.From) != strings.ToLower(newcity.UserAddress) {
					continue
				}
				items = append(items, &etherscan.AddressTxItemResponse{
					From:      newcity.UserAddress,
					To:        newcity.Address,
					Value:     fmt.Sprintf("%f", utils.GetValue(item.Value, 18)),
					UsdtValue: utils.ToUSDT(fmt.Sprintf("%f", utils.GetValue(item.Value, 18)), ethRate),
					Currency:  string(entity.ETH),
					ENS:       ens,
					Avatar:    avatar,
				})
				counting++
			}
			if counting == 0 {
				return nil, errors.New("not balance - " + newcity.Address)
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
				From:      newcity.UserAddress,
				To:        newcity.Address,
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

		cachedExpTime := 12 * 60 * 60

		if oldData {
			cachedExpTime = 30 * 24 * 60 * 60 //a month
		}
		err := u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err7777", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address))
		}

		logger.AtLog.Logger.Info("GetChartDataEthForGMCollection totalEth > 0", zap.Any("result", resp), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

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
				From:      newcity.UserAddress,
				To:        newcity.Address,
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

			cachedExpTime := 12 * 60 * 60 // cache by 1 hour

			if oldData {
				cachedExpTime = 30 * 24 * 60 * 60 //a month
			}
			err := u.Cache.SetDataWithExpireTime(key, resp, cachedExpTime)
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataEthForGMCollection err8888", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address))
			}

			logger.AtLog.Logger.Info("GetChartDataERC20ForGMCollection totalEth == 0", zap.Any("result", resp), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
			return resp, nil
		}
	}

	if newcity.UpdatedAt != nil {
		if time.Now().Add(time.Hour * -12).After(*newcity.UpdatedAt) {
			// cache empty for inactive wallet
			resp := &structure.AnalyticsProjectDeposit{}
			err := u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataEthForGMCollection  newcity.UpdatedAt != nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
			}
		}
	} else {
		if newcity.Status == 1 && time.Now().Add(time.Hour*-12).After(*newcity.CreatedAt) {
			// cache empty for inactive wallet
			resp := &structure.AnalyticsProjectDeposit{}
			err := u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataEthForGMCollection  newcity.UpdatedAt == nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
			}
		}
	}
	return nil, errors.New("not balance - " + newcity.Address)
}

func (u Usecase) GetChartDataBTCForGMCollection(newcity entity.NewCityGm, transferedBTC []string, oldData bool) (*structure.AnalyticsProjectDeposit, error) {
	// try from cache
	key := fmt.Sprintf("gm-collections.deposit.btc4.gmAddress." + newcity.UserAddress + "." + newcity.Address)
	result := &structure.AnalyticsProjectDeposit{}
	if !oldData {
		if newcity.UpdatedAt != nil {
			if time.Now().Add(time.Minute * -120).Before(*newcity.UpdatedAt) {
				u.Cache.Delete(key)
			}
		} else {
			if time.Now().Add(time.Minute * -120).Before(*newcity.CreatedAt) {
				u.Cache.Delete(key)
			}
		}
	}
	cached, err := u.Cache.GetData(key)
	if err == nil {
		err = json.Unmarshal([]byte(*cached), result)
		if err == nil {
			logger.AtLog.Logger.Info("GetChartDataBTCForGMCollection cached", zap.Any("result", result), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
			return result, nil
		}
	}

	// try from cache
	keyRate := fmt.Sprintf("gm-collections.deposit.btc.rate")
	var btcRate float64
	cachedETHRate, err := u.Cache.GetData(keyRate)
	if err == nil {
		btcRate, err = strconv.ParseFloat(*cachedETHRate, 64)
		if err != nil {
			btcRate = 0
		}
		logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection btcRate cached", zap.Float64("btcRate", btcRate))
	}
	if btcRate == 0 {
		btcRate, err := helpers.GetExternalPrice(string(entity.BIT))
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection btcRate", zap.Error(err))
			return nil, err
		}

		logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection btcRate GetExternalPrice", zap.Float64("btcRate", btcRate))
		u.Cache.SetDataWithExpireTime(keyRate, btcRate, 60*60) // cache by 1 hour
	}

	analyticItems := []*etherscan.AddressTxItemResponse{}
	if oldData {
		resp, err := u.MempoolService.AddressTransactions(newcity.Address)
		time.Sleep(time.Millisecond * 500)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection MempoolService", zap.Error(err))
			return nil, err
		}

		vouts := []mempool_space.AddressTxItemResponseVout{}
		for _, item := range resp {
			if item.Status.Confirmed {
				if oldData {
					isContinue := true
					for _, v := range item.Vin {
						if strings.ToLower(v.Prevout.Scriptpubkey_address) == strings.ToLower(newcity.UserAddress) {
							isContinue = false
						}
					}
					if isContinue {
						continue
					}
				}
				vs := item.Vout
				for _, v := range vs {
					if strings.ToLower(v.ScriptpubkeyAddress) == strings.ToLower(newcity.Address) {
						vouts = append(vouts, v)
					}
				}
			}
		}

		total := int64(0)
		for _, vout := range vouts {
			analyticItem := &etherscan.AddressTxItemResponse{
				From:      newcity.UserAddress,
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
			Value:        fmt.Sprintf("%f", amountF),
			Currency:     string(entity.BIT),
			CurrencyRate: btcRate,
			UsdtValue:    usdt,
			Items:        analyticItems,
		}
		u.Cache.SetDataWithExpireTime(key, resp1, 24*60*60*30) // cache by a month

		logger.AtLog.Logger.Info("GetChartDataBTCForGMCollection oldData", zap.Any("result", resp1), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

		return resp1, nil
	} else {
		/*_, bs, err := u.buildBTCClient()
		if err != nil {
			return nil, err
		}
		balance, confirm, err := bs.GetBalance(gmWallet)*/
		walletInfo, err := btc.GetBalanceFromQuickNode(newcity.Address, u.Config.QuicknodeAPI)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("GetChartDataBTCForGMCollection GetBalanceFromQuickNode - %s", newcity.Address), zap.String("newcity.Address", newcity.Address), zap.Error(err))
			return nil, err
		}
		time.Sleep(time.Millisecond * 50)

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

			item := &etherscan.AddressTxItemResponse{
				From:      newcity.UserAddress,
				To:        newcity.Address,
				Value:     fmt.Sprintf("%f", utils.GetValue(fmt.Sprintf("%d", walletInfo.Balance), 8)+transferBtcValue),
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
			err := u.Cache.SetDataWithExpireTime(key, resp1, 12*60*60) // cache by 2 hours
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection  walletInfo.Balance > 0", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
			}

			logger.AtLog.Logger.Info("GetChartDataBTCForGMCollection  walletInfo.Balance > 0", zap.Any("result", resp1), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

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
					From:      newcity.UserAddress,
					To:        newcity.Address,
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
				err := u.Cache.SetDataWithExpireTime(key, resp1, 12*60*60) // cache by 6 hours
				if err != nil {
					logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection  walletInfo.Balance <= 0", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
				}

				logger.AtLog.Logger.Info("GetChartDataBTCForGMCollection  walletInfo.Balance <= 0", zap.Any("result", resp1), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))

				return resp1, nil
			}
		}

		if newcity.UpdatedAt != nil {
			if time.Now().Add(time.Hour * -12).After(*newcity.UpdatedAt) {
				// cache empty for inactive wallet
				resp := &structure.AnalyticsProjectDeposit{}
				err := u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
				if err != nil {
					logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection  newcity.UpdatedAt != nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
				}
			}
		} else {
			if newcity.Status == 1 && time.Now().Add(time.Hour*-12).After(*newcity.CreatedAt) {
				// cache empty for inactive wallet
				resp := &structure.AnalyticsProjectDeposit{}
				err := u.Cache.SetDataWithExpireTime(key, resp, 3*60*60) // cache by 1 day
				if err != nil {
					logger.AtLog.Logger.Error("GetChartDataBTCForGMCollection  newcity.UpdatedAt == nil", zap.Error(err), zap.String("walletAddress", newcity.UserAddress), zap.String("gmAddress", newcity.Address), zap.String("key", key))
				}
			}
		}
		return nil, errors.New("not balance - " + newcity.Address)
	}
}

func (u Usecase) JobGetChartDataForGMCollection() error {
	u.Cache.SetStringData(utils.GM_CRONTAB_PROCESSING_KEY, "1")

	defer func() {
		u.Cache.SetStringData(utils.GM_CRONTAB_PROCESSING_KEY, "0")
	}()

	//clear cache for top 10 items
	//u.ClearCacheTop10GMDashboard()
	u.Logger.Info("JobGetChartDataForGMCollection - JobGetChartDataForGMCollection - GetChartDataForGMCollection: Start")

	//start
	now := time.Now().UTC()
	preText := fmt.Sprintf("[Analytics][Start] - Get chart data for GM Dashboard")
	content := fmt.Sprintf("Start at: %v", now)
	go u.SendGMMEssageToSlack(preText, content)

	data, err := u.GetChartDataForGMCollection(false)
	if err != nil {
		//end
		end := time.Now().UTC()
		preText = fmt.Sprintf("[Analytics][Error] - Get chart data for GM Dashboard")
		content = fmt.Sprintf("End at: %v with Err: %s", end, err.Error())
		go u.SendGMMEssageToSlack(preText, content)
		u.Logger.AtLog().Error("JobGetChartDataForGMCollection - GetChartDataForGMCollection: Error JobGetChartDataForGMCollection", zap.Error(err))
		return err
	}

	//end
	end := time.Now().UTC()
	preText = fmt.Sprintf("[Analytics][End] - Get chart data for GM Dashboard")
	content = fmt.Sprintf("End at: %v with USDT: %f, contributors: %d", end, data.UsdtValue, len(data.Items))
	go u.SendGMMEssageToSlack(preText, content)
	u.Logger.AtLog().Logger.Info("JobGetChartDataForGMCollection - GetChartDataForGMCollection: Complete JobGetChartDataForGMCollection", zap.Any("data", data.UsdtValue))

	return nil
}

func (u Usecase) GetListWallet(walletType string) ([]*structure.WalletResponse, error) {
	res := []*structure.WalletResponse{}
	wallets, err := u.Repo.FindNewCityGmByType(walletType)
	if err != nil {
		return nil, err
	}
	for _, v := range wallets {
		res = append(res, &structure.WalletResponse{
			UserAddress:  v.UserAddress,
			ENS:          v.ENS,
			Avatar:       v.Avatar,
			Address:      v.Address,
			Status:       v.Status,
			Type:         v.Type,
			NativeAmount: v.NativeAmount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}
	return res, nil
}

func (u Usecase) GetChartDataForGMCollection(useCaching bool) (*structure.AnalyticsProjectDeposit, error) {
	key := fmt.Sprintf(keyNotReAllocate)
	result := &structure.AnalyticsProjectDeposit{}

	if useCaching {
		// try get data from reAllocate, check config
		if os.Getenv("GetReallocateData") == "true" {
			dataReAllocate, err := u.GetReallocateData()
			if err == nil && dataReAllocate != nil {
				return dataReAllocate, nil
			}
		}
	}

	cached, err := u.Cache.GetData(key)
	if err == nil {
		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection: u.Cache.GetData", zap.String("key", key), zap.Bool("useCaching", useCaching))
	} else {
		u.Logger.AtLog().Logger.Error("GetChartDataForGMCollection: u.Cache.GetData", zap.String("key", key), zap.Bool("useCaching", useCaching), zap.Error(err))
	}

	if !useCaching || err != nil {
		if useCaching {

			u.Logger.AtLog().Logger.Error("GetChartDataForGMCollection: GetChartDataForGMCollection: useCaching",
				zap.String("key", key),
				zap.Bool("useCaching", useCaching), zap.Error(err))

			return nil, err
		}

		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection: GetChartDataForGMCollection: Starting ...")

		ethDataChan := make(chan structure.AnalyticsProjectDepositChan)
		btcDataChan := make(chan structure.AnalyticsProjectDepositChan)
		//erc20DataChan := make(chan structure.AnalyticsProjectDepositChan)

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
					// eth
					temp, err := u.GetChartDataEthForGMCollection(wallet, wallet.NativeAmount, false, wallet.ENS, wallet.Avatar)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny(fmt.Sprintf("GetChartDataForGMCollection: GetChartDataEthForGMCollection for wallet %s", wallet.UserAddress), zap.Any("err", err))
					}

					// erc20
					temp, err = u.GetChartDataERC20ForGMCollection(wallet, wallet.NativeAmount, wallet.ENS, wallet.Avatar)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny(fmt.Sprintf("GetChartDataForGMCollection: GetChartDataERC20ForGMCollection for wallet %s", wallet.UserAddress), zap.Any("err", err))
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
				temp, err := u.GetChartDataEthForGMCollection(entity.NewCityGm{UserAddress: strings.ToLower(wallet), Address: strings.ToLower(gmAddress)}, []string{}, true, ens, "")
				if err == nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
					data.CurrencyRate = temp.CurrencyRate
				}
				if err != nil {
					u.Logger.ErrorAny(fmt.Sprintf("GetChartDataForGMCollection: GetChartDataEthForGMCollection for wallet %s", wallet), zap.Any("err", err))
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
					temp, err := u.GetChartDataBTCForGMCollection(wallet, wallet.NativeAmount, false)
					if err == nil && temp != nil {
						data.Items = append(data.Items, temp.Items...)
						data.UsdtValue += temp.UsdtValue
						data.Value += temp.Value
						data.CurrencyRate = temp.CurrencyRate
					}
					if err != nil {
						u.Logger.ErrorAny(fmt.Sprintf("GetChartDataForGMCollection: GetChartDataBTCForGMCollection for wallet %s", wallet.UserAddress), zap.Any("err", err))
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
				temp, err := u.GetChartDataBTCForGMCollection(entity.NewCityGm{UserAddress: strings.ToLower(wallet), Address: strings.ToLower(gmAddress)}, []string{}, true)
				if err == nil && temp != nil {
					data.Items = append(data.Items, temp.Items...)
					data.UsdtValue += temp.UsdtValue
					data.Value += temp.Value
					data.CurrencyRate = temp.CurrencyRate
				}
				if err != nil {
					u.Logger.ErrorAny(fmt.Sprintf("GetChartDataForGMCollection: GetChartDataBTCForGMCollection for wallet %s", wallet), zap.Any("err", err))
				}
			}

		}(btcDataChan)

		/*go func(erc20DataChan chan structure.AnalyticsProjectDepositChan) {
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
					temp, err := u.GetChartDataERC20ForGMCollection(wallet, wallet.NativeAmount, wallet.ENS, wallet.Avatar)
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
		}(erc20DataChan)*/

		ethDataFromChan := <-ethDataChan
		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection: ethDataFromChan", zap.Any("ethDataFromChan", ethDataFromChan))

		btcDataFromChan := <-btcDataChan
		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection: btcDataFromChan", zap.Any("btcDataFromChan", btcDataFromChan))
		//erc20DataFromChan := <-erc20DataChan

		result := &structure.AnalyticsProjectDeposit{}
		if ethDataFromChan.Value != nil && len(ethDataFromChan.Value.Items) > 0 {
			u.Logger.Info("Processing data after go routine ethDataFromChan: ", zap.Int("ethDataFromChan", len(ethDataFromChan.Value.Items)))
			result.Items = append(result.Items, ethDataFromChan.Value.Items...)
			result.UsdtValue += ethDataFromChan.Value.UsdtValue
		}

		if btcDataFromChan.Value != nil && len(btcDataFromChan.Value.Items) > 0 {
			u.Logger.Info("Processing data after go routine btcDataFromChan: ", zap.Int("btcDataFromChan", len(btcDataFromChan.Value.Items)))
			result.Items = append(result.Items, btcDataFromChan.Value.Items...)
			result.UsdtValue += btcDataFromChan.Value.UsdtValue
		}

		/*if erc20DataFromChan.Value != nil && len(erc20DataFromChan.Value.Items) > 0 {
			result.Items = append(result.Items, erc20DataFromChan.Value.Items...)
			result.UsdtValue += erc20DataFromChan.Value.UsdtValue
		}*/

		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection - After routing", zap.Int("(result.Items", len(result.Items)))
		if len(result.Items) > 0 {
			u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection -Processing data after go routine", zap.Int("(result.Items", len(result.Items)))

			result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
			result.MapTokensDeposit = make(map[string][]structure.TokensDeposit)
			u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection - Processing data after go routine: build map")

			for i, item := range result.Items {
				u.Logger.AtLog().Logger.Info(fmt.Sprintf("GetChartDataForGMCollection - [%d / %d] Add item to map[string]item - start", i, len(result.Items)),
					zap.Any("item", item),
				)

				item.From = strings.ToLower(item.From)
				_, ok := result.MapItems[item.From]

				u.Logger.AtLog().Logger.Info(fmt.Sprintf("GetChartDataForGMCollection - [%d / %d] Add item to map[string]item: %s", i, len(result.Items), item.From),
					zap.Bool("ok", ok),
				)

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

			u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection - Processing data after go routine: rebuild items", zap.Int("items", len(result.Items)))

			usdtExtra := 0.0
			usdtValue := 0.0

			for i, item := range result.Items {
				u.Logger.AtLog().Logger.Info(fmt.Sprintf("GetChartDataForGMCollection - [%d / %d] Processing total usdtExtra for %s", i, len(result.Items), item.From))

				if item.UsdtValue < 0 {
					item.UsdtValue = 0.0
				}
				item.ExtraPercent = 0.0
				item.UsdtValueExtra = item.UsdtValue
				usdtExtra += item.UsdtValueExtra
				usdtValue += item.UsdtValue
			}

			for i, item := range result.Items {
				u.Logger.AtLog().Logger.Info(fmt.Sprintf("GetChartDataForGMCollection - [%d / %d] Processing GMReceive for %s", i, len(result.Items), item.From))

				item.Percent = item.UsdtValueExtra / usdtExtra * 100
				item.GMReceive = item.Percent * 8000 / 100
				item.GMReceiveString = fmt.Sprintf("%f", utils.ToWei(item.GMReceive, 18))
				if strings.Contains(item.GMReceiveString, ".") {
					item.GMReceiveString = strings.Split(item.GMReceiveString, ".")[0]
				}
			}
			result.UsdtValue = usdtValue

			u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection - End Processing data after go routine",
				zap.Int("contributors", len(result.Items)),
				zap.Float64("usdtValue", usdtValue),
			)
		}

		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection - Unmarshal old caching")
		cachedData := &structure.AnalyticsProjectDeposit{}
		if cached != nil {
			err := json.Unmarshal([]byte(*cached), cachedData)
			if err != nil {
				logger.AtLog.Logger.Error("GetChartDataForGMCollection json.Unmarshal.cachedData", zap.Error(err))
				return nil, err
			}
		}

		go u.BackupGMDashboardCachedData(cachedData, result)

		err = u.Cache.SetDataWithExpireTime(key, result, 60*60*24*3)
		if err != nil {
			u.Logger.AtLog().Logger.Error("GetChartDataForGMCollection: SetDataWithExpireTime ...", zap.Error(err))
		}

		u.Logger.AtLog().Logger.Info("GetChartDataForGMCollection: end! ...")
		return result, nil
	}

	err = json.Unmarshal([]byte(*cached), result)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection json.Unmarshal.cachedData", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (u Usecase) GetReallocateData() (*structure.AnalyticsProjectDeposit, error) {

	result := &structure.AnalyticsProjectDeposit{}
	keyRelocate := fmt.Sprintf(keyReAllocate)
	cachedRelocation, err := u.Cache.GetData(keyRelocate)
	if err == nil && cachedRelocation != nil {
		err = json.Unmarshal([]byte(*cachedRelocation), result)
		if err != nil {
			logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection json.Unmarshal.cachedData gm-collections.deposit.relocate", zap.Error(err))
			return nil, err
		}

		return result, nil
	}

	// get database
	dataFromDB, err := u.Repo.GetTheLatestReAllocated()
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection data from DB", zap.Error(err))
		return nil, err
	}

	url := dataFromDB.BackupFileName
	data, err := u.GCS.ReadFile(url)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection read file from GCS", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataERC20ForGMCollection read data from GCS", zap.Error(err))
		return nil, err
	}

	err = u.Cache.SetDataWithExpireTime(keyRelocate, result, 60*60*24*3)
	if err != nil {
		u.Logger.AtLog().Logger.Error("GetChartDataForGMCollection: SetDataWithExpireTime ...", zap.Error(err))
	}
	return result, nil
}

func (u Usecase) ReAllocateGM() (*structure.AnalyticsProjectDeposit, error) {
	u.Logger.Info("ReAllocateGM: get data from cache")
	key := fmt.Sprintf(keyNotReAllocate)
	result := &structure.AnalyticsProjectDeposit{}
	cached, err := u.Cache.GetData(key)
	//cached = &testData
	if cached == nil {
		logger.AtLog.Logger.Error("ReAllocateGM err json.Unmarshal.cachedData")
		return nil, err
	}
	err = json.Unmarshal([]byte(*cached), result)
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM err json.Unmarshal.cachedData", zap.Error(err))
		return nil, err
	}

	usdtExtra := 0.0
	usdtValue := 0.0
	u.Logger.AtLog().Logger.Info(fmt.Sprintf("ReAllocateGM - start with %d items", len(result.Items)))

	//move out of routine for prevent data race
	chanData := make(chan *etherscan.AddressTxItemResponse)
	for i, item := range result.Items {
		if item.UsdtValue < 0 {
			u.Logger.AtLog().Logger.Info(fmt.Sprintf("ReAllocateGM - maybe refund data, %s", item.From), zap.String("from", item.From), zap.Float64("usdtValue", item.UsdtValue))
			item.UsdtValue = 0.0
			item.UsdtValueExtra = 0.0
		}
		go func(i int, txItem *etherscan.AddressTxItemResponse, dataChan chan *etherscan.AddressTxItemResponse) {
			txItem.ExtraPercent = u.GetExtraPercent(txItem.From)
			if txItem.ExtraPercent == 0 {
				txItem.UsdtValueExtra = txItem.UsdtValue
			} else {
				txItem.UsdtValueExtra = txItem.UsdtValue/100*txItem.ExtraPercent + txItem.UsdtValue
			}

			u.Logger.AtLog().Logger.Info(fmt.Sprintf("ReAllocateGM - process extra percent for %s", txItem.From), zap.String("from", txItem.From), zap.Float64("ExtraPercent", txItem.ExtraPercent), zap.Float64("UsdtValueExtra", txItem.UsdtValueExtra))
			dataChan <- txItem

		}(i, item, chanData)

		if i%100 == 0 {
			time.Sleep(time.Millisecond * 250)
		}
	}

	resp := []*etherscan.AddressTxItemResponse{}

	//usdtExtra, usdtValue and the new updated items are created here
	//move usdtExtra, usdtValue for calculating total
	for {
		dataFromChan := <-chanData

		usdtExtra += dataFromChan.UsdtValueExtra
		usdtValue += dataFromChan.UsdtValue

		item := dataFromChan

		resp = append(resp, item)
		u.Logger.AtLog().Logger.Info(fmt.Sprintf("ReAllocateGM - calculing total usdtExtra, usdtValue: %s", item.From),
			zap.String("from", item.From),
			zap.Float64("usdtValue", item.UsdtValue),
			zap.Float64("usdtValueExtra", item.UsdtValueExtra),
			zap.Float64("extraPercent", item.UsdtValueExtra),
			zap.Int("items", len(resp)))

		if len(resp) == len(result.Items) {
			break
		}
	}

	//calculate via the updated item array (resp)
	totalGMReceive := float64(0)

	for i, item := range resp {
		item.Percent = item.UsdtValueExtra / usdtExtra * 100
		item.GMReceive = item.Percent * 8000 / 100
		item.GMReceiveString = fmt.Sprintf("%f", utils.ToWei(item.GMReceive, 18))
		if strings.Contains(item.GMReceiveString, ".") {
			item.GMReceiveString = strings.Split(item.GMReceiveString, ".")[0]
		}
		totalGMReceive += item.GMReceive
		u.Logger.AtLog().Logger.Info(fmt.Sprintf("ReAllocateGM - processed [%d / %d] - caculcating GM receive for: %s", i, len(resp), item.From),
			zap.String("from", item.From),
			zap.Float64("usdtValue", item.UsdtValue),
			zap.Float64("usdtValueExtra", item.UsdtValueExtra),
			zap.Float64("extraPercent", item.ExtraPercent),
			zap.Float64("percentGM", item.Percent),
			zap.Float64("gMReceive", item.GMReceive),
		)
	}

	u.Logger.AtLog().Logger.Info("ReAllocateGM - Review calculate totalGMReceive",
		zap.Float64("totalGMReceive", totalGMReceive),
		zap.Int("contributors", len(resp)),
		zap.Float64("usdtValue", usdtValue),
		zap.Float64("usdtExtra", usdtExtra))

	result.UsdtValue = usdtValue
	result.UsdtExtra = usdtExtra
	result.TotalGMReceive = totalGMReceive
	result.Items = resp

	err = u.Cache.SetDataWithExpireTime(keyReAllocate, result, 60*60*24*3) // 3 days
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM  set cache", zap.Error(err))
	}

	//backup to DB
	u.SaveReAllocateToDB(result)

	return result, nil
}

func (u Usecase) SaveReAllocateToDB(result *structure.AnalyticsProjectDeposit) {
	dbBackupItem := &entity.CachedGMReAllocatedDashBoard{
		Contributors:   len(result.Items),
		UsdtValue:      result.UsdtValue,
		UsdtExtra:      result.UsdtExtra,
		TotalGMReceive: result.TotalGMReceive,
	}

	dbBackupItem.SetID()
	dbBackupItem.SetCreatedAt()
	objID, err := u.Repo.Create(context.TODO(), dbBackupItem.TableName(), dbBackupItem)
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM backup data to DB Err", zap.Error(err))
	}

	//upload items to GCS
	bytesData, err := json.Marshal(result)
	if err == nil {
		fileName := fmt.Sprintf("items-%s.json", objID.Hex())
		base64Data := helpers.Base64Encode(bytesData)
		uploaded, err := u.GCS.UploadBaseToBucket(base64Data, fmt.Sprintf("backup/dashboard/gm/%s", fileName))
		if err == nil {
			dbBackupItem.BackupURL = fmt.Sprintf("%s%s", os.Getenv("GCS_ENDPOINT"), uploaded.Path)
			dbBackupItem.BackupFilePath = uploaded.Path
			dbBackupItem.BackupFileName = uploaded.Name

			_, err = u.Repo.UpdateOne(dbBackupItem.TableName(), bson.D{{utils.KEY_UUID, objID.Hex()}}, dbBackupItem)
			if err != nil {
				logger.AtLog.Logger.Error("ReAllocateGM update data for Backup", zap.Any("objID", objID), zap.Error(err))
			} else {
				//send message to slack
				end := time.Now().UTC()
				preText := fmt.Sprintf("[ReAllocated][Created backup] - Backup for reallocated has been created")
				title := fmt.Sprintf("UUID %s, URL: %s", dbBackupItem.UUID, dbBackupItem.BackupURL)
				content := fmt.Sprintf("Backup was created at: %v with USDT: %f, USDTExtra: %f, contributors: %d, totalGMReceive: %f - bk UUID: %s", end, dbBackupItem.UsdtValue, dbBackupItem.UsdtExtra, len(result.Items), result.TotalGMReceive, title)

				go u.SendGMMEssageToSlack(preText, content)
			}
		} else {
			logger.AtLog.Logger.Error("ReAllocateGM upload backup file to GCS", zap.Any("objID", objID), zap.Error(err))
		}
	} else {
		logger.AtLog.Logger.Error("ReAllocateGM create base64 data", zap.Any("objID", objID), zap.Error(err))
	}
	logger.AtLog.Logger.Info("ReAllocateGM backup data to DB", zap.Any("objID", objID))
}

func (u Usecase) GetExtraPercent(address string) float64 {
	address = strings.ToLower(address)

	if os.Getenv("GetReallocateData") == "true" {
		dataReAllocate, err := u.GetReallocateData()
		if err == nil && dataReAllocate != nil {
			for _, item := range dataReAllocate.Items {
				if strings.ToLower(item.From) == address {
					return item.ExtraPercent
				}
			}
			return float64(0)
		}
		return float64(0)
	}

	user, err := u.Repo.FindUserByWalletAddressEQ(address)
	if err == nil && user.UUID != "" {
		return 30.0
	}

	for key, value := range newList30 {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 30.0
		}
	}

	for key, value := range kll {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 25.0
		}
	}

	for key, value := range kll2 {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 25.0
		}
	}

	//TODO - move this nod into the other task
	tcBalance, err := u.TcClientPublicNode.GetBalance(context.TODO(), address)
	if err == nil && tcBalance.Cmp(big.NewInt(0)) > 0 {
		return 20.0
	}

	for key, value := range manualAddMore {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 20.0
		}
	}

	for key, value := range newList20 {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 20.0
		}
	}

	for key, value := range allowList {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 10.0
		}
	}

	for key, value := range allowList {
		if strings.ToLower(key) == strings.ToLower(address) && value == true {
			return 10.0
		}
	}

	if strings.ToLower("0x8d73888d01DD73d63b1a93DD447Dd988ec9B5158") == strings.ToLower(address) {
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

/*func (u Usecase) ClearCacheTop10GMDashboard() {
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
}*/

func (u Usecase) BackupGMDashboardCachedData(oldObject, newObject *structure.AnalyticsProjectDeposit) {
	if os.Getenv("ENV") != "mainnet" {
		return
	}

	inserted, err := u.UploadCachedAndCreateBk(oldObject, newObject)
	if err != nil {
		logger.AtLog.Logger.Info("BackupGMDashboardCachedData UploadCachedAndCreateBk", zap.Error(err))
	}
	logger.AtLog.Logger.Info("BackupGMDashboardCachedData", zap.Any("inserted", inserted))

	end := time.Now().UTC()
	preText := fmt.Sprintf("[Analytics][Created backup] - Get chart data for GM Dashboard")
	title := fmt.Sprintf("UUID %s, URL: %s", inserted.UUID, inserted.BackupURL)
	content := fmt.Sprintf("Backup was created at: %v with USDT: %f, contributors: %d - bk UUID: %s", end, newObject.UsdtValue, len(newObject.Items), title)

	go u.SendGMMEssageToSlack(preText, content)
	u.Logger.Info("Complete JobGetChartDataForGMCollection", zap.Any("data", newObject.UsdtValue))

	return
}

func (u Usecase) RestoreGMDashboardCachedData(uuid string) (*structure.AnalyticsProjectDeposit, error) {
	resp := &structure.AnalyticsProjectDeposit{}
	key := fmt.Sprintf(keyNotReAllocate)
	var cached *entity.CachedGMDashBoardNew
	var err error
	if uuid == "" {
		cached, err = u.Repo.GetTheLatestGMDashboardNewCached()
		if err != nil {
			return nil, err
		}
	} else {
		cached, err = u.Repo.FindGMDashboardNewCached(uuid)
		if err != nil {
			return nil, err
		}
	}

	data, err := u.GCS.ReadFile(cached.BackupFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	err = u.Cache.SetDataWithExpireTime(key, resp, 60*60*1)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func (u Usecase) UploadCachedAndCreateBk(oldData, newData *structure.AnalyticsProjectDeposit) (*entity.CachedGMDashBoardNew, error) {
	type uploadChan struct {
		Data *googlecloud.GcsUploadedObject
		Err  error
	}

	uploadFunc := func(id string, data *structure.AnalyticsProjectDeposit, uploadOldChan chan uploadChan) {
		uploaded := &googlecloud.GcsUploadedObject{}
		var err error

		defer func() {
			uploadOldChan <- uploadChan{
				Data: uploaded,
				Err:  err,
			}
		}()

		uploaded, err = u.UploadCachedDataToGCS(data, id)
		if err != nil {
			logger.AtLog.Logger.Error("UploadCachedAndCreateBk UploadCachedDataToGCS", zap.Error(err))
			return
		}

	}
	dbBackupItem := &entity.CachedGMDashBoardNew{}
	err := u.Repo.InsertGMDashboardNewCached(dbBackupItem)
	if err != nil {
		logger.AtLog.Logger.Error("UploadCachedAndCreateBk InsertGMDashboardNewCached", zap.Error(err))
		return nil, err
	}

	id := dbBackupItem.ID.Hex()
	uploadOldChan := make(chan uploadChan)
	uploadNewChan := make(chan uploadChan)

	go uploadFunc(id, oldData, uploadOldChan)
	go uploadFunc(id, newData, uploadNewChan)

	oldDataFromChan := <-uploadOldChan
	newDataFromChan := <-uploadNewChan

	if oldDataFromChan.Err != nil {
		logger.AtLog.Logger.Error("UploadCachedAndCreateBk oldDataFromChan", zap.Error(oldDataFromChan.Err))
		return nil, err
	}

	if newDataFromChan.Err != nil {
		logger.AtLog.Logger.Error("UploadCachedAndCreateBk newDataFromChan", zap.Error(newDataFromChan.Err))
		return nil, err
	}

	dbBackupItem.UsdtValue = newData.UsdtValue
	dbBackupItem.Contributors = len(newData.Items)

	dbBackupItem.OldBackupURL = fmt.Sprintf("%s%s", os.Getenv("GCS_ENDPOINT"), oldDataFromChan.Data.Path)
	dbBackupItem.OldBackupFilePath = oldDataFromChan.Data.Path
	dbBackupItem.OldBackupFileName = oldDataFromChan.Data.Name

	dbBackupItem.BackupURL = fmt.Sprintf("%s%s", os.Getenv("GCS_ENDPOINT"), newDataFromChan.Data.Path)
	dbBackupItem.BackupFilePath = newDataFromChan.Data.Path
	dbBackupItem.BackupFileName = newDataFromChan.Data.Name

	_, err = u.Repo.UpdateOne(dbBackupItem.TableName(), bson.D{{utils.KEY_UUID, id}}, dbBackupItem)
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM update data for Backup", zap.Any("objID", id), zap.Error(err))
		return nil, err
	}

	return dbBackupItem, nil
}

func (u Usecase) UploadCachedDataToGCS(data *structure.AnalyticsProjectDeposit, fileName string) (*googlecloud.GcsUploadedObject, error) {
	uploadedTime := time.Now().UTC().UnixNano()

	//upload items to GCS
	bytesData, err := json.Marshal(data)
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM create base64 data", zap.Int64("objID", uploadedTime), zap.Error(err))
		return nil, err
	}

	uploadedFileName := fmt.Sprintf("%s-items-%d.json", fileName, uploadedTime)
	base64Data := helpers.Base64Encode(bytesData)
	uploaded, err := u.GCS.UploadBaseToBucket(base64Data, fmt.Sprintf("backup/dashboard/gm-not-allocated/%s", uploadedFileName))
	if err != nil {
		logger.AtLog.Logger.Error("ReAllocateGM upload backup file to GCS", zap.Any("objID", uploadedTime), zap.Error(err))
		return nil, err
	}

	return uploaded, nil
}

/*func (u Usecase) ChartForGMDashboard() (*structure.GMDashBoardPercent, error) {
	past := time.Now().UTC().Add(time.Hour * -24)

	pastData := entity.AggregatedGMDashBoard{}
	pastDataArr, err := u.Repo.AggregateGMDashboardCachedDataByTime(&past)
	if err != nil {
		return nil, err
	}

	if len(pastDataArr) == 1 {
		pastData = pastDataArr[0]
	}

	if pastData.Usdt == 0 {
		pastData.Usdt = float64(1700000.000)
	}
	if pastData.Contributors == 0 {
		pastData.Contributors = int64(2100)
	}

	data, err := u.GetChartDataForGMCollection(true)
	if err != nil {
		return nil, err
	}

	contributors := float64(len(data.Items))
	usdt := data.UsdtValue

	percentContributor := ((contributors - float64(pastData.Contributors)) / contributors) * 100
	percentUsdt := ((usdt - pastData.Usdt) / usdt) * 100

	resp := &structure.GMDashBoardPercent{
		PastContributors:   pastData.Contributors,
		PastUSDT:           pastData.Usdt,
		USDT:               usdt,
		Contributor:        int64(len(data.Items)),
		PercentUSDT:        percentUsdt,
		PercentContributor: percentContributor,
	}

	return resp, nil
}*/

// DUYNQ get old
func (u Usecase) GetDataOld() (*structure.AnalyticsProjectDeposit, error) {
	result := &structure.AnalyticsProjectDeposit{}
	if true {
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
				temp, err := u.GetChartDataEthForGMCollection(entity.NewCityGm{UserAddress: strings.ToLower(wallet), Address: strings.ToLower(gmAddress)}, []string{}, true, ens, "")
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
				temp, err := u.GetChartDataBTCForGMCollection(entity.NewCityGm{UserAddress: strings.ToLower(wallet), Address: strings.ToLower(wallet)}, []string{}, true)
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
		//erc20DataFromChan := <-erc20DataChan

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

		return result, nil
	}

	return result, nil
}

// DUYNQ backup api
func (u Usecase) GetChartDataForGMCollectionBackup() (*structure.AnalyticsProjectDeposit, error) {
	// try get data from reAllocate, check config
	if os.Getenv("GetReallocateData") == "true" {
		dataReAllocate, err := u.GetReallocateData()
		if err == nil && dataReAllocate != nil {
			return dataReAllocate, nil
		}
		return nil, err
	}

	fullUrl := "https://www.fprotocol.io/api/gm/deposit"
	statusCode, req, err := request.GetRequest(fullUrl)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataForGMCollectionBackup", zap.Error(err), zap.Int("statusCode", statusCode))
		return nil, err
	}

	if statusCode != 200 {
		err := errors.New(fmt.Sprintf("Response with status: %d", statusCode))
		logger.AtLog.Logger.Error("GetChartDataForGMCollectionBackup", zap.Error(err), zap.Int("statusCode", statusCode))
		return nil, err
	}

	rsp := &structure.AnalyticsProjectDepositExternal{}
	err = json.Unmarshal(req, rsp)
	if err != nil {
		logger.AtLog.Logger.Error("GetChartDataForGMCollectionBackup", zap.Error(err), zap.Int("statusCode", statusCode))
		return nil, err
	}

	err = u.Cache.SetDataWithExpireTime(keyNotReAllocate, rsp.Data, 60*60*24*3)
	if err != nil {
		u.Logger.AtLog().Logger.Error("GetChartDataForGMCollection: SetDataWithExpireTime ...", zap.Error(err))
	}

	logger.AtLog.Logger.Info("GetChartDataForGMCollectionBackup", zap.Float64("usdt", rsp.Data.UsdtValue), zap.Int("items", len(rsp.Data.Items)))

	return &rsp.Data, nil
}
