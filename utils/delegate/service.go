package delegate

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"math/big"
)

type DelegationType int32

const (
	DelegationTypeUndefined DelegationType = iota
	DelegationTypeAll
	DelegationTypeContract
	DelegationTypeToken
)

// See supported chain here: https://docs.delegate.cash/delegatecash/contract-addresses
const delegateCashAddress = "0x00000000000076a84fef008cdabe6409d2fe638b"

type Service struct {
	contractAddress string

	delegateEthService *Delegate
}

func NewService(
	ethClient *ethclient.Client) (*Service, error) {

	hexAddress := common.HexToAddress(delegateCashAddress)
	delegateEthService, err := NewDelegate(hexAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &Service{
		contractAddress:    delegateCashAddress,
		delegateEthService: delegateEthService,
	}, nil
}

// GetDelegationsByDelegate Returns all active delegations a given delegate is able to claim on behalf of.
// delegate is the delegate that you would like to retrieve delegations for
func (s *Service) GetDelegationsByDelegate(ctx context.Context, delegateAddr string) (delegations []IDelegationRegistryDelegationInfo, err error) {
	hexAddress := common.HexToAddress(delegateAddr)
	delegations = make([]IDelegationRegistryDelegationInfo, 0)

	polygonDelegations, err := s.delegateEthService.GetDelegationsByDelegate(&bind.CallOpts{Context: ctx, Pending: false}, hexAddress)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	delegations = append(delegations, polygonDelegations...)

	return delegations, nil
}

// CheckDelegateForContract returns true if the address is delegated to act on your behalf for a token contract or an
// entire vault.
func (s *Service) CheckDelegateForContract(ctx context.Context, delegateAddr string, vaultAddr string, contractAddr string) (bool, error) {
	delegate := common.HexToAddress(delegateAddr)
	vault := common.HexToAddress(vaultAddr)
	contract := common.HexToAddress(contractAddr)
	ok, err := s.delegateEthService.CheckDelegateForContract(&bind.CallOpts{Context: ctx, Pending: false}, delegate, vault, contract)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if ok {
		return ok, nil
	}

	return false, nil
}

// CheckDelegateForToken returns true if the address is delegated to act on your behalf for a specific token, the
// token's contract or an entire vault.
func (s *Service) CheckDelegateForToken(ctx context.Context, delegateAddr string, vaultAddr string, contractAddr string, token *big.Int) (bool, error) {
	delegate := common.HexToAddress(delegateAddr)
	vault := common.HexToAddress(vaultAddr)
	contract := common.HexToAddress(contractAddr)
	ok, err := s.delegateEthService.CheckDelegateForToken(&bind.CallOpts{Context: ctx, Pending: false}, delegate, vault, contract, token)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if ok {
		return ok, nil
	}

	return false, nil
}

// CheckDelegateForAll returns true if the address is delegated to act on the entire vault.
// delegate is the hot wallet to act on your behalf.
// vault is the cold wallet who issued the delegation.
func (s *Service) CheckDelegateForAll(ctx context.Context, delegateAddr string, vaultAddr string) (bool, error) {
	delegate := common.HexToAddress(delegateAddr)
	vault := common.HexToAddress(vaultAddr)
	ok, err := s.delegateEthService.CheckDelegateForAll(&bind.CallOpts{Context: ctx, Pending: false}, delegate, vault)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if ok {
		return ok, nil
	}

	return false, nil

}
