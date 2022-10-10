package web3

import (
	"rederinghub.io/pkg/third-party/crypto/web3/covalenthq"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

type ethereumClientImpl struct {
	covalenthqClient *covalenthq.Client
}

func (e ethereumClientImpl) NFTFeeds(req *BalanceRequest) (interface{}, error) {
	return e.covalenthqClient.GetNFTFeeds(req.Address, req.ChainID)
}

func (e ethereumClientImpl) BalanceOf(req *BalanceRequest) (*nftdata.NFTCustomerInfo, error) {
	info, err := e.covalenthqClient.GetBalance(req.Address, req.ChainID, false /* noFetchNFTMetadata */)
	if err != nil {
		// retry and no fetch nft metadata in case of error
		info, err = e.covalenthqClient.GetBalance(req.Address, req.ChainID, true /* noFetchNFTMetadata */)
		if err != nil {
			return nil, err
		}
	}

	return info.GetNFTCustomerInfo(), nil
}

func NewClientEthereum(baseURL, apiKey string) Strategy {
	return &ethereumClientImpl{
		covalenthqClient: covalenthq.NewClient(baseURL, apiKey),
	}
}
