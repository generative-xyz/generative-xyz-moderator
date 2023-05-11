package usecase

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
