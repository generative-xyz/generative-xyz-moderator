package entity

type NFT struct {
	Name              string            `json:"name"`
	TokenId           string            `json:"token_id"`
	CollectionAddress string            `json:"collection_address"`
	TokenUri          string            `json:"token_uri"`
	Owner             string            `json:"owner"`
	Attributes        []TokenUriAttrStr `json:"attributes"`
}
