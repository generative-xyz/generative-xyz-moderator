package binance

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type BinanceTestSuite struct {
	suite.Suite
}

func (t *BinanceTestSuite) SetupTest() {

}

func TestBinanceSwap(t *testing.T) {

	secret := os.Getenv("BNB_SECRET")
	apiKey := os.Getenv("BNB_API_KEY")

	bs := NewBinanceService(apiKey, secret)

	pair := "ETHBTC"

	order, err := bs.SwapEth2Btc("0.007", pair)

	if err != nil {
		fmt.Println("bs.SwapEth2Btc => err: ", err.Error())
	} else {
		fmt.Println("bs.order => order: ", order.OrderID)
		fmt.Println("bs.Status => order: ", order.Status)
	}

	assert.Equal(t, false, true)

}

func TestBinanceService(t *testing.T) {

	secret := os.Getenv("BNB_SECRET")
	apiKey := os.Getenv("BNB_API_KEY")

	bs := NewBinanceService(apiKey, secret)

	pair := "ETHBTC"
	var orderID int64 = 1

	status, err := bs.GetOrderStatus(orderID, pair)

	fmt.Println("bs.GetOrderStatus => status: ", status)

	if err != nil {
		fmt.Println("bs.GetOrderStatus => status, err: ", err.Error())
	}

	assert.Equal(t, false, true)

}
