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

type ExistedInAllowList struct {
	Existed bool `json:"existed"`
}

type ProjectResp struct {
	BaseResponse
	ContractAddress           string           `json:"contractAddress"`
	TokenID                   string           `json:"tokenID"`
	MaxSupply                 int64            `json:"maxSupply"`
	Limit                     int64            `json:"limit"`
	MintPrice                 string           `json:"mintPrice"`
	MintPriceEth              string           `json:"mintPriceEth"`
	NetworkFee                string           `json:"networkFee"`
	NetworkFeeEth             string           `json:"networkFeeEth"`
	MintPriceAddr             string           `json:"mintPriceAddr"`
	Name                      string           `json:"name"`
	CreatorAddrr              string           `json:"creatorAddr"`
	CreatorAddrrBTC           string           `json:"creatorAddrrBTC"`
	Categories                []string         `json:"categories"`
	License                   string           `json:"license"`
	Desc                      string           `json:"desc"`
	Image                     string           `json:"image"`
	ScriptType                []string         `json:"scriptType"`
	Social                    interface{}      `json:"social"`
	Scripts                   []string         `json:"scripts"`
	Styles                    string           `json:"styles"`
	CompleteTime              int64            `json:"completeTime"`
	GenNFTAddr                string           `json:"genNFTAddr"`
	ItemDesc                  string           `json:"itemDesc"`
	Status                    bool             `json:"status"`
	IsFullChain               bool             `json:"isFullChain"`
	NftTokenURI               string           `json:"projectURI"`
	MintingInfo               NftMintingDetail `json:"mintingInfo"`
	Royalty                   int              `json:"royalty"`
	Reservers                 []string         `json:"reservers"`
	CreatorProfile            ProfileResponse  `json:"creatorProfile"`
	BlockNumberMinted         *string          `json:"blockNumberMinted"`
	MintedTime                *time.Time       `json:"mintedTime"`
	Stats                     ProjectStatResp  `json:"stats"`
	TraitStat                 []TraitStat      `json:"traitStat"`
	Priority                  int              `json:"priority"`
	OpenMintUnixTimestamp     int              `json:"openMintUnixTimestamp"`
	CloseMintUnixTimestamp    int              `json:"closeMintUnixTimestamp"`
	WhiteListEthContracts     []string         `json:"whiteListEthContracts"`
	IsHidden                  bool             `json:"isHidden"`
	EditableIsHidden          bool             `json:"editableIsHidden"`
	TotalImages               int              `json:"totalImages"`
	ReportUsers               []*ReportProject `json:"reportUsers"`
	AnimationHtml             *string          `json:"animationHtml"`
	MaxFileSize               int64            `json:"maxFileSize"`
	CaptureThumbnailDelayTime int              `json:"captureThumbnailDelayTime"`
	FromAuthentic             bool             `json:"fromAuthentic"`
	TokenAddress              string           `json:"tokenAddress"`
	TokenId                   string           `json:"nftTokenId"`
	OwnerOf                   string           `json:"ownerOf"`
	OrdinalsTx                string           `json:"ordinalsTx"`
	InscribedBy               string           `json:"inscribedBy"`
	HtmlFile                  string           `json:"htmlFile"`
	LimitMintPerProcess       int              `json:"limitMintPerProcess"`
	FileExtension             string           `json:"fileExtension"`
	BtcFloorPrice             uint64           `json:"btcFloorPrice"`
	ReserveMintPrice          string           `json:"reserveMintPrice"`
	ReserveMintLimit          int              `json:"reserveMintLimit"`
	IsReviewing               bool             `json:"isReviewing"`
	CanCreateProposal         bool             `json:"canCreateProposal"`
	ProposalSeqId             *uint            `json:"proposalSeqId"`
	IsGenerative              bool             `json:"isGenerative"` // art by generative code or html file
}

type ReportProject struct {
	OriginalLink      string `json:"originalLink"`
	ReportUserAddress string `json:"reportUserAddress"`
}

type ProjectStatResp struct {
	UniqueOwnerCount   uint32 `json:"uniqueOwnerCount"`
	TotalTradingVolumn string `json:"totalTradingVolumn"`
	FloorPrice         string `json:"floorPrice"`
	BestMakeOfferPrice string `json:"bestMakeOfferPrice"`
	ListedPercent      int32  `json:"listedPercent"`
}

type NftMintingDetail struct {
	Index        int64 `json:"index"`
	IndexReserve int64 `json:"indexReserve"`
}

type ProjectMarketplaceData struct {
	Listed          uint64  `json:"listed"`
	FloorPrice      uint64  `json:"floor_price"`
	TotalVolume     uint64  `json:"volume"`
	MintVolume      uint64  `json:"mint_volume"`
	CEXVolume       uint64  `json:"cex_volume"`
	FirstSaleVolume float64 `json:"first_sale_volume"`
}
