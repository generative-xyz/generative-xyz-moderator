package contractaddress

import (
	"math"
	"math/big"

	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
)

func balanceToFloat64(balance *big.Int, coin string) float64 {
	return float64(balance.Uint64()) / math.Pow10(cryptocurrency.RoundPlacesByCurrency[coin])
}

// Float64ToBalance converts a float64 to a big.Int by coin
func Float64ToBalance(value float64, coin string) *big.Int {
	int64Value := int64(value * math.Pow10(cryptocurrency.RoundPlacesByCurrency[coin]))
	return big.NewInt(int64Value)
}
