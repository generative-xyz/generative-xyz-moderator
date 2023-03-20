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
