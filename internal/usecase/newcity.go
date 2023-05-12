package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"rederinghub.io/external/opensea"

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

	secretKeyName := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")
	if len(secretKeyName) == 0 {
		return nil, errors.New("please config google key first!")
	}

	keyToEncrypt, err := GetGoogleSecretKey(secretKeyName)
	if err != nil {
		return nil, errors.New("can't not get secretKey from key name")
	}

	if len(keyToEncrypt) == 0 {
		return nil, errors.New("please config key first!")
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
		privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, keyToEncrypt)
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
			KeyVersion:  1, // 2 from now
		}

		err = u.Repo.InsertNewCityGm(itemEth)

		if err != nil {
			return nil, err
		}

		go func(item *entity.NewCityGm) {
			ens, errENS := u.EthClient.GetEns(addressInput)
			if errENS == nil {
				if len(ens) > 0 {
					itemEth.ENS = ens
				}
			}

			avatar, errAvatar := opensea.OpenseaService{}.GetProfileAvatar(addressInput)
			if errAvatar == nil {
				if len(avatar) > 0 {
					itemEth.Avatar = avatar
				}
			}
			u.Repo.UpdateNewCityGmENSAvatar(item)
		}(itemEth)

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
		privateKeyEnCryptBtc, err := encrypt.EncryptToString(privateKeyBtc, keyToEncrypt)
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
			KeyVersion:  1, // 2 from now
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

	if false {
		if len(os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")) == 0 {
			return nil, errors.New("key to get key is empty")
		}
		secretKeyName := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")

		secretKey, err := GetGoogleSecretKey(secretKeyName)
		if err != nil {
			return nil, errors.New("can't not get secretKey from key name")
		}

		// try to encrypt
		privKeyTest := "hello 123"

		privateKeyDeCrypt, err := encrypt.EncryptToString(privKeyTest, secretKey)
		if err != nil {
			return nil, errors.New("can't not EncryptToString" + err.Error())

		}
		// try to decrypt
		valueDecrypt, err := encrypt.DecryptToString(privateKeyDeCrypt, secretKey)
		if err != nil {
			return nil, errors.New("can't not DecryptToString" + err.Error())

		}
		ok := strings.EqualFold(valueDecrypt, privKeyTest)

		if !ok {
			return nil, errors.New("can't not compare")
		}

		return privateKeyDeCrypt, nil
	}

	var returnData []*entity.NewCityGm

	list, err := u.Repo.ListNewCityGmByStatus([]int{1}) // 1 pending, 2: tx, 3 confirm

	if err != nil {
		return nil, err
	}

	ethWithdrawAddrses := os.Getenv("GM_ETH_WITHDRAW_ADDRESS")

	if len(ethWithdrawAddrses) == 0 {
		return nil, errors.New("GM_ETH_WITHDRAW_ADDRESS not found")
	}

	btcWithdrawAddrses := os.Getenv("GM_BTC_WITHDRAW_ADDRESS")

	if len(btcWithdrawAddrses) == 0 {
		return nil, errors.New("GM_BTC_WITHDRAW_ADDRESS not found")
	}

	keyToDecodeV1 := os.Getenv("SECRET_KEY")

	keyToDecodeGoogle := os.Getenv("GENERATIVE_ENCRYPT_SECRET_KEY_NAME")
	keyToDecodeV2, err := GetGoogleSecretKey(keyToDecodeGoogle)
	if err != nil {
		return nil, errors.New("can't not get secretKey from key name")
	}

	if len(list) > 0 {
		for _, item := range list {

			if item.Type == utils.NETWORK_ETH {

				// check balance:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				balance, _ := u.EthClient.GetBalance(ctx, item.Address)

				if balance != nil && balance.Cmp(big.NewInt(0)) > 0 {

					fmt.Println("balance: ", balance)

					// update balance
					item.NativeAmount = append(item.NativeAmount, balance.String())
					_, err := u.Repo.UpdateNewCityGm(item)
					if err != nil {
						return nil, err
					}

					// hardcode for test withdraw funds:
					// if item.UserAddress == ethWithdrawAddrses {
					// send all:
					keyToDecode := keyToDecodeV1

					if item.KeyVersion == 1 {
						keyToDecode = keyToDecodeV2
					}
					if len(keyToDecode) == 0 {
						return nil, errors.New("key to decrypt is empty")
					}

					privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, keyToDecode)
					if err != nil {
						logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.Decrypt.%s.Error", item.Address), zap.Error(err))
						go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.DecryptToString", err.Error(), true)
						time.Sleep(300 * time.Millisecond)
						continue
					}
					tx, value, err := u.EthClient.TransferMax(privateKeyDeCrypt, ethWithdrawAddrses)
					if err != nil {

						// check if not enough balance:
						if strings.Contains(err.Error(), "rlp: cannot encode negative big.Int") {
							item.Status = -1
							u.Repo.UpdateNewCityGm(item)
						}

						logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.ethClient.TransferMax.%s.Error", item.Address), zap.Error(err))
						go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.ethClient.TransferMax", err.Error(), true)
						time.Sleep(300 * time.Millisecond)
						continue
					}
					_ = value

					item.TxNatives = append(item.TxNatives, tx)
					item.Status = 2 // tx pending
					_, err = u.Repo.UpdateNewCityGm(item)
					if err != nil {
						return nil, err
					}

					returnData = append(returnData, item)
					// }

				}
				time.Sleep(300 * time.Millisecond)

			} else if item.Type == utils.NETWORK_BTC {
				// todo

				balanceQuickNode, _ := btc.GetBalanceFromQuickNode(item.Address, u.Config.QuicknodeAPI)
				if balanceQuickNode != nil {
					balance := balanceQuickNode.Balance
					if balance > 0 {

						keyToDecode := keyToDecodeV1

						if item.KeyVersion == 1 {
							keyToDecode = keyToDecodeV2
						}
						if len(keyToDecode) == 0 {
							return nil, errors.New("key to decrypt is empty")
						}

						// send max:
						// send master now:
						privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, keyToDecode)
						if err != nil {
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.Decrypt.%s.Error", item.Address), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.DecryptToString", err.Error(), true)
							continue
						}

						tx, err := u.BsClient.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.Address, btcWithdrawAddrses, -1, btc.PreferenceMedium)
						if err != nil {

							// check if not enough balance:
							if strings.Contains(err.Error(), "insufficient priority and fee for relay") {
								item.Status = -1
								_, err = u.Repo.UpdateNewCityGm(item)
								if err != nil {
									return nil, err
								}

							}

							if strings.Contains(err.Error(), "already exists") {
								item.Status = -2
								_, err = u.Repo.UpdateNewCityGm(item)
								if err != nil {
									return nil, err
								}

							}
							logger.AtLog.Logger.Error(fmt.Sprintf("ApiAdminCrawlFunds.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", btcWithdrawAddrses), zap.Error(err))
							go u.trackMintNftBtcHistory(item.UUID, "ApiAdminCrawlFunds", item.TableName(), item.Status, "ApiAdminCrawlFunds.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
							time.Sleep(300 * time.Millisecond)
							continue
						}
						// save tx:
						item.TxNatives = append(item.TxNatives, tx)
						item.Status = 2 // tx pending
						_, err = u.Repo.UpdateNewCityGm(item)
						if err != nil {
							return nil, err
						}

						returnData = append(returnData, item)
					}
				}
				time.Sleep(300 * time.Millisecond)
			}

		}
	}
	return returnData, nil
}
