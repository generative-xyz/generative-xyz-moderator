package web3

import "rederinghub.io/pkg/third-party/crypto/web3/nftdata"

type Strategy interface {
	BalanceOf(req *BalanceRequest) (*nftdata.NFTCustomerInfo, error)
	NFTFeeds(req *BalanceRequest) (interface{}, error)
}
