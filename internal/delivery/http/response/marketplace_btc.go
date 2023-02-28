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
	ServiceFee string `json:"serviceFee"`
	RoyaltyFee string `json:"royaltyFee"`
}
