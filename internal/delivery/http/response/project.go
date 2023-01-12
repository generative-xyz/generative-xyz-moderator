package response

import "time"

type TraitValueStat struct {
	Value  string `json:"value"`
	Rarity int32  `json:"rarity"`
}

type TraitStat struct {
	TraitName       string           `json:"traitName"`
	TraitValuesStat []TraitValueStat `json:"traitValuesStat"`
}

type ProjectResp struct{
	BaseResponse
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	MaxSupply int64 `json:"maxSupply"`
	Limit int64 `json:"limit"`
	MintPrice string `json:"mintPrice"`
	MintPriceAddr string `json:"mintPriceAddr"`
	Name string `json:"name"`
	CreatorAddr string `json:"creatorAddr"`
	License string `json:"license"`
	Desc string `json:"desc"`
	Image string `json:"image"`
	ScriptType []string `json:"scriptType"`
	Social interface{} `json:"social"`
	Scripts []string `json:"scripts"`
	Styles string `json:"styles"`
	CompleteTime int64  `json:"completeTime"`
	GenNFTAddr string `json:"genNFTAddr"`
	ItemDesc string `json:"itemDesc"`
	Status bool `json:"status"`
	NftTokenURI string `json:"projectURI"`
	MintingInfo NftMintingDetail `json:"mintingInfo"`
	Royalty int `json:"royalty"`
	Reservers []string `json:"reservers"`
	CreatorProfile ProfileResponse `json:"creatorProfile"`
	BlockNumberMinted *string           `json:"blockNumberMinted"`
	MintedTime       *time.Time        `json:"mintedTime"`
	Stats ProjectStatResp `json:"stats"`
	TraitStat []TraitStat `json:"traitStat"`
}

type ProjectStatResp struct {
	UniqueOwnerCount uint32 `json:"uniqueOwnerCount"`
	TotalTradingVolumn string `json:"totalTradingVolumn"`
	FloorPrice         string `json:"floorPrice"`
	BestMakeOfferPrice string `json:"bestMakeOfferPrice"`
	ListedPercent      int32 `json:"listedPercent"`
}

type NftMintingDetail struct {
	Index           int64 `json:"index"`
	IndexReserve    int64 `json:"indexReserve"`
}
