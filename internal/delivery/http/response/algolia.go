package response

type SearchResponse struct {
	ObjectType  string             `json:"objectType"`
	Inscription *SearhcInscription `json:"inscription"`
	Project     *SearchProject     `json:"project"`
	Artist      *SearchArtist      `json:"artist"`
	TokenUri    *SearchTokenUri    `json:"tokenUri"`
}

type SearhcInscription struct {
	ObjectId      string `json:"objectId"`
	InscriptionId string `json:"inscriptionId"`
	Number        int64  `json:"number"`
	Sat           string `json:"sat"`
	Chain         string `json:"chain"`
	GenesisFee    int64  `json:"genesisFee"`
	GenesisHeight int64  `json:"genesisHeight"`
	Timestamp     string `json:"timestamp"`
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
	ObjectId    string   `json:"objectId"`
	TokenId     string   `json:"tokenId"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	CreatorAddr string   `json:"creatorAddr"`
	Categories  []string `json:"categories"`
}
