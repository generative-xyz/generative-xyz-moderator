package response

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
