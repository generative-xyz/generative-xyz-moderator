package structure

type BctWalletAddressData struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID     string `json:"projectID"`
}
type BctMintData struct {
	Address string `json:"address"` //ord_walletaddress
}

type FilterBctWalletAddresses struct {
	BaseFilters
}

type CheckBalance struct {
	Address string `json:"address"`
}

type MintingInscription struct {
	Status    string `json:"status"`
	FileURI   string `json:"fileURI"`
	ProjectID string `json:"projectID"`
}
