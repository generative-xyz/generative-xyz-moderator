package request

type CreateProjectReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
}


type CreateBTCProjectReq struct {
	MaxSupply             int64              `json:"maxSupply"`
	LimitSupply           int64              `json:"limitSupply"`
	MintPrice             string             `json:"mintPrice"`
	Name                  string             `json:"name"`
	CreatorName           string             `json:"creatorName"`
	CreatorAddrr          string             `json:"creatorAddrr"`
	CreatorAddrrBTC   		string        		 `json:"creatorAddrrBTC"`
	Description           string             `json:"description"`
	OpenMintUnixTimestamp  int                `json:"openMintUnixTimestamp"`
	CloseMintUnixTimestamp int                `json:"closeMintUnixTimestamp"`
	Thumbnail             string             `json:"thumbnail"`
	ThirdPartyScripts     []string           `json:"thirdPartyScripts"`
	Scripts               []string           `json:"scripts"`
	TokenDescription      string             `json:"tokenDescription"`
	Styles                string             `json:"styles"`
	SocialWeb             string             `json:"socialWeb"`
	SocialTwitter         string             `json:"socialTwitter"`
	SocialDiscord         string             `json:"socialDiscord"`
	SocialMedium          string             `json:"socialMedium"`
	SocialInstagram       string             `json:"socialInstagram"`
	License               string             `json:"license"`
	Tags                  []string           `json:"tags"`
	Categories            []string           `json:"categories"`
	ZipLink				  *string			  `json:"zipLink"`
	AnimationURL		  *string			  `json:"animationURL"`
	Royalty		  int			  `json:"royalty"`
	IsFullChain        bool        `json:"isFullChain"` 
}

type UpdateBTCProjectReq struct {
	ProjectID                 	 *string             `json:"projectID"`
	Name                  *string             `json:"name"`
	Description           *string             `json:"description"`
	Thumbnail             *string             `json:"thumbnail"`
	IsHidden 			*bool	 `json:"isHidden"`
	Royalty 				*int64	 `json:"royalty"`
	MintPrice 				*string	 `json:"mintPrice"`
	MaxSupply 				*int64	 `json:"maxSupply"`
}

type UpdateProjectReq struct {
	Priority *int `json:"priority"`
}
