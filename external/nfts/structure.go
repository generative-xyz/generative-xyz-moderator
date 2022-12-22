package nfts

type  MoralisFilter struct {
	Chain *string `json:"chain"`
	Format *string `json:"format"`
	Limit *int `json:"limit"`
	TotalRanges *int `json:"totalRanges"`
	Range *int `json:"range"`
	Cursor *string `json:"cursor"`
	NormalizeMetadata *bool `json:"normalizeMetadata"`
}

type MoralisTokensResp struct {
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	Cursor string `json:"cursor"`
	Result []MoralisToken `json:"result"`
}

type MoralisToken struct {
	TokenAddress string `json:"token_address"`
	TokenID string `json:"token_id"`
	Amount string `json:"amount"`
	Owner string `json:"owner_of"`
	TokenHash string `json:"token_hash"`
	ContractType string `json:"contract_type"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	TokenUri string `json:"token_uri"`
	MetadataString *string `json:"metadata"`
	Metadata *MoralisTokenMetadata `json:"-"`
}

type MoralisTokenMetadata struct {
	Image string `json:"image"`
	Name string `json:"name"`
	Description string `json:"description"`
	ExternalLink string `json:"external_link"`
	AnimationUrl string `json:"animation_url"`
	Traits interface{} `json:"traits"`
}