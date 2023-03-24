package structure

type DexListingOWCollections struct {
	Collections []DexListingOWCollectionInfo `json:"collections"`
	Total       int                          `json:"total"`
}

type DexListingOWCollectionInfo struct {
	Active                bool    `json:"active"`
	ChangeWeek            float64 `json:"change_week"`
	CreatorAddress        any     `json:"creator_address"`
	Description           string  `json:"description"`
	FeaturedPriority      int     `json:"featured_priority"`
	FloorPrice            int     `json:"floor_price"`
	HighestInscriptionNum int     `json:"highest_inscription_num"`
	Icon                  string  `json:"icon"`
	ID                    string  `json:"id"`
	Listed                int     `json:"listed"`
	LowestInscriptionNum  int     `json:"lowest_inscription_num"`
	Name                  string  `json:"name"`
	Slug                  string  `json:"slug"`
	SponsoredPriority     int     `json:"sponsored_priority"`
	TotalSupply           int     `json:"total_supply"`
	VolumeWeek            int     `json:"volume_week"`
	Socials               struct {
		Discord string `json:"discord"`
		Twitter string `json:"twitter"`
		Website string `json:"website"`
	} `json:"socials,omitempty"`
}

type DexListingOWCollectionItem struct {
	// Collection struct {
	// 	CreatorAddress any    `json:"creator_address"`
	// 	Name           string `json:"name"`
	// 	Slug           string `json:"slug"`
	// } `json:"collection"`
	ContentType string `json:"content_type"`
	Escrow      struct {
		BoughtAt      string `json:"bought_at"`
		SatoshiPrice  int    `json:"satoshi_price"`
		SellerAddress string `json:"seller_address"`
	} `json:"escrow"`
	ID string `json:"id"`
	// Meta struct {
	// 	Attributes any    `json:"attributes"`
	// 	Name       string `json:"name"`
	// } `json:"meta"`
	Num int `json:"num"`
}

type DexListingOWPurchaseRespond struct {
	Purchase string `json:"purchase"`
	Setup    string `json:"setup"`
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
}

// {
//     "purchase": "81162e76847495f3800ccd207ae6a7f6375106900f7f8531a97e618d563de89b",
//     "setup": "dbb4976c7bd3667cc5ee1dbed62fe23cae26dbe76eeec118f743854ec5d0d08b",
//     "success": true
// }
