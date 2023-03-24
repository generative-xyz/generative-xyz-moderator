package binance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2"
)

func NewBinanceService(apiKey, secretKey string) *BinanceService {
	return &BinanceService{
		binanceClient: binance.NewClient(apiKey, secretKey),
	}

}

type BinanceService struct {
	binanceClient *binance.Client
}

// pair: "ETHBTC"

func (bs *BinanceService) SwapEth2Btc(ethAmount string, pair string) (*binance.CreateOrderResponse, error) {

	// ethAmount: the amount of ETH you want to swap

	// Get the current ETH/BTC exchange rate:
	// ticker, err := bs.binanceClient.NewListPriceChangeStatsService().
	// 	Symbol(pair).
	// 	Do(context.Background())
	// if err != nil {
	// 	return orderStatus, err
	// }
	// lastPrice, err := strconv.ParseFloat(ticker[0].LastPrice, 64)
	// if err != nil {
	// 	return orderStatus, err
	// }
	// ethPrice := lastPrice

	// // Calculate the amount of BTC you will receive based on the amount of ETH you want to swap:
	// _ = ethAmount * ethPrice

	// Place the order to swap ETH to BTC:
	order, err := bs.binanceClient.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(ethAmount).
		Do(context.Background())
	if err != nil {
		return order, err

	}

	fmt.Println("order status new: ", order.Status)

	return nil, nil

}

func (bs *BinanceService) GetOrderStatus(orderID int64, pair string) (string, error) {

	// Check the status of your order:
	order, err := bs.binanceClient.NewGetOrderService().
		Symbol(pair).
		OrderID(orderID).
		Do(context.Background())
	if err != nil {
		return "", err

	}
	fmt.Println("order status check: ", order.Status)
	return string(order.Status), nil
}

func (bs *BinanceService) WithdrawAsset(amount, withdrawAddress, coin string) (string, error) {
	withdrawOrder, err := bs.binanceClient.NewCreateWithdrawService().
		Coin(coin).
		Address(withdrawAddress).
		Amount(amount).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	withdrawID := withdrawOrder.ID
	fmt.Println("Withdraw ID:", withdrawID)

	return withdrawOrder.ID, nil
}

func (bs *BinanceService) GetWithdrawStatus(withdrawID string) (int, error) {

	withdrawHistory, err := bs.binanceClient.NewListWithdrawsService().Do(context.Background())
	if err != nil {
		return 0, err
	}

	for _, withdraw := range withdrawHistory {
		if strings.EqualFold(withdraw.ID, withdrawID) {
			fmt.Println("Withdraw Status check: ", withdraw.Status)
			return withdraw.Status, nil
		}
	}

	return 0, errors.New("order not found")
}
