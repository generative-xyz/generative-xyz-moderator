package web3

import (
	"math/big"

	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

// getCustomerInfoFromBalance returns customer info from balance
func getCustomerInfoFromBalance(balance *big.Int, address string) (*nftdata.NFTCustomerInfo, error) {
	customerInfo := nftdata.NewNFTCustomerInfo(nftdata.DefaultItem, address).WithBalance(balance.Int64())
	return customerInfo, nil
}
