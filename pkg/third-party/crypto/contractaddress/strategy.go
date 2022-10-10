package contractaddress

import (
	"context"
	"math/big"
)

type Strategy interface {
	GenerateContractAddress() (*ContractAddress, error)
	CheckBalance(address string) (float64, error)
	Transfer(req *TransferRequest) (string, error)
	SuggestGasPrice(context context.Context) (*big.Int, error)
	EstimateGas(req *TransferRequest) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionID string) (*TransactionReceipt, error)
	GetNativeStrategy() Strategy
}
