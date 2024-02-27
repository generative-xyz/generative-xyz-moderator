package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"math/big"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"strings"
	"sync"
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

	orderIDs := []string{}
	amount := make(map[string]string)
	for _, item := range d.Orders {
		orderIDs = append(orderIDs, item.Id)
		amount[item.Id] = item.Amount
	}

	orders, err := u.Repo.FindOrderByIDs(orderIDs)
	if err != nil {
		return nil, err
	}

	orderDetails := make(map[string]entity.OrdersAddress)
	for _, i := range orders {
		orderDetails[i.OrderID] = *i
	}

	for _, item := range d.Orders {
		d1, ok := orderDetails[item.Id] // from DB
		if ok {
			item.PayType = d1.AddressType
			item.Status = int(d1.Status)
			item.PaymentAddress = d1.Address
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

func (u Usecase) GetOrderByIDFromApi(id string) (*structure.ApiOrderItemResp, error) {
	cachedKey := fmt.Sprintf("order.%s", id)
	respData := &structure.ApiOrderItemResp{}

	err := u.Cache.GetObjectData(cachedKey, respData)
	if err != nil || respData == nil {
		grailAPI := os.Getenv("GRAIL_API")
		if grailAPI == "" {
			grailAPI = "https://generative.xyz/api/v1"
		}

		_url := fmt.Sprintf("%s/order/by-id/%s", grailAPI, id)
		_b, _, _, err := helpers.HttpRequest(_url, "GET", map[string]string{}, nil)
		if err != nil {
			return nil, err
		}

		resp := &structure.ApiOrderDetailResp{}
		err = json.Unmarshal(_b, resp)
		if err != nil {
			return nil, err
		}

		if resp.Message != nil {
			err = errors.New(resp.Message.Message)
			return nil, err
		}

		respData = &resp.Data.Order
		u.Cache.SetDataWithExpireTime(cachedKey, respData, 86400*10) // 10 days
		return respData, nil
	}

	return respData, nil
}

func (u Usecase) ProcessOrderPaymentApi(orderID, txHash, walletAddress string) (interface{}, error) {
	//TODO - implement me
	//grailAPI := os.Getenv("GRAIL_API")
	//if grailAPI == "" {
	//	grailAPI = "https://generative.xyz/api/v1"
	//}
	//
	//postData := make(map[string]interface{})
	//postData["order_id"] = orderID
	//postData["tx_hash"] = txHash
	//postData["wallet_address"] = walletAddress
	//
	//_url := fmt.Sprintf("%s/order/complete", grailAPI)
	//_b, _, _, err := helpers.HttpRequest(_url, "POST", map[string]string{}, postData)
	//if err != nil {
	//	return nil, err
	//}
	//
	//resp := &structure.ApiOrderDetailResp{}
	//err = json.Unmarshal(_b, resp)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if resp.Message != nil {
	//	err = errors.New(resp.Message.Message)
	//	return nil, err
	//}
	//
	//respData = &resp.Data.Order
	//return respData, nil

	return nil, nil
}

func (u Usecase) ProcessCancelOrderApi(orderID string) (interface{}, error) {
	//TODO - implement me
	//grailAPI := os.Getenv("GRAIL_API")
	//if grailAPI == "" {
	//	grailAPI = "https://generative.xyz/api/v1"
	//}
	//
	//postData := make(map[string]interface{})
	//postData["order_id"] = orderID
	//postData["tx_hash"] = txHash
	//postData["wallet_address"] = walletAddress
	//
	//_url := fmt.Sprintf("%s/order/complete", grailAPI)
	//_b, _, _, err := helpers.HttpRequest(_url, "POST", map[string]string{}, postData)
	//if err != nil {
	//	return nil, err
	//}
	//
	//resp := &structure.ApiOrderDetailResp{}
	//err = json.Unmarshal(_b, resp)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if resp.Message != nil {
	//	err = errors.New(resp.Message.Message)
	//	return nil, err
	//}
	//
	//respData = &resp.Data.Order
	//return respData, nil

	return nil, nil
}

func (u Usecase) CheckOrderStatus(wg *sync.WaitGroup, input chan entity.OrdersAddress, output chan structure.OrderStatusChan) {
	order := <-input
	var err error
	statusP := new(int)
	status := 0
	statusP = &status

	defer wg.Done()

	defer func() {

		isUpdateStatus := false
		if *statusP != int(order.Status) {
			isUpdateStatus = true
		}

		output <- structure.OrderStatusChan{
			OrderID:        order.OrderID,
			Err:            err,
			Status:         *statusP,
			PayType:        string(order.AddressType),
			OrderAmount:    order.Amount,
			PaymentAddress: order.Address,
			IsUpdated:      isUpdateStatus,
		}

	}()

	if strings.EqualFold(order.AddressType, string(entity.ETH)) {
		//address := "0x13BB7Bf390B55A7d5bF44c4dcEcdFEB1Dd2a884a"
		address := order.Address
		balance, err := u.EthClient.GetBalance(context.TODO(), address)
		if err != nil {
			return
		}

		aF, _ := big.NewFloat(0).SetString(order.Amount)
		aF = aF.Mul(aF, big.NewFloat(1e18))
		balanceF := big.NewFloat(0).SetInt(balance)

		if balanceF.Cmp(aF) >= 0 { //balance >= amount
			status = int(entity.Order_Paid)
		}

	}

	if strings.EqualFold(order.AddressType, string(entity.BTC)) {
		_, bs, err := u.buildBTCClient()
		if err != nil {
			return
		}
		//address := "bc1pv47nhns0xeljuzkdtvdk3qxm5zk42cmmwgrz3tt4x4he6kvwcjzqm2ml2h"
		address := order.Address
		balance, _, err := bs.GetBalance(address)
		if err != nil {
			// get balance from quicknode:
			var balanceQuickNode *structure.BlockCypherWalletInfo
			balanceQuickNode, err = btc.GetBalanceFromQuickNode(order.Address, u.Config.QuicknodeAPI)
			if err == nil {
				if balanceQuickNode != nil {
					balance = big.NewInt(int64(balanceQuickNode.Balance))
				}
			}
		}

		aF, _ := big.NewFloat(0).SetString(order.Amount)
		aF = aF.Mul(aF, big.NewFloat(1e8))
		balanceF := big.NewFloat(0).SetInt(balance)

		if balanceF.Cmp(aF) >= 0 { //balance >= amount
			status = int(entity.Order_Paid)
		}
	}

}

func (u Usecase) JobSyncOrderPaymentStatus() error {
	orders, err := u.Repo.FindOrderByStatus([]entity.OrderStatus{entity.Order_Pending})
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("JobSyncOrderPaymentStatus", zap.Int("orders", len(orders)))

	inchan := make(chan entity.OrdersAddress, len(orders))
	outChan := make(chan structure.OrderStatusChan, len(orders))
	wg := sync.WaitGroup{}
	wg.Add(len(orders))

	for range orders {
		go u.CheckOrderStatus(&wg, inchan, outChan)
	}

	for _, i := range orders {
		od, err := u.GetOrderByIDFromApi(i.OrderID)
		if err == nil {
			i.Amount = od.Amount
		} else {
			i.Amount = "99999" // cannot update status with this amount
		}

		inchan <- *i
	}

	for range orders {
		data := <-outChan
		if data.Err != nil {
			continue
		}

		if data.IsUpdated {
			_, err := u.Repo.UpdateOrderStatus(data.OrderID, entity.OrderStatus(data.Status))
			if err != nil {
				continue
			}

			//call API
			switch entity.OrderStatus(data.Status) {
			case entity.Order_Paid:
				u.ProcessOrderPaymentApi(data.OrderID, fmt.Sprintf("https://etherscan.io/address/%s", data.PaymentAddress), data.PaymentAddress)
			case entity.Order_Cancel:
				u.ProcessCancelOrderApi(data.OrderID)
			default:
				spew.Dump("handle default")
			}
		}
	}

	wg.Wait()
	return nil
}
