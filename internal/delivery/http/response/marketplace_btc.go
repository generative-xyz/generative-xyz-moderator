package response

type CreateMarketplaceBTCBuyOrder struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
	Price          string `json:"price"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
}

type ListingFee struct {
	ServiceFee     string `json:"serviceFee"`
	RoyaltyFee     string `json:"royaltyFee"`
	RoyaltyAddress string `json:"royaltyAddress"`
	ServiceAddress string `json:"serviceAddress"`
	ProjectID      string `json:"projectID"`
}

type StatFirstSale struct {
	Amount          string            `json:"amount"`
	AmountByPaytype map[string]string `json:"amountByPaytype"`
	ProjectID       string            `json:"projectID"`
}
