package binance

import (
	"context"
	"fmt"
	"strconv"

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

func (bs *BinanceService) SwapEth2Btc(ethAmount float64, pair string) (string, error) {

	// ethAmount: the amount of ETH you want to swap

	orderStatus := ""

	// Get the current ETH/BTC exchange rate:
	ticker, err := bs.binanceClient.NewListPriceChangeStatsService().
		Symbol(pair).
		Do(context.Background())
	if err != nil {
		return orderStatus, err
	}
	lastPrice, err := strconv.ParseFloat(ticker[0].LastPrice, 64)
	if err != nil {
		// handle error
	}
	ethPrice := lastPrice

	// Calculate the amount of BTC you will receive based on the amount of ETH you want to swap:
	btcAmount := ethAmount * ethPrice

	// Place the order to swap ETH to BTC:
	order, err := bs.binanceClient.NewCreateOrderService().
		Symbol(pair).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.8f", btcAmount)).
		Do(context.Background())
	if err != nil {
		return orderStatus, err

	}

	fmt.Println("order status new: ", order.Status)

	return string(order.Status), nil

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
