package contractaddress

import "math/big"

type TransferRequest struct {
	Payer             *ContractAddress
	FromAddress       *ContractAddress
	ToAddress         string
	Amount            uint64
	GasPrice          *big.Int
	GasLimit          uint64
	Data              []byte
	UseNativeStrategy bool
}
