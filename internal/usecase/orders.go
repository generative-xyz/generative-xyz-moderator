package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
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

func (u Usecase) ListOrders(f structure.FilterOrders) (interface{}, error) {
	d, err := u.ListOrdersFromApi(f)
	if err != nil {
		return nil, err
	}

	//TODO - process here
	for i, item := range d.Orders {
		item.PayType = string(entity.ETH)

		item.Status = int(entity.Order_Pending)
		if i%2 == 0 {
			item.Status = int(entity.Order_Paid)
		}

	}

	return d, nil
}

func (u Usecase) ListOrdersFromApi(f structure.FilterOrders) (*structure.ApiOrderDataResp, error) {
	grailAPI := os.Getenv("GRAIL_API")
	if grailAPI == "" {
		grailAPI = "https://generative.xyz/api/v1"
	}

	if f.Email == nil {
		return nil, errors.New("email is required")
	}

	if *f.Email == "" {
		return nil, errors.New("email is not empty")
	}

	_url := fmt.Sprintf("%s/order/by-email/list?email=%s", grailAPI, *f.Email)
	_b, _, _, err := helpers.HttpRequest(_url, "GET", map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	resp := &structure.ApiOrderResp{}
	err = json.Unmarshal(_b, resp)
	if err != nil {
		return nil, err
	}

	if resp.Message != nil {
		err = errors.New(resp.Message.Message)
		return nil, err
	}

	_d := resp.Data
	return &_d, nil
}
