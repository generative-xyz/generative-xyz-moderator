package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

const BITCOIN_PROJECT_ID_START_WITH = 1000001

type TraitValueStat struct {
	Value  string `bson:"value" json:"value"`
	Rarity int32  `bson:"rarity" json:"rarity"`
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
	MintedCount        int32  `bson:"minted_count" json:"minted_count"`
}

type Projects struct {
	BaseEntity`bson:",inline"`
	ContractAddress string `bson:"contractAddress"`
	TokenID string `bson:"tokenid"`
	TokenIDInt int64 `bson:"tokenIDInt"`
	MaxSupply int64 `bson:"maxSupply"`
	LimitSupply int64 `bson:"limitSupply"`
	MintPrice string `bson:"mintPrice"`
	Name string `bson:"name"`
	CreatorName string `bson:"creatorName"`
	CreatorAddrr string `bson:"creatorAddress"`
	Description string `bson:"description"`
	Thumbnail string `bson:"thumbnail"`
	ThirdPartyScripts []string `bson:"thirdPartyScripts"`
	Scripts []string `bson:"scripts"`
	ReservationList []string `bson:"reservationList"`
	MintFee int `bson:"mintFee"`
	OpenMintUnixTimestamp int `bson:"openMintUnixTimestamp"`
	TokenDescription string `bson:"tokenDescription"`
	Styles string `bson:"styles"`
	Royalty int `bson:"royalty"`
	SocialWeb string `bson:"socialWeb"`
	SocialTwitter string `bson:"socialTwitter"`
	SocialDiscord string `bson:"socialDiscord"`
	SocialMedium string `bson:"socialMedium"`
	SocialInstagram string `bson:"socialInstagram"`
	License string `bson:"license"`
	GenNFTAddr string `bson:"genNFTAddr"`
	MintTokenAddress string `bson:"mintTokenAddress"`
	Hash string `bson:"hash"`
	Tags []string `bson:"tags"`
	Categories []string `bson:"categories"`
	Status bool `bson:"status"`
	NftTokenUri string `bson:"nftTokenUri"`
	IsSynced bool `bson:"isSynced"`
	MintingInfo ProjectMintingInfo `bson:",inline"`
	CompleteTime int64  `bson:"completeTime"`
	Reservers []string `bson:"reservers"`
	CreatorProfile Users `bson:"creatorProfile"`
	BlockNumberMinted *string `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime *time.Time `bson:"minted_time" json:"minted_time"`
	Stats                 ProjectStat        `bson:"stats"`
	TraitsStat         []TraitStat `bson:"traitsStat" json:"traitsStat"`
	Priority *int `bson:"priority"`
	IsHidden bool `bson:"isHidden"`
}

type	ProjectMintingInfo struct {
	Index int64 `bson:"index"`
	IndexReverse int64 `bson:"indexReverse"`
}

type FilterProjects struct {
	BaseFilters
	WalletAddress *string
	Name *string
}

func (u Projects) TableName() string { 
	return utils.COLLECTION_PROJECTS
}

func (u Projects) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
