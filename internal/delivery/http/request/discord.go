package request

type SendNewBidNotifyRequest struct {
	WalletAddress string  `json:"wallet_address"`
	BidPrice      float64 `json:"bid_price"`
}
