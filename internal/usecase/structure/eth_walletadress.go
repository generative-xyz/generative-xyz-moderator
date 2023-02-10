package structure

type EthWalletAddressData struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID     string `json:"projectID"`
}
type EthMintData struct {
	Address string `json:"address"` //ord_walletaddress
}

type FilterEthWalletAddresses struct {
	BaseFilters
}
