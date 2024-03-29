package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

const BITCOIN_PROJECT_ID_START_WITH = 1000001
const DEFAULT_FEE_RATE = 15
const DEFAULT_CAPTURE_TIME int = 20
const DEFAULT_DELAY_OPEN_MINT_TIME_IN_HOUR = 3 // hours

type TraitValueStat struct {
	Value   string  `bson:"value" json:"value"`
	Rarity  int32   `bson:"rarity" json:"rarity"`
	RarityF float64 `bson:"rarity_number" json:"rarity_number"`
}

type TraitStat struct {
	TraitName       string           `bson:"traitName" json:"traitName"`
	TraitValuesStat []TraitValueStat `bson:"traitValuesStat" json:"traitValuesStat"`
}

type ProjectStat struct {
	LastTimeSynced   *time.Time `bson:"lastTimeSynced" json:"lastTimeSynced"`
	UniqueOwnerCount uint32     `bson:"uniqueOwnerCount" json:"uniqueOwnerCount"`
	// TODO add other stats here
	TotalTradingVolumn string `bson:"totalTradingVolumn" json:"totalTradingVolumn"`
	FloorPrice         string `bson:"floorPrice" json:"floorPrice"`
	BestMakeOfferPrice string `bson:"bestMakeOfferPrice" json:"bestMakeOfferPrice"`
	ListedPercent      int32  `bson:"listedPercent" json:"listedPercent"`
	MintedCount        int32  `bson:"minted_count" json:"mintedCount"`
	TrendingScore      int64  `bson:"trending_score" json:"trendingScore"`

	//volume, buyable for sorting
	Volume  float64 `bson:"volume" json:"volume"`
	Buyable bool    `bson:"buyable" json:"buyable"`
}

type MaxProjectID struct {
	ID int64 `bson:"_id"`
}

