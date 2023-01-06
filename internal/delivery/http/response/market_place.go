package response

type MarketplaceListing struct {
	BaseResponse         `json:",inline"`
	ID                 string `json:"id"`
	OfferingId         string `json:"offeringID"`
	CollectionContract string `json:"collectionContract"`
	TokenId            string `json:"tokenID"`
	Seller             string `json:"seller"`
	Erc20Token         string `json:"erc20Token"`
	Price              string `json:"price"`
	Closed             bool   `json:"closed"`
	Finished           bool   `json:"finished"`
	DurationTime       string `json:"durationTime"`
	Token InternalTokenURIResp `json:"token"`
}

type MarketplaceOffer struct {
	BaseResponse         `json:",inline"`
	ID                 string `json:"id"`
	OfferingId         string `json:"offeringID"`
	CollectionContract string `json:"collectionContract"`
	TokenId            string `json:"tokenID"`
	Buyer             string `json:"buyer"`
	Erc20Token         string `json:"erc20Token"`
	Price              string `json:"price"`
	Closed             bool   `json:"closed"`
	Finished           bool   `json:"finished"`
	DurationTime       string `json:"durationTime"`
	Token InternalTokenURIResp `json:"token"`
}

type MarketplaceStatResp struct {
	Stats ProjectStatResp `json:"stats"`
}
