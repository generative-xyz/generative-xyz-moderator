package response

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TokenURIResp struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	AnimationURL string `json:"animation_url"`
	Attributes interface{} `json:"attributes"`
}

type TokenTraitsResp struct{
	Attributes interface{} `json:"attributes"`
}

type ProjectResp struct{
	MaxSupply big.Int `json:"maxSupply"`
	Limit big.Int `json:"limit"`
	MintPrice big.Int `json:"mintPrice"`
	MintPriceAddr common.Address `json:"mintPriceAddr"`
	Name string `json:"name"`
	Creator string `json:"creator"`
	CreatorAddr common.Address `json:"creatorAddr"`
	License string `json:"license"`
	Desc string `json:"desc"`
	Image string `json:"image"`
	ScriptType []string `json:"scriptType"`
	Social interface{} `json:"social"`
	Scripts []string `json:"scripts"`
	Styles string `json:"styles"`
	CompleteTime big.Int  `json:"completeTime"`
	GenNFTAddr common.Address `json:"genNFTAddr"`
	ItemDesc string `json:"itemDesc"`
	Status bool `json:"status"`
	NftTokenURI string `json:"nftTokenURI"`
}
