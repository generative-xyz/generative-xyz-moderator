package structure

import "time"

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
	ID           string     `json:"id"`
	CreatedAt    *time.Time `json:"createdAt"`
	Status       string     `json:"status"`
	FileURI      string     `json:"fileURI"`
	ProjectImage string     `json:"projectImage"`
	ProjectID    string     `json:"projectID"`
	ProjectName  string     `json:"projectName"`
}
