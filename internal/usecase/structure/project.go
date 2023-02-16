package structure

type FilterProjects struct {
	BaseFilters
	WalletAddress *string
	Name *string
	CategoryIds []string
}

type FilterProposal struct {
	BaseFilters
	Proposer *string
	State *int
	ProposalID *string
}

type FilterProposalVote struct {
	BaseFilters
	Voter *string
	Support *int
	ProposalID *string
}

type CreateProjectReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
}

type CreateBtcProjectReq struct {
	MaxSupply             int64              `json:"maxSupply"`
	LimitSupply           int64              `json:"limitSupply"`
	MintPrice             string             `json:"mintPrice"`
	Name                  string             `json:"name"`
	CreatorName           string             `json:"creatorName"`
	CreatorAddrr          string             `json:"creatorAddrr"`
	CreatorAddrrBTC          string           `json:"creatorAddrrBTC"`
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
	ZipLink				   *string			  `json:"zipLink"`
	AnimationURL		*string			  `json:"animationURL"`
	Royalty		  		int			  `json:"royalty"`
	IsFullChain        	bool        `json:"isFullChain"`
}

type CreateProposaltReq struct {
	Title string `json:"title"`
	Description string `json:"description"`
	TokenType string `json:"tokenType"`
	Amount string `json:"amount"`
	ReceiverAddress string `json:"receiverAddress"`
}

type UpdateProjectReq struct {
	TokenID string `json:"tokenID"`
	Priority *int `json:"priority"`
	ContracAddress string `json:"contractAddress"`
}

type GetProjectReq struct {
	ContractAddr string
	TokenID string
}

type ProjectAnimationUrl struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	AnimationUrl string `json:"animation_url"`
}
