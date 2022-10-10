package web3

import (
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/third-party/crypto/contractaddress"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

type solanaClientImpl struct {
	contract *contractaddress.Client
}

func (s solanaClientImpl) NFTFeeds(req *BalanceRequest) (interface{}, error) {
	return nil, nil
}

func (s solanaClientImpl) BalanceOf(req *BalanceRequest) (*nftdata.NFTCustomerInfo, error) {
	balanceFloat64, err := s.contract.CheckBalance(cryptocurrency.Solana, req.Address, false)
	if err != nil {
		return nil, err
	}

	balance := contractaddress.Float64ToBalance(balanceFloat64, cryptocurrency.Solana)
	return getCustomerInfoFromBalance(balance, req.Address)
}

func NewClientSolana(isProduction bool) Strategy {
	solonaChain := contractaddress.NewSolanaChain(isProduction)
	solanaContract := contractaddress.NewClient()
	solanaContract.RegisterCoins(solonaChain, cryptocurrency.Solana)

	return &solanaClientImpl{
		contract: solanaContract,
	}
}
