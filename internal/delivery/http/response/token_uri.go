package response

import (
	"rederinghub.io/internal/entity"
	"time"

	"rederinghub.io/internal/usecase/structure"
)

type TokenURIResp struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Image        *string     `json:"image"`
	AnimationURL string      `json:"animation_url"`
	Attributes   interface{} `json:"attributes"`
}

type TokenTraitsResp struct {
	Attributes interface{} `json:"attributes"`
}

type TokenStat struct {
	Price *string `json:"price,omitempty"`
}

type ExternalTokenURIResp struct {
	BaseResponse
	TokenID          string `json:"tokenID"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Image            string `json:"image"`
	AnimationURL     string `json:"animationUrl"`
	InscriptionIndex string `json:"inscriptionIndex"`
	Buyable          bool   `json:"buyable"`
	IsCompleted      bool   `json:"isCompleted"`
	PriceBTC         string `json:"priceBTC"`
	OrderID          string `json:"orderID"`
	ProjectName      string `json:"projectName"`
	ProjectID        string `json:"projectID"`
	Thumbnail        string `json:"thumbnail"`
	Priority         int    `json:"priority"`
}

type InternalTokenURIResp struct {
	BaseResponse
	TokenID               string           `json:"tokenID"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Image                 string           `json:"image"`
	AnimationURL          string           `json:"animationUrl"`
	Attributes            interface{}      `json:"attributes"`
	MintedTime            time.Time        `json:"mintedTime"`
	GenNFTAddr            string           `json:"genNFTAddr"`
	OwnerAddr             string           `json:"ownerAddr"`
	Owner                 *ProfileResponse `json:"owner"`
	Project               *ProjectResp     `json:"project"`
	Creator               *ProfileResponse `json:"creator"`
	Thumbnail             string           `json:"thumbnail"`
	Priority              int              `json:"priority"`
	Stats                 TokenStat        `json:"stats"`
	InscriptionIndex      string           `json:"inscriptionIndex"`
	OrderInscriptionIndex int              `json:"orderInscriptionIndex"`
	OrdinalsData          *OrdinalsData    `json:"ordinalsData"`
	NftTokenId            string           `json:"nftTokenId"`
	InscribedBy           string           `json:"inscribedBy"`
	Seed                  string           `json:"seed"`

	// for buyable:
	Buyable      bool   `json:"buyable"`
	IsCompleted  bool   `json:"isCompleted"`
	SellVerified bool   `json:"sell_verified"`
	PriceBTC     string `json:"priceBTC"`
	PriceETH     string `json:"priceETH"`
	OrderID      string `json:"orderID"`

	ListingDetail *structure.MarketplaceNFTDetail `json:"listingDetail"`
	PriceBrc20    entity.PriceBRC20Obj            `json:"priceBrc20""`

	TokenIDData string `json:"tokenIDData"`
}

type OrdinalsData struct {
	Sat           float64 `json:"sat"`
	ContentType   string  `json:"contentType"`
	ContentLength string  `json:"contentLength"`
	Timestamp     string  `json:"timestamp"`
	Block         int64   `json:"block"`
}

type InternalTokenTraitsResp struct {
	Attributes interface{} `json:"attributes"`
}

type TokenActivityResp struct {
	UUID          string                `json:"uuid"`
	Type          int                   `json:"type"`
	Title         string                `json:"title"`
	UserAAddress  string                `json:"user_a_address"`
	UserA         *ProfileResponse      `json:"user_a"`
	UserBAddress  string                `json:"user_b_address"`
	UserB         *ProfileResponse      `json:"user_b"`
	Amount        int64                 `json:"amount"`
	Time          *time.Time            `json:"time"`
	TokenInfo     *InternalTokenURIResp `json:"token_info"`
	InscriptionID string                `json:"inscription_id"`
	ProjectID     string                `json:"project_id"`
}
