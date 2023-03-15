package response

type ProjectListing struct {
	ObjectID               string                  `json:"objectID"`
	ContractAddress        string                  `json:"contractAddress"`
	Project                *ProjectInfo            `json:"project"`
	TotalSupply            int64                   `json:"totalSupply"`
	NumberOwners           int64                   `json:"numberOwners"`
	NumberOwnersPercentage float64                 `json:"numberOwnersPercentage"`
	FloorPrice             string                  `json:"floorPrice"`
	FloorPriceOneDay       *VolumneObject          `json:"floorPriceOneDay"`
	FloorPriceOneWeek      *VolumneObject          `json:"floorPriceOneWeek"`
	VolumeFifteenMinutes   *VolumneObject          `json:"volumeFifteenMinutes"`
	VolumeOneDay           *VolumneObject          `json:"volumeOneDay"`
	VolumeOneWeek          *VolumneObject          `json:"volumeOneWeek"`
	ProjectMarketplaceData *ProjectMarketplaceData `json:"projectMarketplaceData"`
	Owner                  *OwnerInfo              `json:"owner"`
}

type OwnerInfo struct {
	WalletAddress           string `json:"walletAddress,omitempty"`
	WalletAddressPayment    string `json:"walletAddress_payment,omitempty"`
	WalletAddressBTC        string `json:"walletAddress_btc,omitempty"`
	WalletAddressBTCTaproot string `json:"walletAddress_btc_taproot,omitempty"`
	DisplayName             string `json:"displayName,omitempty"`
	Avatar                  string `json:"avatar"`
}

type ProjectInfo struct {
	Name            string `json:"name"`
	TokenId         string `json:"tokenId"`
	Thumbnail       string `json:"thumbnail"`
	ContractAddress string `json:"contractAddress"`
}

type VolumneObject struct {
	Amount            string  `json:"amount"`
	PercentageChanged float64 `json:"percentageChanged"`
}

type SearchResponse struct {
	ObjectType  string                `json:"objectType"`
	Inscription *SearhcInscription    `json:"inscription"`
	Project     *ProjectResp          `json:"project"`
	Artist      *ArtistResponse       `json:"artist"`
	TokenUri    *InternalTokenURIResp `json:"tokenUri"`
}

type SearhcInscription struct {
	ObjectId       string          `json:"objectId"`
	InscriptionId  string          `json:"inscriptionId"`
	ContentType    string          `json:"contentType"`
	Number         int64           `json:"number"`
	Sat            float64         `json:"sat"`
	Chain          string          `json:"chain"`
	Address        string          `json:"address"`
	GenesisFee     int64           `json:"genesisFee"`
	GenesisHeight  int64           `json:"genesisHeight"`
	Timestamp      string          `json:"timestamp"`
	ProjectName    string          `json:"projectName"`
	ProjectTokenId string          `json:"projectTokenId"`
	Buyable        bool            `json:"buyable"`
	PriceBTC       string          `json:"priceBtc"`
	Owner          *ArtistResponse `json:"owner"`
}

type SearchArtist struct {
	ObjectId      string `json:"objectId"`
	WalletAddress string `json:"walletAddress"`
	DisplayName   string `json:"displayName"`
	Bio           string `json:"bio"`
	Avatar        string `json:"avatar"`
}

type SearchTokenUri struct {
	ObjectId         string `json:"objectId"`
	TokenId          string `json:"tokenId"`
	InscriptionIndex string `json:"inscriptionIndex"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Image            string `json:"image"`
	ProjectName      string `json:"projectName"`
	ProjectId        string `json:"projectId"`
	Thumbnail        string `json:"thumbnail"`
}

type SearchProject struct {
	ObjectId        string   `json:"objectId"`
	TokenId         string   `json:"tokenId"`
	Name            string   `json:"name"`
	Image           string   `json:"image"`
	CreatorAddrr    string   `json:"creatorAddrr"`
	ContractAddress string   `json:"contractAddress"`
	Categories      []string `json:"categories"`
	Index           int64    `json:"index"`
	MintPrice       string   `json:"mintPrice"`
	MaxSupply       int64    `json:"maxSupply"`
}
