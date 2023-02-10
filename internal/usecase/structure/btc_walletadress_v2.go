package structure

type BctWalletAddressDataV2 struct {
	WalletAddress string `json:"walletAddress"`
	Name string
	File string
	FeeRate int32
}
type BctMintDataV2 struct {
	Address string `json:"address"` //ord_walletaddress
}

type FilterBctWalletAddressesV2 struct {
	BaseFilters
}

type CheckBalanceV2 struct {
	Address string `json:"address"`
}