type Projects struct {
	BaseEntity               `bson:",inline" json:"-"`
	ContractAddress          string `bson:"contractAddress" json:"contractAddress"`
	TokenID                  string `bson:"tokenid" json:"tokenID"`
	TokenIDInt               int64  `bson:"tokenIDInt" json:"tokenIDInt"`
	MaxSupply                int64  `bson:"maxSupply" json:"maxSupply"`
	LimitSupply              int64  `bson:"limitSupply" json:"limitSupply"`
	MintPrice                string `bson:"mintPrice" json:"mintPrice"`
	MintPriceEth             string
	NetworkFeeEth            string
	NetworkFee               string              `bson:"networkFee" json:"networkFee"`
	MaxFileSize              int64               `bson:"maxFileSize" json:"maxFileSize"`
	Name                     string              `bson:"name" json:"name"`
	CreatorName              string              `bson:"creatorName" json:"creatorName"`
	CreatorAddrr             string              `bson:"creatorAddress" json:"creatorAddrr"`
	CreatorAddrrBTC          string              `bson:"creatorAddrrBTC" json:"creatorAddrrBTC"`
	Description              string              `bson:"description" json:"description"`
	OpenMintUnixTimestamp    int                 `bson:"openMintUnixTimestamp" json:"openMintUnixTimestamp"`
	CloseMintUnixTimestamp   int                 `bson:"closeMintUnixTimestamp" json:"closeMintUnixTimestamp"`
	Thumbnail                string              `bson:"thumbnail" json:"thumbnail"`
	ThirdPartyScripts        []string            `bson:"thirdPartyScripts" json:"thirdPartyScripts"`
	Scripts                  []string            `bson:"scripts" json:"scripts"`
	MintFee                  int                 `bson:"mintFee" json:"mintFee"`
	TokenDescription         string              `bson:"tokenDescription" json:"tokenDescription"`
	Styles                   string              `bson:"styles" json:"styles"`
	Royalty                  int                 `bson:"royalty" json:"royalty"`
	SocialWeb                string              `bson:"socialWeb" json:"socialWeb"`
	SocialTwitter            string              `bson:"socialTwitter" json:"socialTwitter"`
	SocialDiscord            string              `bson:"socialDiscord" json:"socialDiscord"`
	SocialMedium             string              `bson:"socialMedium" json:"socialMedium"`
	SocialInstagram          string              `bson:"socialInstagram" json:"socialInstagram"`
	License                  string              `bson:"license" json:"license"`
	GenNFTAddr               string              `bson:"genNFTAddr" json:"genNFTAddr"`
	MintTokenAddress         string              `bson:"mintTokenAddress" json:"mintTokenAddress"`
	Hash                     string              `bson:"hash" json:"hash"`
	Tags                     []string            `bson:"tags" json:"tags"`
	Categories               []string            `bson:"categories" json:"categories"`
	Status                   bool                `bson:"status" json:"status"`
	NftTokenUri              string              `bson:"nftTokenUri" json:"nftTokenUri"`
	IsSynced                 bool                `bson:"isSynced" json:"isSynced"`
	MintingInfo              ProjectMintingInfo  `bson:",inline" json:"mintingInfo"`
	CompleteTime             int64               `bson:"completeTime" json:"completeTime"`
	Reservers                []string            `bson:"reservers" json:"reservers"`
	CreatorProfile           Users               `bson:"creatorProfile" json:"creatorProfile"`
	BlockNumberMinted        *string             `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime               *time.Time          `bson:"minted_time" json:"minted_time"`
	Stats                    ProjectStat         `bson:"stats" json:"stats"`
	TraitsStat               []TraitStat         `bson:"traitsStat" json:"traitsStat"`
	Priority                 *int                `bson:"priority" json:"priority"`
	IsHidden                 bool                `bson:"isHidden" json:"isHidden"`
	Images                   []string            `bson:"images" json:"images"`                               //if user uses links instead of animation URL
	WhiteListEthContracts    []string            `bson:"whiteListEthContracts" json:"whiteListEthContracts"` //if user uses links instead of animation URL
	ProcessingImages         []string            `bson:"processingImages" json:"processingImages"`
	MintedImages             []MintedImages      `bson:"mintedImages" json:"mintedImages"` //if user uses links instead of animation URL
	IsFullChain              bool                `bson:"isFullChain" json:"isFullChain"`
	TraceID                  string              `bson:"traceID" json:"traceID"` //TO find log easily
	ReportUsers              []*ReportProject    `bson:"reportUsers" json:"reportUsers"`
	InscriptionIcon          string              `bson:"inscription_icon" json:"inscriptionIcon"`
	CreatedByCollectionMeta  bool                `bson:"created_by_collection_meta" json:"created_by_collection_meta"`
	Source                   string              `bson:"source" json:"source"`
	AnimationHtml            *string             `bson:"animation_html"`
	CatureThumbnailDelayTime *int                `bson:"cature_thumbnail_delay_time"`
	ReserveMintPrice         string              `bson:"reserveMintPrice" json:"reserveMintPrice"`
	ReserveMintLimit         int                 `bson:"reserveMintLimit" json:"reserveMintLimit"`
	FromAuthentic            bool                `bson:"fromAuthentic"`
	TokenAddress             string              `bson:"tokenAddress"`
	TokenId                  string              `bson:"tokenId"`
	OwnerOf                  string              `bson:"ownerOf"`
	OrdinalsTx               string              `bson:"ordinalsTx"`
	InscribedBy              string              `bson:"inscribedBy"`
	HtmlFile                 string              `bson:"htmlFile"`
	LimitMintPerProcess      int                 `bson:"limitMintPerProcess"`
	TxHash                   string              `bson:"txhash"`
	TxHex                    string              `bson:"txHex"`
	CommitTxHash             string              `bson:"commitTxHash"`
	RevealTxHash             string              `bson:"revealTxHash"`
	IsBigFile                bool                `bson:"isBigFile"`
	HasZipFile               bool                `bson:"hasZipFile"`
	AuctionWinnerList        []AuctionWinnerList `bson:"-" json:"auctionWinnerList"`

	// GM whitelist holder
	IsSupportGMHolder bool   `json:"isSupportGMHolder" bson:"isSupportGMHolder"`
	MinimumGMSupport  string `json:"minimumGMSupport" bson:"minimumGMSupport"`

	CurrentLoginUserID string `json:"-" bson:"-"`
}

func (p *Projects) IsMintTC() bool {
	return p.TokenIDInt < 1000000
}

type ProjectsHaveMinted struct {
	TokenID      string `bson:"tokenid" json:"tokenID"`
	Name         string `bson:"name" json:"name"`
	Index        int    `bson:"index" json:"index"`
	MintPrice    string `bson:"mintPrice"`
	MintPriceEth string `bson:"mintpriceeth"`
	CreatorAddrr string `bson:"creatorAddress" json:"creatorAddrr"`
}

type ReportProject struct {
	OriginalLink      string `bson:"originalLink" json:"originalLink"`
	ReportUserAddress string `bson:"reportUserAddress" json:"reportUserAddress"`
}

type MintedImages struct {
	URI         string     `bson:"uri"`
	Commit      string     `bson:"commit"`
	Inscription string     `bson:"inscription"`
	Reveal      string     `bson:"reveal"`
	Fees        int        `bson:"fees"`
	IsSent      bool       `bson:"isSent"`
	MintedAt    *time.Time `bson:"mintedAt"`
}

type ProjectMintingInfo struct {
	Index        int64 `bson:"index"`
	IndexReverse int64 `bson:"indexReverse"`
}

type FilterProjects struct {
	BaseFilters
	Search          *string
	WalletAddress   *string
	ContractAddress *string
	Name            *string
	IsHidden        *bool
	Status          *bool
	IsSynced        *bool
	CategoryIds     []string
	TokenIds        []string
	Ids             []string
	CustomQueries   map[string]bson.M
	TxHash          *string
	TxHex           *string
	CommitTxHash    *string
	RevealTxHash    *string
}

type ProjectListed struct {
	ID     string `json:"_id" bson:"_id"`
	Listed int    `bson:"listed" json:"listed"`
}

type ProjectFloorPrice struct {
	ID     string  `json:"_id" bson:"_id"`
	Amount float64 `bson:"total_amount" json:"total_amount"`
}

type ProjectVolume struct {
	ID     string  `json:"_id" bson:"_id"`
	Amount float64 `bson:"Amount" json:"Amount"`
}

func (u Projects) TableName() string {
	return utils.COLLECTION_PROJECTS
}

func (u Projects) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
