package nftdata

import "fmt"

type (
	NFTCustomerInfo struct {
		CompanyDomain string
		Name          string
		Email         string
		CompanyName   string
		Balance       int64
	}
)

func (i *NFTCustomerInfo) WithBalance(balance int64) *NFTCustomerInfo {
	i.Balance = balance
	return i
}

func NewNFTCustomerInfo(nftItem *Item, address string) *NFTCustomerInfo {
	domain := nftItem.Domain
	name := getNameFromAddress(address)
	email := fmt.Sprintf("%s@%s", name, domain)
	return &NFTCustomerInfo{
		CompanyDomain: domain,
		Name:          name,
		Email:         email,
		CompanyName:   nftItem.CompanyName,
	}
}
