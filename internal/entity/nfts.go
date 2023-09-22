package entity

type NFT struct {
	Name              string            `json:"name"`
	TokenId           string            `json:"token_id"`
	CollectionAddress string            `json:"collection_address"`
	TokenUri          string            `json:"token_uri"`
	Owner             string            `json:"owner"`
	Attributes        []TokenUriAttrStr `json:"attributes"`
	Metadata          *NFTMetadata      `json:"metadata"`
}

type NFTMetadata struct {
	Attributes   interface{} `json:"attributes"`
	Description  string      `json:"description"`
	ExternalUrl  interface{} `json:"external_url"`
	Image        string      `json:"image"`
	AnimationURL *string     `json:"animation_url"`
	Name         string      `json:"name"`
}
