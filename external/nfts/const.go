package nfts

const (
	KeyOwner       string = "owner"
	KeyLimit       string = "limit"
	KeyOffset      string = "offset"
	KeyCurrsor     string = "currsor"
	KeyChain       string = "chain"
	KeyFormat      string = "format"
	KeyTotalRanges string = "totalRanges"
	KeyRange       string = "range"
	URLAssets      string = "assets"
	URLNft         string = "nft"
)

var ChainToChainID = map[string]string{
	"mumbai": "80001",
	"goerli": "5",
}
