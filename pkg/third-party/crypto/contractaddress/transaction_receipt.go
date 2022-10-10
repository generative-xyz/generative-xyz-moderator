package contractaddress

import "rederinghub.io/pkg/third-party/crypto/constants/cryptotransactionstatus"

type TransactionReceipt struct {
	Status        uint64
	TransactionID string
	Coin          string
	RawData       string
}

func (t *TransactionReceipt) IsSuccess() bool {
	return t.Status == cryptotransactionstatus.Success
}

func (t *TransactionReceipt) IsFailure() bool {
	return t.Status == cryptotransactionstatus.Failure
}
