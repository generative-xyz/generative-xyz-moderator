package entity

type ItemListing struct {
	InscriptionId         string         `json:"inscription_id"`
	VolumeOneHour         *VolumneObject `json:"volumeOneHour"`
	VolumeOneDay          *VolumneObject `json:"volumeOneDay"`
	VolumeOneWeek         *VolumneObject `json:"volumeOneWeek"`
	VolumeOneMonth        *VolumneObject `json:"volumeOneMonth"`
	SellerAddress         string         `json:"sellerAddress"`
	SellerDisplayName     string         `json:"sellerDisplayName"`
	BuyerAddress          string         `json:"buyerAddress"`
	BuyerDisplayName      string         `json:"buyerDisplayName"`
	Name                  string         `json:"name"`
	Image                 string         `json:"image"`
	ContentType           string         `json:"contentType"`
	OrderInscriptionIndex float64        `json:"orderInscriptionIndex"`
	InscriptionIndex      float64        `json:"inscriptionIndex"`
}

type VolumneObject struct {
	Amount            string  `json:"amount"`
	PercentageChanged float64 `json:"-"`
}
