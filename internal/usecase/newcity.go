package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiCreateNewGM(addressInput string) (interface{}, error) {

	if !eth.ValidateAddress(addressInput) {
		return nil, errors.New("you address invalid")
	}

	// get temp address from db:
	itemEth, err := u.Repo.FindNewCityGmByUserAddress(addressInput, utils.NETWORK_ETH)

	if err != nil {
		if !strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, err
		}
	}

	fmt.Println("itemEth: ", itemEth)

	if itemEth == nil {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err := ethClient.GenerateAddress()
		if err != nil {
			logger.AtLog.Logger.Error("ApiCreateNewGM.ethClient.GenerateAddress", zap.Error(err))
			return nil, err
		}
		privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
			return nil, err
		}
		itemEth = &entity.NewCityGm{
			UserAddress: addressInput,
			Type:        utils.NETWORK_ETH,
			Status:      1,
			Address:     receiveAddress, // temp address for the user send to
			PrivateKey:  privateKeyEnCrypt,
		}

		err = u.Repo.InsertNewCityGm(itemEth)

		if err != nil {
			return nil, err
		}
	}

	itemBtc, err := u.Repo.FindNewCityGmByUserAddress(addressInput, utils.NETWORK_BTC)
	if err != nil {
		if !strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, err
		}
	}

	fmt.Println("itemBtc: ", itemBtc)

	if itemBtc == nil {

		privateKeyBtc, _, receiveAddressBtc, err := btc.GenerateAddressSegwit()
		if err != nil {
			logger.AtLog.Logger.Error("u.ApiCreateNewGM.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}
		privateKeyEnCryptBtc, err := encrypt.EncryptToString(privateKeyBtc, os.Getenv("SECRET_KEY"))
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
			return nil, err
		}
		itemBtc = &entity.NewCityGm{
			UserAddress: addressInput,
			Type:        utils.NETWORK_BTC,
			Status:      1,
			Address:     receiveAddressBtc, // temp address for the user send to
			PrivateKey:  privateKeyEnCryptBtc,
		}

		err = u.Repo.InsertNewCityGm(itemBtc)

		if err != nil {
			return "", err
		}
	}

	return map[string]string{
		"eth": itemEth.Address,
		"btc": itemBtc.Address,
	}, nil
}

// admin
func (u Usecase) ApiAdminCrawlFunds() (interface{}, error) {

	var returnData []*entity.NewCityGm

	list, _ := u.Repo.ListNewCityGmByStatus([]int{1}) // 1 pending, 2: tx, 3 confirm

	ethWithdrawAddrses := os.Getenv("GM_ETH_WITHDRAW_ADDRESS")

	if len(ethWithdrawAddrses) == 0 {
		return nil, errors.New("GM_ETH_WITHDRAW_ADDRESS not found")
	}

	if len(list) > 0 {
		for _, item := range list {
			// hardcode for test withdraw funds:

			if item.Type == utils.NETWORK_ETH {

				// check balance:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				balance, _ := u.EthClient.GetBalance(ctx, item.Address)

				if balance != nil && balance.Cmp(big.NewInt(0)) > 0 {

					// update balance
					item.NativeAmount = append(item.NativeAmount, balance.String())
					_, err := u.Repo.UpdateNewCityGm(item)
					if err != nil {
						return nil, err
					}

					if item.UserAddress == ethWithdrawAddrses {
						// send all:
						privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
						if err != nil {
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.Decrypt.%s.Error", item.Address), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.DecryptToString", err.Error(), true)
							continue
						}
						tx, value, err := u.EthClient.TransferMax(privateKeyDeCrypt, ethWithdrawAddrses)
						if err != nil {
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.ethClient.TransferMax.%s.Error", item.Address), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.ethClient.TransferMax", err.Error(), true)
							continue
						}
						_ = value

						item.TxNatives = append(item.NativeAmount, tx)
						_, err = u.Repo.UpdateNewCityGm(item)
						if err != nil {
							return nil, err
						}

						returnData = append(returnData, item)
					}

				}

			} else if item.Type == utils.NETWORK_BTC {
				// todo
			}

		}
	}
	return returnData, nil
}
