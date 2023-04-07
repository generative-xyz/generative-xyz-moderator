package request

type SendNewBidNotifyRequest struct {
	WalletAddress       string  `json:"wallet_address"`
	BidPrice            float64 `json:"bid_price"`
	Quantity            int     `json:"quantity"`
	CollectorRedirectTo string  `json:"collector_redirect_to"`
}
