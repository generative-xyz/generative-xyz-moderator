package usecase

import (
	"errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/logger"
	"strings"
)

func (u Usecase) CreateOrderReceiveAddress(input structure.OrderBtcData) (*entity.OrdersAddress, error) {
	walletAddress := &entity.OrdersAddress{}
	receiveAddress := ""
	privateKey := ""
	var err error

	if input.OrderID == "" {
		err = errors.New("order_id is required")
		return nil, err
	}

	// verify paytype:
	if input.PayType != utils.NETWORK_BTC && input.PayType != utils.NETWORK_ETH {
		err = errors.New("only support payType is eth or btc")
		return nil, err
	}

	ord, _ := u.Repo.FindOrderBy(input.OrderID, input.PayType)
	if ord != nil {
		return ord, nil
	}

	// check type:
	if input.PayType == utils.NETWORK_BTC {
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}

	} else if input.PayType == utils.NETWORK_ETH {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			logger.AtLog.Logger.Error("CreateMintReceiveAddress.ethClient.GenerateAddress", zap.Error(err))
			return nil, err
		}
	}

	if len(receiveAddress) == 0 || len(privateKey) == 0 {
		err = errors.New("can not create the wallet")
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.GenerateAddress", zap.Error(err))
		return nil, err
	}

	walletAddress.AddressType = input.PayType
	walletAddress.Address = strings.ToLower(receiveAddress)
	walletAddress.PrivateKey = privateKey
	walletAddress.OrderID = input.OrderID

	// insert now:
	err = u.Repo.InsertOrder(walletAddress)
	if err != nil {
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.InsertMintNftBtc", zap.Error(err))
		return nil, err
	}

	return walletAddress, nil
}
