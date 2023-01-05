package structure

type FilterProjects struct {
	BaseFilters
	WalletAddress *string
}

type CreateProjectReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
}

type UpdateProjectReq struct {
	MaxSupply int `bson:"maxSupply"`
	LimitSupply int `bson:"limitSupply"`
	MintPrice string `bson:"mintPrice"`
	Name string `bson:"name"`
	CreatorName string `bson:"creatorName"`
	Description string `bson:"description"`
	Thumbnail string `bson:"thumbnail"`
	ThirdPartyScripts []string `bson:"thirdPartyScripts"`
	Scripts []string `bson:"scripts"`
	ReservationList []string `bson:"reservationList"`
	MintFee int `bson:"mintFee"`
	OpenMintUnixTimestamp int `bson:"openMintUnixTimestamp"`
	TokenDescription string `bson:"tokenDescription"`
}

type GetProjectReq struct {
	ContractAddr string
	TokenID string
}
