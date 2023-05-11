package usecase

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ApiCreateNewGM(addressInput, typeReq string) (string, error) {

	if !eth.ValidateAddress(addressInput) {
		return "", errors.New("you address invalid")
	}

	if typeReq != utils.NETWORK_BTC && typeReq != utils.NETWORK_ETH {
		return "", errors.New("network type invalid")
	}

	receiveAddress := ""
	privateKey := ""
	var err error

	// get temp address from db:
	item, _ := u.Repo.FindNewCityGmByUserAddress(addressInput, typeReq)
	if item != nil {
		return item.Address, nil
	}

	if typeReq == utils.NETWORK_BTC {
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			logger.AtLog.Logger.Error("u.ApiCreateNewGM.GenerateAddressSegwit", zap.Error(err))
			return "", err
		}

	} else if typeReq == utils.NETWORK_ETH {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			logger.AtLog.Logger.Error("ApiCreateNewGM.ethClient.GenerateAddress", zap.Error(err))
			return "", err
		}
	}

	privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
		return "", err
	}

	item = &entity.NewCityGm{
		UserAddress: addressInput,
		Type:        typeReq,

		Address:    receiveAddress, // temp address for the user send to
		PrivateKey: privateKeyEnCrypt,
	}

	err = u.Repo.InsertNewCityGm(item)

	if err != nil {
		return "", err
	}

	return item.Address, nil
}
